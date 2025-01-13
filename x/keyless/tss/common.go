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
    if party.PartyID() == msg.GetFrom() {
        return
    }
    bz, _, err := msg.WireBytes()
    if err != nil {
        errCh <- party.WrapError(err)
        return
    }
    pMsg, err := tss.ParseWireMessage(bz, msg.GetFrom(), msg.IsBroadcast())
    if err != nil {
        errCh <- party.WrapError(err)
        return
    }
    if _, err := party.Update(pMsg); err != nil {
        errCh <- err
    }
}

// GeneratePartyID creates a new party ID for TSS
func GeneratePartyID(index int) (*tss.PartyID, error) {
    // Convert index to big.Int for party ID
    key := big.NewInt(int64(index))
    
    // Create a unique moniker for the party
    moniker := fmt.Sprintf("P%d", index)
    
    // Create and sort party ID
    id := tss.NewPartyID(key.String(), moniker, key)
    return id, nil
}

// GeneratePartyIDs creates a sorted list of party IDs
func GeneratePartyIDs(count int) ([]*tss.PartyID, error) {
    var parties []*tss.PartyID
    
    for i := 0; i < count; i++ {
        id, err := GeneratePartyID(i)
        if err != nil {
            return nil, fmt.Errorf("failed to generate party ID %d: %w", i, err)
        }
        parties = append(parties, id)
    }
    
    // Sort party IDs
    sorted := tss.SortPartyIDs(parties)
    return sorted, nil
}
