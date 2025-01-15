package tss

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/bnb-chain/tss-lib/v2/common"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/signing"
	"github.com/bnb-chain/tss-lib/v2/tss"
)

// SignResult contains the result of signing
type SignResult struct {
	R *big.Int
	S *big.Int
}

// SignMessage signs a message using TSS
func SignMessage(ctx context.Context, msg []byte, party1Data, party2Data *keygen.LocalPartySaveData) (*SignResult, error) {
	// Create communication channels with sufficient buffer
	outCh1 := make(chan tss.Message, TotalPartyCount*4) // Increased buffer size
	outCh2 := make(chan tss.Message, TotalPartyCount*4)
	endCh := make(chan *common.SignatureData, TotalPartyCount)
	tssErrCh := make(chan *tss.Error, 2*TotalPartyCount)
	errChan := make(chan error, 2)

	// Create party IDs with consistent indices
	p1ID := tss.NewPartyID("P1", "", new(big.Int).SetInt64(1))
	p2ID := tss.NewPartyID("P2", "", new(big.Int).SetInt64(2))
	parties := tss.SortPartyIDs([]*tss.PartyID{p1ID, p2ID})

	// Create peer context
	peerCtx := tss.NewPeerContext(parties)
	params1 := tss.NewParameters(tss.S256(), peerCtx, p1ID, len(parties), Threshold)
	params2 := tss.NewParameters(tss.S256(), peerCtx, p2ID, len(parties), Threshold)

	// Convert message to big.Int
	msgBigInt := new(big.Int).SetBytes(msg)

	// Create signing parties
	p1 := signing.NewLocalParty(msgBigInt, params1, *party1Data, outCh1, endCh).(*signing.LocalParty)
	p2 := signing.NewLocalParty(msgBigInt, params2, *party2Data, outCh2, endCh).(*signing.LocalParty)

	// Start both parties with error handling
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := p1.Start(); err != nil {
			errChan <- fmt.Errorf("party 1 error: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := p2.Start(); err != nil {
			errChan <- fmt.Errorf("party 2 error: %w", err)
		}
	}()

	// Message routing for party 1
	go func() {
		defer close(outCh1)
		for msg := range outCh1 {
			dest := msg.GetTo()
			if dest == nil { // broadcast
				SharedPartyUpdater(p2, msg, tssErrCh)
			} else {
				for _, to := range dest {
					if to.Index == p2.PartyID().Index {
						SharedPartyUpdater(p2, msg, tssErrCh)
						break
					}
				}
			}
		}
	}()

	// Message routing for party 2
	go func() {
		defer close(outCh2)
		for msg := range outCh2 {
			dest := msg.GetTo()
			if dest == nil { // broadcast
				SharedPartyUpdater(p1, msg, tssErrCh)
			} else {
				for _, to := range dest {
					if to.Index == p1.PartyID().Index {
						SharedPartyUpdater(p1, msg, tssErrCh)
						break
					}
				}
			}
		}
	}()

	// Wait for completion or error with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errChan:
		return nil, err
	case err := <-tssErrCh:
		return nil, fmt.Errorf("signing error: %w", err)
	case data := <-endCh:
		// Wait for both parties to finish
		select {
		case <-done:
			return &SignResult{
				R: new(big.Int).SetBytes(data.R),
				S: new(big.Int).SetBytes(data.S),
			}, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
