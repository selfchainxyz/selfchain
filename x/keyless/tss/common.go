package tss

import (
	"fmt"
	"math/big"

	"github.com/bnb-chain/tss-lib/v2/tss"
)

const (
    // TotalPartyCount is the total number of parties in the TSS protocol
    TotalPartyCount = 2
    // Threshold is the minimum number of parties required to sign
    Threshold = TotalPartyCount - 1
)

// SharedPartyUpdater is a helper function to update a party with a message
func SharedPartyUpdater(party tss.Party, msg tss.Message, errCh chan<- *tss.Error) {
    // do not send a message from this party back to itself
    if party.PartyID().Index == msg.GetFrom().Index {
        return
    }

    // Convert message to wire format
    wireBytes, _, err := msg.WireBytes()
    if err != nil {
        errCh <- party.WrapError(fmt.Errorf("failed to convert message to wire format: %w", err))
        return
    }

    // Parse wire message
    parsedMsg, err := tss.ParseWireMessage(wireBytes, msg.GetFrom(), msg.IsBroadcast())
    if err != nil {
        errCh <- party.WrapError(fmt.Errorf("failed to parse wire message: %w", err))
        return
    }

    // Validate message
    if parsedMsg.GetFrom() == nil {
        errCh <- party.WrapError(fmt.Errorf("message has nil sender"))
        return
    }

    if parsedMsg.GetTo() != nil {
        // Point-to-point message
        isForMe := false
        for _, to := range parsedMsg.GetTo() {
            if to.Index == party.PartyID().Index {
                isForMe = true
                break
            }
        }
        if !isForMe {
            return
        }
    }

    // Update party state with the message
    if _, err := party.Update(parsedMsg); err != nil {
        errCh <- party.WrapError(fmt.Errorf("failed to update party state: %w", err))
        return
    }
}

// GeneratePartyID creates a new party ID for TSS
func GeneratePartyID(index int) *tss.PartyID {
    // Convert index to big.Int for party ID
    key := new(big.Int).SetInt64(int64(index + 1))
    
    // Create a unique moniker for the party
    moniker := fmt.Sprintf("P%d", index)
    
    // Create party ID with empty metadata
    return tss.NewPartyID(moniker, "", key)
}

// GeneratePartyIDs creates a sorted list of party IDs
func GeneratePartyIDs(count int) ([]*tss.PartyID, error) {
    if count < 2 {
        return nil, fmt.Errorf("minimum 2 parties required, got %d", count)
    }

    var parties []*tss.PartyID
    for i := 0; i < count; i++ {
        id := GeneratePartyID(i)
        parties = append(parties, id)
    }

    return tss.SortPartyIDs(parties), nil
}
