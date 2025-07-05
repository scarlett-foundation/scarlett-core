package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"scarlett-core/x/contracts/keeper"
	module "scarlett-core/x/contracts/module"
	"scarlett-core/x/contracts/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
}

// mockWasmKeeper is a minimal mock for testing
type mockWasmKeeper struct{}

func (m *mockWasmKeeper) GetContractInfo(ctx context.Context, contractAddress sdk.AccAddress) interface{} {
	return nil
}

func (m *mockWasmKeeper) GetCodeInfo(ctx context.Context, codeID uint64) interface{} {
	return nil
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	// For testing, we'll use nil for the wasmd keeper since we're not testing wasm functionality
	k := keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		log.NewNopLogger(),
		nil, // wasmKeeper - nil for testing
		nil, // authKeeper - nil for testing
		nil, // bankKeeper - nil for testing
		nil, // accountKeeper - nil for testing
	)

	return &fixture{
		ctx:          ctx,
		keeper:       k,
		addressCodec: addressCodec,
	}
}
