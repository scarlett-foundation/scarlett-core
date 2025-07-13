package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"scarlett-core/x/proofofdegen/keeper"
	"scarlett-core/x/proofofdegen/types"
)

func createNCampaign(keeper keeper.Keeper, ctx context.Context, n int) []types.Campaign {
	items := make([]types.Campaign, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)
		items[i].Name = strconv.Itoa(i)
		items[i].Active = false
		items[i].TotalAllocation = uint64(i)
		_ = keeper.Campaign.Set(ctx, items[i].Index, items[i])
	}
	return items
}

func TestCampaignQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNCampaign(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetCampaignRequest
		response *types.QueryGetCampaignResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetCampaignRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetCampaignResponse{Campaign: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetCampaignRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetCampaignResponse{Campaign: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetCampaignRequest{
				Index: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetCampaign(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestCampaignQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNCampaign(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllCampaignRequest {
		return &types.QueryAllCampaignRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListCampaign(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Campaign), step)
			require.Subset(t, msgs, resp.Campaign)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListCampaign(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Campaign), step)
			require.Subset(t, msgs, resp.Campaign)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListCampaign(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Campaign)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListCampaign(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
