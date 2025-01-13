package tss

import (
    "crypto/rand"
    "math/big"
    "strconv"

    "github.com/bnb-chain/tss-lib/v2/tss"
)

const (
    // We use 2-party TSS setup (user device and chain)
    TotalPartyCount = 2
    Threshold       = 1 // t+1 must be <= n, where n is TotalPartyCount
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

// GeneratePartyID creates a new party ID for TSS operations
func GeneratePartyID(index int, partyCount int) (*tss.PartyID, error) {
    // Generate random key
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return nil, err
    }

    // Convert to big.Int for PartyID
    keyBigInt := new(big.Int).SetBytes(key)
    
    // Create a unique party ID with string index
    moniker := "P" + strconv.FormatInt(int64(index), 10)
    id := tss.NewPartyID(moniker, moniker, keyBigInt)
    id.Index = index
    return id, nil
}
