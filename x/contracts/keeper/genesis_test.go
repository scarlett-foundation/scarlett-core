package keeper_test

import (
	"testing"

	"scarlett-core/x/contracts/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got := f.keeper.ExportGenesis(f.ctx)
	require.NotNil(t, got)

	require.Equal(t, genesisState.Params, got.Params)
}
