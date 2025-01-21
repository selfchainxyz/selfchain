package tss

import (
    "testing"

    "github.com/stretchr/testify/require"
)

func TestGeneratePartyID(t *testing.T) {
    id := GeneratePartyID(0)
    require.NotNil(t, id)
    require.Equal(t, "P0", id.Moniker)
    require.Equal(t, int64(1), id.KeyInt().Int64())
}
