package tss

import (
    "testing"

    "github.com/stretchr/testify/require"
)

func TestGeneratePartyID(t *testing.T) {
    tests := []struct {
        name    string
        index   int
        wantErr bool
    }{
        {
            name:    "valid party ID generation",
            index:   0,
            wantErr: false,
        },
        {
            name:    "valid party ID generation with index 1",
            index:   1,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            partyID, err := GeneratePartyID(tt.index)
            if tt.wantErr {
                require.Error(t, err)
                return
            }

            require.NoError(t, err)
            require.NotNil(t, partyID)
            require.Equal(t, tt.index, partyID.Index)
            require.Equal(t, "P"+string(rune('0'+tt.index)), partyID.Moniker)
            require.NotNil(t, partyID.Key)
            require.NotEmpty(t, partyID.Key) // Check if key is not empty
        })
    }
}
