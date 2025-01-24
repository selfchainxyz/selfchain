package tss

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
	"time"

	"selfchain/x/keyless/crypto/signing/ecdsa"
	"selfchain/x/keyless/crypto/signing/format"
	"selfchain/x/keyless/crypto/signing/types"
	ktypes "selfchain/x/keyless/types"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/cosmos/cosmos-sdk/codec"
)

// Protocol implements the TSSProtocol interface for distributed key generation and signing
type Protocol struct {
	codec           codec.BinaryCodec
	sessions        sync.Map
	timeout         time.Duration
	cleanupInterval time.Duration
}

// NewProtocol creates a new TSS protocol instance
func NewProtocol(codec codec.BinaryCodec) ktypes.TSSProtocol {
	p := &Protocol{
		codec:           codec,
		timeout:         5 * time.Minute,
		cleanupInterval: 10 * time.Minute,
	}

	// Start cleanup goroutine
	go p.cleanupSessions()

	return p
}

// Session represents an active TSS session
type Session struct {
	ID            string
	Type          ktypes.SessionType
	PartyData     map[string]*ktypes.PartyData
	PublicKey     []byte
	SecurityLevel ktypes.SecurityLevel
	Status        ktypes.SessionStatus
	StartTime     time.Time
	mu            sync.RWMutex
	// Additional fields for signing
	Message   []byte
	WalletID  string
	Signature *format.SignatureResult
}

// cleanupSessions periodically removes completed or failed sessions
func (p *Protocol) cleanupSessions() {
	ticker := time.NewTicker(p.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		toDelete := make([]string, 0)

		// Find sessions to delete
		p.sessions.Range(func(key, value interface{}) bool {
			session := value.(*Session)
			session.mu.RLock()

			// Delete if:
			// 1. Session is completed or failed
			// 2. Session has timed out
			if session.Status == ktypes.SessionStatus_SESSION_STATUS_COMPLETED ||
				session.Status == ktypes.SessionStatus_SESSION_STATUS_FAILED ||
				now.Sub(session.StartTime) > p.timeout {
				toDelete = append(toDelete, key.(string))
			}

			session.mu.RUnlock()
			return true
		})

		// Delete sessions
		for _, id := range toDelete {
			p.sessions.Delete(id)
		}
	}
}

