package testutil

import (
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"selfchain/app"
)

// Setup initializes a new App for testing
func Setup(t testing.TB, isCheckTx bool) *app.App {
	t.Helper()

	// Create app
	app := app.New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		t.TempDir(),
		0,
		app.MakeEncodingConfig(),
		EmptyAppOptions{},
		[]wasm.Option{},
	)

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			ChainId:         "test-chain",
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams(),
		},
	)

	if !isCheckTx {
		// Create a new block
		app.BeginBlock(abci.RequestBeginBlock{
			Header: tmproto.Header{
				ChainID: "test-chain",
				Height:  app.LastBlockHeight() + 1,
				AppHash: app.LastCommitID().Hash,
			},
		})
	}

	return app
}
