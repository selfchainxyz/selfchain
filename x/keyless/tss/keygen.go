package tss

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/bnb-chain/tss-lib/v2/crypto"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/btcsuite/btcd/btcec/v2"
	"selfchain/x/keyless/types"
)

// KeygenResult contains the result of key generation
type KeygenResult struct {
	Party1Data     *keygen.LocalPartySaveData
	Party2Data     *keygen.LocalPartySaveData
	ChainID        string
	PublicKeyBytes []byte
	PublicKey      *ecdsa.PublicKey
	SecurityLevel  types.SecurityLevel
}

// EncryptedShare represents an encrypted key share
type EncryptedShare struct {
	EncryptedData []byte
	ChainID       string
}

// createParties creates two TSS parties for key generation
func createParties(preParams *keygen.LocalPreParams) (tss.Party, tss.Party, error) {
	// Create party IDs
	p1ID := tss.NewPartyID("1", "P1", big.NewInt(1))
	p2ID := tss.NewPartyID("2", "P2", big.NewInt(2))
	parties := tss.SortPartyIDs([]*tss.PartyID{p1ID, p2ID})

	// Create peer context
	peerCtx := tss.NewPeerContext(parties)
	params1 := tss.NewParameters(tss.S256(), peerCtx, p1ID, len(parties), Threshold)
	params2 := tss.NewParameters(tss.S256(), peerCtx, p2ID, len(parties), Threshold)

	// Create parties
	outCh := make(chan tss.Message, TotalPartyCount)
	endCh := make(chan *keygen.LocalPartySaveData, TotalPartyCount)

	p1 := keygen.NewLocalParty(params1, outCh, endCh, *preParams)
	p2 := keygen.NewLocalParty(params2, outCh, endCh, *preParams)

	return p1, p2, nil
}

// convertECPointToECDSA converts a TSS ECPoint to an ECDSA public key
func convertECPointToECDSA(point *crypto.ECPoint) *ecdsa.PublicKey {
	if point == nil {
		return nil
	}
	return &ecdsa.PublicKey{
		Curve: point.Curve(),
		X:     point.X(),
		Y:     point.Y(),
	}
}

// serializePublicKey converts a TSS public key to bytes
func serializePublicKey(pubKey *ecdsa.PublicKey) []byte {
	if pubKey == nil {
		return nil
	}

	// Convert ECDSA public key to btcec public key
	x := new(btcec.FieldVal)
	x.SetByteSlice(pubKey.X.Bytes())
	y := new(btcec.FieldVal)
	y.SetByteSlice(pubKey.Y.Bytes())
	btcecPubKey := btcec.NewPublicKey(x, y)

	// Convert to compressed format
	return btcecPubKey.SerializeCompressed()
}

// EncryptShare encrypts a party's save data
func EncryptShare(key []byte, data *keygen.LocalPartySaveData) (*EncryptedShare, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: expected 32 bytes, got %d", len(key))
	}

	if data == nil {
		return nil, fmt.Errorf("share data cannot be nil")
	}

	// Marshal the save data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal share data: %w", err)
	}

	// TODO: Implement proper encryption using key
	return &EncryptedShare{
		EncryptedData: jsonData,
		ChainID:       "", // ChainID is set during network setup
	}, nil
}

// DecryptShare decrypts an encrypted share
func DecryptShare(key []byte, encryptedShare *EncryptedShare) (*keygen.LocalPartySaveData, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: expected 32 bytes, got %d", len(key))
	}

	if encryptedShare == nil {
		return nil, fmt.Errorf("encrypted share cannot be nil")
	}

	// TODO: Implement proper decryption using key
	var data keygen.LocalPartySaveData
	if err := json.Unmarshal(encryptedShare.EncryptedData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal share data: %w", err)
	}

	return &data, nil
}

// GenerateKey generates a key using TSS
func GenerateKey(ctx context.Context, preParams *keygen.LocalPreParams, chainID string) (*KeygenResult, error) {
	if preParams == nil {
		return nil, fmt.Errorf("preParams cannot be nil")
	}

	if chainID == "" {
		return nil, fmt.Errorf("chainID cannot be empty")
	}

	// Create parties
	p1, p2, err := createParties(preParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create parties: %v", err)
	}

	// Start key generation
	if err := p1.Start(); err != nil {
		return nil, fmt.Errorf("failed to start party 1: %v", err)
	}
	if err := p2.Start(); err != nil {
		return nil, fmt.Errorf("failed to start party 2: %v", err)
	}

	// Create channels for message passing
	outCh := make(chan tss.Message, TotalPartyCount)
	endCh := make(chan *keygen.LocalPartySaveData, TotalPartyCount)
	errCh := make(chan *tss.Error, TotalPartyCount)

	// Process messages until both parties are done
	var party1Data, party2Data *keygen.LocalPartySaveData
	for {
		select {
		case msg := <-outCh:
			dest := msg.GetTo()
			if dest == nil {
				continue
			}

			// Route message to appropriate party
			var err error
			wireBytes, _, wireErr := msg.WireBytes()
			if wireErr != nil {
				return nil, fmt.Errorf("failed to get wire bytes: %v", wireErr)
			}

			if dest[0].Index == p1.PartyID().Index {
				_, err = p1.UpdateFromBytes(wireBytes, msg.GetFrom(), msg.IsBroadcast())
			} else {
				_, err = p2.UpdateFromBytes(wireBytes, msg.GetFrom(), msg.IsBroadcast())
			}
			if err != nil {
				return nil, fmt.Errorf("failed to update party: %v", err)
			}

		case save := <-endCh:
			if save == nil {
				continue
			}

			// Store save data based on party ID
			if save.ShareID.Cmp(big.NewInt(1)) == 0 {
				party1Data = save
			} else {
				party2Data = save
			}

			// Both parties are done
			if party1Data != nil && party2Data != nil {
				// Extract public key from party1's data
				ecdsaPubKey := convertECPointToECDSA(party1Data.ECDSAPub)
				if ecdsaPubKey == nil {
					return nil, fmt.Errorf("failed to convert public key")
				}

				// Create result
				result := &KeygenResult{
					Party1Data:     party1Data,
					Party2Data:     party2Data,
					ChainID:        chainID,
					PublicKey:      ecdsaPubKey,
					PublicKeyBytes: serializePublicKey(ecdsaPubKey),
					SecurityLevel:  types.SecurityLevel_SECURITY_LEVEL_STANDARD, // Default to standard
				}

				return result, nil
			}

		case err := <-errCh:
			return nil, fmt.Errorf("key generation error: %v", err)

		case <-ctx.Done():
			return nil, fmt.Errorf("key generation cancelled")
		}
	}
}