// GenerateKeyShares initiates a new key generation session
func (p *Protocol) GenerateKeyShares(ctx context.Context, req *ktypes.KeyGenRequest) (*ktypes.KeyGenResponse, error) {
	if req == nil {
		return nil, errors.New("nil request")
	}

	// Create new session
	session := &Session{
		ID:            generateSessionID(),
		Type:          ktypes.SessionType_SESSION_TYPE_KEYGEN,
		PartyData:     make(map[string]*ktypes.PartyData),
		SecurityLevel: req.SecurityLevel,
		Status:        ktypes.SessionStatus_SESSION_STATUS_PENDING,
		StartTime:     time.Now(),
		WalletID:      req.WalletId,
	}

	// Store session
	p.sessions.Store(session.ID, session)

	// Initialize key generation
	if err := p.initializeKeyGen(ctx, session, req); err != nil {
		session.Status = ktypes.SessionStatus_SESSION_STATUS_FAILED
		return nil, fmt.Errorf("failed to initialize key generation: %v", err)
	}

	session.Status = ktypes.SessionStatus_SESSION_STATUS_ACTIVE

	// Create response channel
	resultCh := make(chan *ktypes.KeyGenResponse, 1)
	errCh := make(chan error, 1)

	// Start monitoring goroutine
	go func() {
		for {
			session.mu.RLock()
			status := session.Status
			session.mu.RUnlock()

			switch status {
			case ktypes.SessionStatus_SESSION_STATUS_COMPLETED:
				now := time.Now()
				resultCh <- &ktypes.KeyGenResponse{
					WalletId:  req.WalletId,
					PublicKey: session.PublicKey,
					Metadata: &ktypes.KeyMetadata{
						CreatedAt:     now,
						LastRotated:   now,
						LastUsed:      now,
						UsageCount:    0,
						BackupStatus:  ktypes.BackupStatus_BACKUP_STATUS_NONE,
						SecurityLevel: req.SecurityLevel,
					},
				}
				return
			case ktypes.SessionStatus_SESSION_STATUS_FAILED:
				errCh <- errors.New("key generation failed")
				return
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Wait for completion or timeout
	select {
	case <-ctx.Done():
		session.Status = ktypes.SessionStatus_SESSION_STATUS_FAILED
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case result := <-resultCh:
		return result, nil
	case <-time.After(p.timeout):
		session.Status = ktypes.SessionStatus_SESSION_STATUS_FAILED
		return nil, errors.New("key generation timeout")
	}
}

// ProcessKeyGenRound processes a round of key generation
func (p *Protocol) ProcessKeyGenRound(ctx context.Context, sessionID string, partyData *ktypes.PartyData) error {
	if sessionID == "" || partyData == nil {
		return errors.New("invalid input parameters")
	}

	// Get session
	sessionObj, ok := p.sessions.Load(sessionID)
	if !ok {
		return errors.New("session not found")
	}
	session := sessionObj.(*Session)

	session.mu.Lock()
	defer session.mu.Unlock()

	// Check session status
	if session.Status != ktypes.SessionStatus_SESSION_STATUS_ACTIVE {
		return fmt.Errorf("invalid session status: %v", session.Status)
	}

	// Validate party data
	if err := validatePartyData(partyData); err != nil {
		return fmt.Errorf("invalid party data: %v", err)
	}

	// Store party data
	session.PartyData[partyData.PartyId] = partyData

	// Check if all parties have submitted data
	if len(session.PartyData) == getPartyCount(session.SecurityLevel) {
		if err := p.finalizeKeyGen(ctx, session); err != nil {
			session.Status = ktypes.SessionStatus_SESSION_STATUS_FAILED
			return fmt.Errorf("failed to finalize key generation: %v", err)
		}
		session.Status = ktypes.SessionStatus_SESSION_STATUS_COMPLETED
	}

	return nil
}

// InitiateSigning starts a new signing session
func (p *Protocol) InitiateSigning(ctx context.Context, msg []byte, walletID string) (*ktypes.SigningResponse, error) {
	if len(msg) == 0 || walletID == "" {
		return nil, errors.New("invalid input parameters")
	}

	// Create signing session
	session := &Session{
		ID:        generateSessionID(),
		Type:      ktypes.SessionType_SESSION_TYPE_SIGNING,
		PartyData: make(map[string]*ktypes.PartyData),
		Status:    ktypes.SessionStatus_SESSION_STATUS_PENDING,
		StartTime: time.Now(),
		Message:   msg,
		WalletID:  walletID,
	}

	// Store session
	p.sessions.Store(session.ID, session)

	// Initialize signing
	if err := p.initializeSigning(ctx, session, msg, walletID); err != nil {
		session.Status = ktypes.SessionStatus_SESSION_STATUS_FAILED
		return nil, fmt.Errorf("failed to initialize signing: %v", err)
	}

	session.Status = ktypes.SessionStatus_SESSION_STATUS_ACTIVE

	// Create response channel
	resultCh := make(chan *ktypes.SigningResponse, 1)
	errCh := make(chan error, 1)

	// Start monitoring goroutine
	go func() {
		for {
			session.mu.RLock()
			status := session.Status
			sig := session.Signature
			session.mu.RUnlock()

			switch status {
			case ktypes.SessionStatus_SESSION_STATUS_COMPLETED:
				if sig == nil {
					errCh <- errors.New("signature not generated")
					return
				}

				// Format signature based on chain requirements
				sigBytes, err := format.FormatCosmosSignature(sig)
				if err != nil {
					errCh <- fmt.Errorf("failed to format signature: %v", err)
					return
				}

				now := time.Now()
				resultCh <- &ktypes.SigningResponse{
					WalletId:  walletID,
					Signature: sigBytes,
					Metadata: &ktypes.SignatureMetadata{
						Timestamp: &now,
						ChainId:   "", // Set from request
						SignType:  ktypes.SignatureType_SIGNATURE_TYPE_ECDSA,
					},
				}
				return
			case ktypes.SessionStatus_SESSION_STATUS_FAILED:
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
		session.Status = ktypes.SessionStatus_SESSION_STATUS_FAILED
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case result := <-resultCh:
		return result, nil
	case <-time.After(p.timeout):
		session.Status = ktypes.SessionStatus_SESSION_STATUS_FAILED
		return nil, errors.New("signing timeout")
	}
}

// Helper functions

func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func validatePartyData(data *ktypes.PartyData) error {
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

func getPartyCount(level ktypes.SecurityLevel) int {
	switch level {
	case ktypes.SecurityLevel_SECURITY_LEVEL_STANDARD:
		return 2
	case ktypes.SecurityLevel_SECURITY_LEVEL_HIGH:
		return 3
	case ktypes.SecurityLevel_SECURITY_LEVEL_ENTERPRISE:
		return 5
	default:
		return 2
	}
}

func (p *Protocol) initializeKeyGen(ctx context.Context, session *Session, req *ktypes.KeyGenRequest) error {
	// Generate initial key pair for this party
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// Create party data with our public key
	partyData := &ktypes.PartyData{
		PartyId:   generateSessionID(), // Use random party ID
		PublicKey: privKey.PubKey().SerializeCompressed(),
		ChainId:   req.ChainId,
		Status:    "active",
	}

	// Store our party data
	session.PartyData[partyData.PartyId] = partyData

	return nil
}

func (p *Protocol) finalizeKeyGen(ctx context.Context, session *Session) error {
	// In a real implementation, we would:
	// 1. Verify all party signatures
	// 2. Combine public keys using threshold scheme
	// 3. Generate and distribute key shares

	// For now, we'll just combine the first public key
	for _, party := range session.PartyData {
		session.PublicKey = party.PublicKey
		break
	}

	return nil
}

func (p *Protocol) initializeSigning(ctx context.Context, session *Session, msg []byte, walletID string) error {
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
	signer := ecdsa.NewECDSASigner(privKey, privKey.PubKey())
	sig, err := signer.Sign(ctx, hash[:], types.ECDSA)
	if err != nil {
		return fmt.Errorf("failed to create test signature: %v", err)
	}

	// Store the signature
	session.mu.Lock()
	session.Signature = &format.SignatureResult{
		R:     sig.R,
		S:     sig.S,
		Bytes: sig.Bytes,
	}
	session.mu.Unlock()

	return nil
}

func combinePublicKeys(partyData map[string]*ktypes.PartyData) ([]byte, error) {
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
