package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"selfchain/x/migration/keeper"
	"selfchain/x/migration/types"
)

func SimulateMsgUpdateConfig(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgUpdateConfig{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the UpdateConfig simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "UpdateConfig simulation not implemented"), nil, nil
	}
}
