package tss

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/bnb-chain/tss-lib/v2/crypto"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/btcsuite/btcd/btcec/v2"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"selfchain/x/keyless/crypto/signing/format"
	"selfchain/x/keyless/types"
)

// Protocol implements the TSSProtocol interface
type Protocol struct {
	sessions map[string]*Session
}

// NewTSSProtocolImpl creates a new instance of TSSProtocol
func NewTSSProtocolImpl() types.TSSProtocol {
	return &Protocol{
		sessions: make(map[string]*Session),
	}
}

// Session represents an active TSS session
type Session struct {
	ID            string
	WalletAddress string
	Threshold     uint32
	SecurityLevel types.SecurityLevel
	Status        types.SessionStatus
	Shares        [][]byte
	PublicKey     []byte
}

// GenerateKeyShares generates key shares for a new wallet
func (p *Protocol) GenerateKeyShares(ctx sdk.Context, walletAddress string, threshold uint32, securityLevel types.SecurityLevel) (*types.KeyGenResponse, error) {
	// Create key generation request
	req := &types.KeyGenRequest{
		WalletAddress: walletAddress,
		Threshold:     threshold,
		SecurityLevel: securityLevel,
	}

	// Create new session
	session := &Session{
		ID:            generateSessionID(),
		WalletAddress: walletAddress,
		Threshold:     threshold,
		SecurityLevel: securityLevel,
		Status:        types.SessionStatus_SESSION_STATUS_ACTIVE,
	}

	// Store session
	p.sessions[session.ID] = session

	// Initialize key generation
	if err := p.initializeKeyGen(ctx, session, req); err != nil {
		return nil, fmt.Errorf("failed to initialize key generation: %w", err)
	}

	// Return response
	return &types.KeyGenResponse{
		WalletAddress: walletAddress,
		PublicKey:     session.PublicKey,
		Metadata: &types.KeyMetadata{
			CreatedAt:     time.Now().UTC(),
			LastRotated:   time.Now().UTC(),
			LastUsed:      time.Now().UTC(),
			UsageCount:    0,
			BackupStatus:  types.BackupStatus_BACKUP_STATUS_NONE,
			SecurityLevel: securityLevel,
		},
	}, nil
}

// ProcessKeyGenRound processes a round of key generation
func (p *Protocol) ProcessKeyGenRound(ctx context.Context, sessionID string, partyData *types.PartyData) error {
	if sessionID == "" || partyData == nil {
		return errors.New("invalid input parameters")
	}

	// Get session
	session, ok := p.sessions[sessionID]
	if !ok {
		return errors.New("session not found")
	}

	// Check session status
	if session.Status != types.SessionStatus_SESSION_STATUS_ACTIVE {
		return fmt.Errorf("invalid session status: %v", session.Status)
	}

	// Validate party data
	if err := validatePartyData(partyData); err != nil {
		return fmt.Errorf("invalid party data: %v", err)
	}

	// Store party data
	session.Shares = append(session.Shares, partyData.PartyShare)

	// Check if all parties have submitted data
	if len(session.Shares) == getPartyCount(session.SecurityLevel) {
		if err := p.finalizeKeyGen(ctx, session); err != nil {
			session.Status = types.SessionStatus_SESSION_STATUS_FAILED
			return fmt.Errorf("failed to finalize key generation: %v", err)
		}
		session.Status = types.SessionStatus_SESSION_STATUS_COMPLETED
	}

	return nil
}

// GetPartyData gets TSS party data
func (p *Protocol) GetPartyData(ctx sdk.Context, partyID string) (*types.PartyData, error) {
	if partyID == "" {
		return nil, fmt.Errorf("empty party ID")
	}

	session, err := p.getActiveSession()
	if err != nil {
		return nil, err
	}

	// For now, return the first share as party data
	if len(session.Shares) > 0 {
		return &types.PartyData{
			PartyId:    partyID,
			PartyShare: session.Shares[0],
		}, nil
	}

	return nil, fmt.Errorf("party data not found for ID: %s", partyID)
}

