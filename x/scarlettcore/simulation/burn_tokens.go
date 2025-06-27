package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"scarlett-core/x/scarlettcore/keeper"
	"scarlett-core/x/scarlettcore/types"
)

func SimulateMsgBurnTokens(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgBurnTokens{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handle the BurnTokens simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "BurnTokens simulation not implemented"), nil, nil
	}
}
