package tss

import (
	"context"
	"fmt"
	"math/big"

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
	// Create communication channels
	outCh1 := make(chan tss.Message, TotalPartyCount)
	outCh2 := make(chan tss.Message, TotalPartyCount)
	endCh := make(chan *common.SignatureData, TotalPartyCount)
	tssErrCh := make(chan *tss.Error, 2*TotalPartyCount)

	// Create party IDs with different IDs and moniker strings
	p1ID := tss.NewPartyID(party1Data.ShareID.String(), "P1", party1Data.ShareID)
	p2ID := tss.NewPartyID(party2Data.ShareID.String(), "P2", party2Data.ShareID)
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

	// Wait for completion or error
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-tssErrCh:
		return nil, fmt.Errorf("signing error: %w", err)
	case data := <-endCh:
		return &SignResult{
			R: new(big.Int).SetBytes(data.R),
			S: new(big.Int).SetBytes(data.S),
		}, nil
	}
}