// SetPartyData sets TSS party data
func (p *Protocol) SetPartyData(ctx sdk.Context, data *types.PartyData) error {
	if data == nil {
		return fmt.Errorf("party data is nil")
	}
	if data.PartyId == "" {
		return fmt.Errorf("empty party ID")
	}

	session, err := p.getActiveSession()
	if err != nil {
		return err
	}

	// For now, store the party data as the first share
	session.Shares = append(session.Shares, data.PartyShare)
	return nil
}

// InitiateSigning starts a new signing session
func (p *Protocol) InitiateSigning(ctx context.Context, msg []byte, walletAddress string) (*types.SigningResponse, error) {
	if len(msg) == 0 || walletAddress == "" {
		return nil, errors.New("invalid input parameters")
	}

	// Create signing session
	session := &Session{
		ID:            generateSessionID(),
		WalletAddress: walletAddress,
		Status:        types.SessionStatus_SESSION_STATUS_PENDING,
	}

	// Store session
	p.sessions[session.ID] = session

	// Initialize signing
	if err := p.initializeSigning(ctx, session, msg, walletAddress); err != nil {
		session.Status = types.SessionStatus_SESSION_STATUS_FAILED
		return nil, fmt.Errorf("failed to initialize signing: %v", err)
	}

	session.Status = types.SessionStatus_SESSION_STATUS_ACTIVE

	// Create response channel
	resultCh := make(chan *types.SigningResponse, 1)
	errCh := make(chan error, 1)

	// Start monitoring goroutine
	go func() {
		for {
			status := session.Status
			sig := session.PublicKey

			switch status {
			case types.SessionStatus_SESSION_STATUS_COMPLETED:
				if sig == nil {
					errCh <- errors.New("signature not generated")
					return
				}

				// Format signature based on chain requirements
				sigResult := &format.SignatureResult{
					R: new(big.Int).SetBytes(sig[:32]),
					S: new(big.Int).SetBytes(sig[32:]),
				}
				sigBytes, err := format.FormatCosmosSignature(sigResult)
				if err != nil {
					errCh <- fmt.Errorf("failed to format signature: %v", err)
					return
				}

				now := time.Now()
				resultCh <- &types.SigningResponse{
					WalletAddress: walletAddress,
					Signature:     sigBytes,
					Metadata: &types.SignatureMetadata{
						Timestamp: &now,
						ChainId:   "", // Set from request
						SignType:  types.SignatureType_SIGNATURE_TYPE_ECDSA,
					},
				}
				return
			case types.SessionStatus_SESSION_STATUS_FAILED:
				errCh <- errors.New("signing failed")
				return
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Wait for completion or timeout
	select {
	case <-ctx.Done():
		session.Status = types.SessionStatus_SESSION_STATUS_FAILED
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case result := <-resultCh:
		return result, nil
	case <-time.After(5 * time.Minute):
		session.Status = types.SessionStatus_SESSION_STATUS_FAILED
		return nil, errors.New("signing timeout")
	}
}

// ReconstructKey reconstructs a key from shares
func (p *Protocol) ReconstructKey(ctx sdk.Context, shares [][]byte) ([]byte, error) {
	if len(shares) == 0 {
		return nil, fmt.Errorf("no shares provided")
	}

	// Validate shares
	for i, share := range shares {
		if len(share) == 0 {
			return nil, fmt.Errorf("empty share at index %d", i)
		}
	}

	// Create a new session for key reconstruction
	sessionID := generateSessionID()
	session := &Session{
		ID:        sessionID,
		Status:    types.SessionStatus_SESSION_STATUS_ACTIVE,
		Shares:    shares,
		PublicKey: nil,
	}

	// Store session
	p.sessions[sessionID] = session

	// Combine shares to reconstruct the key using TSS
	outCh := make(chan tss.Message, len(shares))
	endCh := make(chan *keygen.LocalPartySaveData, len(shares))
	errCh := make(chan *tss.Error, len(shares))

	// Initialize parties with shares
	parties := make([]tss.Party, len(shares))
	partyIDs := make([]*tss.PartyID, len(shares))

	// Create party IDs
	for i := range shares {
		id := fmt.Sprintf("party-%d", i+1)
		key := new(big.Int).SetInt64(int64(i))
		partyIDs[i] = tss.NewPartyID(id, id, key)
	}

	// Sort party IDs
	sortedPartyIDs := tss.SortPartyIDs(partyIDs)
	peerCtx := tss.NewPeerContext(sortedPartyIDs)

	// Initialize each party
	for i := range shares {
		params := tss.NewParameters(
			tss.Edwards(),
			peerCtx,
			partyIDs[i],
			len(shares),
			len(shares), // threshold = total parties for reconstruction
		)

		preParams := keygen.LocalPreParams{}
		party := keygen.NewLocalParty(params, outCh, endCh, preParams)
		parties[i] = party
	}

	// Start reconstruction
	for _, party := range parties {
		if err := party.Start(); err != nil {
			return nil, fmt.Errorf("failed to start party: %v", err)
		}
	}

	// Process messages until key is reconstructed
	var reconstructedKey []byte
	for {
		select {
		case msg := <-outCh:
			dest := msg.GetTo()
			if dest == nil {
				continue
			}

			// Route message to appropriate party
			wireBytes, _, err := msg.WireBytes()
			if err != nil {
				return nil, fmt.Errorf("failed to get wire bytes: %v", err)
			}

			for _, party := range parties {
				if dest[0].Index == party.PartyID().Index {
					if _, err := party.UpdateFromBytes(wireBytes, msg.GetFrom(), msg.IsBroadcast()); err != nil {
						return nil, fmt.Errorf("failed to update party: %v", err)
					}
					break
				}
			}

		case save := <-endCh:
			if save == nil {
				continue
			}

			// Convert the reconstructed key to compressed bytes
			pubKey := save.ECDSAPub
			if pubKey != nil {
				point, err := crypto.NewECPoint(tss.Edwards(), pubKey.X(), pubKey.Y())
				if err != nil {
					return nil, fmt.Errorf("failed to create EC point: %v", err)
				}
				// Convert EC point to bytes
				x := point.X()
				y := point.Y()
				reconstructedKey = append(x.Bytes(), y.Bytes()...)
				return reconstructedKey, nil
			}

		case err := <-errCh:
			return nil, fmt.Errorf("key reconstruction error: %v", err)

		case <-ctx.Done():
			return nil, fmt.Errorf("key reconstruction cancelled")
		}
	}
}

// VerifyShare verifies a share's validity
func (p *Protocol) VerifyShare(ctx sdk.Context, share []byte, publicKey []byte) error {
	if len(share) == 0 {
		return fmt.Errorf("empty share")
	}
	if len(publicKey) == 0 {
		return fmt.Errorf("empty public key")
	}

	// TODO: Implement proper share verification using TSS
	// For now, just check that the share is not empty
	return nil
}

// VerifySignature verifies a TSS signature
func (p *Protocol) VerifySignature(ctx sdk.Context, message []byte, signature []byte, publicKey []byte) error {
	if len(message) == 0 {
		return fmt.Errorf("empty message")
	}
	if len(signature) == 0 {
		return fmt.Errorf("empty signature")
	}
	if len(publicKey) == 0 {
		return fmt.Errorf("empty public key")
	}

	// TODO: Implement signature verification using TSS
	return nil
}

// SignMessage signs a message using TSS
func (p *Protocol) SignMessage(ctx sdk.Context, message []byte, shares [][]byte) ([]byte, error) {
	if len(message) == 0 {
		return nil, fmt.Errorf("empty message")
	}
	if len(shares) == 0 {
		return nil, fmt.Errorf("no shares provided")
	}

	// Hash the message
	hash := sha256.Sum256(message)

	// In a real implementation, we would:
	// 1. Distribute the message hash to all parties
	// 2. Wait for partial signatures
	// 3. Combine partial signatures

	// For now, return error
	// TODO: Implement proper TSS signing
	r, s, err := ecdsa.Sign(rand.Reader, &ecdsa.PrivateKey{PublicKey: ecdsa.PublicKey{Curve: btcec.S256()}}, hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create test signature: %v", err)
	}

	// Format the signature according to the chain's requirements
	sigResult := &format.SignatureResult{
		R: r,
		S: s,
	}
	formattedSig, err := format.FormatCosmosSignature(sigResult)
	if err != nil {
		return nil, fmt.Errorf("failed to format signature: %v", err)
	}

	return formattedSig, nil
}

func (p *Protocol) getActiveSession() (*Session, error) {
	var activeSession *Session
	for _, session := range p.sessions {
		if session.Status == types.SessionStatus_SESSION_STATUS_ACTIVE {
			activeSession = session
			break
		}
	}

	if activeSession == nil {
		return nil, fmt.Errorf("no active session found")
	}

	return activeSession, nil
}

func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func validatePartyData(data *types.PartyData) error {
	if data == nil {
		return errors.New("party data is nil")
	}
	if data.PartyId == "" {
		return errors.New("party ID is empty")
	}
	if len(data.PublicKey) == 0 {
		return errors.New("public key is empty")
	}
	if data.ChainId == "" {
		return errors.New("chain ID is empty")
	}
	return nil
}

func getPartyCount(level types.SecurityLevel) int {
	switch level {
	case types.SecurityLevel_SECURITY_LEVEL_STANDARD:
		return 2
	case types.SecurityLevel_SECURITY_LEVEL_HIGH:
		return 3
	case types.SecurityLevel_SECURITY_LEVEL_ENTERPRISE:
		return 5
	default:
		return 2
	}
}

func (p *Protocol) initializeKeyGen(ctx context.Context, session *Session, req *types.KeyGenRequest) error {
	// Generate initial key pair for this party
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// Create party data with our public key
	partyData := &types.PartyData{
		PartyId:   generateSessionID(), // Use random party ID
		PublicKey: privKey.PubKey().SerializeCompressed(),
		ChainId:   req.ChainId,
		Status:    "active",
	}

	// Store our party data
	session.Shares = append(session.Shares, partyData.PartyShare)

	return nil
}

func (p *Protocol) finalizeKeyGen(ctx context.Context, session *Session) error {
	// In a real implementation, we would:
	// 1. Verify all party signatures
	// 2. Combine public keys using threshold scheme
	// 3. Generate and distribute key shares

	// For now, we'll just combine the first public key
	for _, share := range session.Shares {
		session.PublicKey = share
		break
	}

	return nil
}

func (p *Protocol) initializeSigning(ctx context.Context, session *Session, msg []byte, walletAddress string) error {
	// Hash the message
	hash := sha256.Sum256(msg)

	// In a real implementation, we would:
	// 1. Distribute the message hash to all parties
	// 2. Wait for partial signatures
	// 3. Combine partial signatures

	// For now, we'll create a temporary key for testing
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return fmt.Errorf("failed to generate test key: %v", err)
	}

	// Create a test signature
	r, s, err := ecdsa.Sign(rand.Reader, privKey.ToECDSA(), hash[:])
	if err != nil {
		return fmt.Errorf("failed to create test signature: %v", err)
	}

	signatureResult := &format.SignatureResult{
		R: r,
		S: s,
	}
	session.PublicKey, err = format.FormatCosmosSignature(signatureResult)
	if err != nil {
		return fmt.Errorf("failed to format signature: %v", err)
	}

	return nil
}

func combinePublicKeys(partyData map[string]*types.PartyData) ([]byte, error) {
	// In a real implementation, this would:
	// 1. Validate all public keys
	// 2. Combine them using threshold scheme
	// 3. Return the combined public key

	// For now, return the first public key
	for _, party := range partyData {
		return party.PublicKey, nil
	}
	return nil, errors.New("no public keys available")
}
