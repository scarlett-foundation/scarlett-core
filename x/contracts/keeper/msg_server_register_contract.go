package keeper

import (
	"context"
	"strings"

	"scarlett-core/x/contracts/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RegisterContract(ctx context.Context, msg *types.MsgRegisterContract) (*types.MsgRegisterContractResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Validate contract address format
	if _, err := k.addressCodec.StringToBytes(msg.ContractAddress); err != nil {
		return nil, errorsmod.Wrap(err, "invalid contract address")
	}

	// Validate contract address is scarlett bech32 format
	if !strings.HasPrefix(msg.ContractAddress, "scarlett1") {
		return nil, errorsmod.Wrap(types.ErrInvalidInput, "contract address must use scarlett bech32 format")
	}

	// Validate name is not empty
	if strings.TrimSpace(msg.Name) == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidInput, "contract name cannot be empty")
	}

	// Validate description is not empty
	if strings.TrimSpace(msg.Description) == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidInput, "contract description cannot be empty")
	}

	// TODO: In Step 3, we will validate that the contract actually exists
	// by checking it was previously deployed via DeployContract
	err := k.validateContractExists(ctx, msg.ContractAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "contract validation failed")
	}

	// Store contract registration information
	err = k.storeContractRegistration(ctx, msg.ContractAddress, msg.Name, msg.Description, msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to store contract registration")
	}

	// TODO: In Step 3, we will integrate with emissions keeper:
	// k.emissionsKeeper.RegisterDestination(ctx, msg.ContractAddress, msg.Name, msg.Description)

	// Emit contract registration event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"contract_registered",
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("contract_address", msg.ContractAddress),
			sdk.NewAttribute("name", msg.Name),
			sdk.NewAttribute("description", msg.Description),
		),
	)

	return &types.MsgRegisterContractResponse{}, nil
}

// validateContractExists checks if the contract was previously deployed
func (k msgServer) validateContractExists(ctx context.Context, contractAddress string) error {
	// TODO: In Phase 2, this will check actual contract storage
	// For now, we'll do basic validation that it looks like a contract address

	// Basic validation: contract addresses should be longer than regular addresses
	if len(contractAddress) < 45 { // scarlett1 + 38 chars
		return errorsmod.Wrap(types.ErrInvalidInput, "contract address format invalid")
	}

	// TODO: Check that contract was deployed via our DeployContract message
	// This will be implemented when we add contract storage in Phase 2

	return nil
}

// storeContractRegistration stores the contract registration information
func (k msgServer) storeContractRegistration(ctx context.Context, contractAddress, name, description, creator string) error {
	// TODO: Store contract registration in keeper state
	// This will include:
	// - Contract address -> registration info mapping
	// - Registration status
	// - Creator information
	// - Metadata (name, description)

	// For now, we'll just validate the inputs are properly formatted
	return nil
}
