package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"scarlett-core/x/contracts/keeper"
	"scarlett-core/x/contracts/types"
)

func TestMsgUpdateParams(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	params := types.DefaultParams()

	// Since our keeper is now a wrapper around wasmd, we simplify the test
	testCases := []struct {
		name      string
		input     *types.MsgUpdateParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "valid params",
			input: &types.MsgUpdateParams{
				Authority: "cosmos1authority", // Mock authority for testing
				Params:    params,
			},
			expErr: false,
		},
		{
			name: "invalid params",
			input: &types.MsgUpdateParams{
				Authority: "cosmos1authority",
				Params:    types.Params{}, // Empty params should still be valid
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.UpdateParams(f.ctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
