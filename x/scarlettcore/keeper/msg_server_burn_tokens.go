package keeper

import (
	"context"
	"strconv"

	"scarlett-core/x/scarlettcore/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) BurnTokens(ctx context.Context, msg *types.MsgBurnTokens) (*types.MsgBurnTokensResponse, error) {
	// 1. Validate burner address
	burnerAddr, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid burner address")
	}

	// 2. Parse and validate amount
	amount, err := strconv.ParseUint(msg.Amount, 10, 64)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "amount must be a positive integer")
	}
	if amount == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "amount must be greater than 0")
	}

	// 3. Validate denomination
	if msg.Denom == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidDenom, "denomination cannot be empty")
	}

	// 4. Create coin from amount and denom
	burnCoin := sdk.NewCoin(msg.Denom, math.NewInt(int64(amount)))
	burnCoins := sdk.NewCoins(burnCoin)

	// 5. Check user balance
	balance := k.bankKeeper.GetBalance(ctx, burnerAddr, msg.Denom)
	if balance.Amount.LT(burnCoin.Amount) {
		return nil, errorsmod.Wrapf(types.ErrInsufficientBalance,
			"insufficient balance: have %s, need %s", balance.Amount.String(), burnCoin.Amount.String())
	}

	// 6. Transfer tokens from user to module account (Transfer-then-Burn pattern)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, burnerAddr, types.ModuleName, burnCoins); err != nil {
		return nil, errorsmod.Wrap(err, "failed to transfer tokens to module account")
	}

	// 7. Burn tokens from module account
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins); err != nil {
		return nil, errorsmod.Wrap(err, "failed to burn tokens from module account")
	}

	// 8. Emit burn event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBurnTokens,
			sdk.NewAttribute(types.AttributeKeyBurner, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
		),
	)

	return &types.MsgBurnTokensResponse{}, nil
}
