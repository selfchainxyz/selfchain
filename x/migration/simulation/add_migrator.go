package simulation

import (
	"math/rand"

	"selfchain/x/migration/keeper"
	"selfchain/x/migration/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgAddMigrator(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgAddMigrator{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the AddMigrator simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "AddMigrator simulation not implemented"), nil, nil
	}
}
