package tss

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"selfchain/x/keyless/crypto"
)

// KeygenResult contains the result of key generation
type KeygenResult struct {
	Party1Data *keygen.LocalPartySaveData
	Party2Data *keygen.LocalPartySaveData
	ChainID    string
}

// EncryptedShare represents an encrypted key share
type EncryptedShare struct {
	EncryptedData string
	ChainID      string
}

// EncryptShare encrypts a party's save data
func EncryptShare(key crypto.EncryptionKey, data *keygen.LocalPartySaveData) (*EncryptedShare, error) {
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

	// Encrypt the JSON data
	encryptedData, err := crypto.Encrypt(key, jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt share data: %w", err)
	}

	return &EncryptedShare{
		EncryptedData: encryptedData,
	}, nil
}

// DecryptShare decrypts an encrypted share
func DecryptShare(key crypto.EncryptionKey, encryptedShare *EncryptedShare) (*keygen.LocalPartySaveData, error) {
	// Decrypt the data
	decryptedData, err := crypto.Decrypt(key, encryptedShare.EncryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt share data: %w", err)
	}

	// Unmarshal the JSON data
	var saveData keygen.LocalPartySaveData
	if err := json.Unmarshal(decryptedData, &saveData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal share data: %w", err)
	}

	return &saveData, nil
}

// GenerateKey generates a key using TSS
func GenerateKey(ctx context.Context, preParams *keygen.LocalPreParams, chainID string) (*KeygenResult, error) {
	// Create separate communication channels for each party
	outCh1 := make(chan tss.Message, TotalPartyCount)
	outCh2 := make(chan tss.Message, TotalPartyCount)
	endCh := make(chan *keygen.LocalPartySaveData, TotalPartyCount)
	tssErrCh := make(chan *tss.Error, 2*TotalPartyCount)

	// Create party IDs with different IDs and moniker strings
	p1ID := tss.NewPartyID("party1", "P1", big.NewInt(1))
	p2ID := tss.NewPartyID("party2", "P2", big.NewInt(2))
	parties := tss.SortPartyIDs([]*tss.PartyID{p1ID, p2ID})

	// Create peer context
	peerCtx := tss.NewPeerContext(parties)
	params1 := tss.NewParameters(tss.S256(), peerCtx, p1ID, len(parties), Threshold)
	params2 := tss.NewParameters(tss.S256(), peerCtx, p2ID, len(parties), Threshold)

	// Generate separate pre-parameters for each party
	preParams2, err := keygen.GeneratePreParams(time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate pre-parameters for party 2: %w", err)
	}

	// Create keygen parties with unique pre-parameters
	p1 := keygen.NewLocalParty(params1, outCh1, endCh, *preParams).(*keygen.LocalParty)
	p2 := keygen.NewLocalParty(params2, outCh2, endCh, *preParams2).(*keygen.LocalParty)

	// Start both parties
	go func() {
		if err := p1.Start(); err != nil {
			tssErrCh <- err
			return
		}
	}()

	go func() {
		if err := p2.Start(); err != nil {
			tssErrCh <- err
			return
		}
	}()

	// Message routing for party 1
	go func() {
		for msg := range outCh1 {
			dest := msg.GetTo()
			if dest == nil { // broadcast
				SharedPartyUpdater(p2, msg, tssErrCh)
			} else if dest[0].Index == p2.PartyID().Index {
				SharedPartyUpdater(p2, msg, tssErrCh)
			}
		}
	}()

	// Message routing for party 2
	go func() {
		for msg := range outCh2 {
			dest := msg.GetTo()
			if dest == nil { // broadcast
				SharedPartyUpdater(p1, msg, tssErrCh)
			} else if dest[0].Index == p1.PartyID().Index {
				SharedPartyUpdater(p1, msg, tssErrCh)
			}
		}
	}()

	// Wait for both parties to complete
	var party1Data, party2Data *keygen.LocalPartySaveData
	for i := 0; i < 2; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-tssErrCh:
			return nil, fmt.Errorf("keygen error: %w", err)
		case data := <-endCh:
			if data.ShareID.Cmp(big.NewInt(1)) == 0 {
				party1Data = data
			} else {
				party2Data = data
			}
		}
	}

	if party1Data == nil || party2Data == nil {
		return nil, fmt.Errorf("failed to receive data from both parties")
	}

	return &KeygenResult{
		Party1Data: party1Data,
		Party2Data: party2Data,
		ChainID:    chainID,
	}, nil
}
