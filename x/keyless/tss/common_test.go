package tss

import (
    "testing"

    "github.com/stretchr/testify/require"
)

func TestGeneratePartyID(t *testing.T) {
    id := GeneratePartyID(0)
    require.NotNil(t, id)
    require.Equal(t, "", id.Moniker) // We use empty moniker in the current implementation
    require.Equal(t, "P0", id.Id)    // The ID format is what we care about
    require.Equal(t, int64(1), id.KeyInt().Int64())
}
