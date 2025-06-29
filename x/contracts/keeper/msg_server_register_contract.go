package keeper

import (
	"context"
	"strings"

	"scarlett-core/x/contracts/types"
	emissionstypes "scarlett-core/x/emissions/types"

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

	// Validate that the contract actually exists by checking it was previously deployed
	err := k.validateContractExists(ctx, msg.ContractAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "contract validation failed")
	}

	// Check if contract is already registered with emissions
	if k.emissionsKeeper.IsModuleRegistered(ctx, msg.ContractAddress) {
		return nil, errorsmod.Wrap(types.ErrInvalidInput, "contract is already registered for emissions")
	}

	// Store contract registration information
	err = k.storeContractRegistration(ctx, msg.ContractAddress, msg.Name, msg.Description, msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to store contract registration")
	}

	// Register contract with emissions keeper as a valid destination
	// Use the proper constructor from emissions types
	registeredModule := emissionstypes.NewRegisteredModule(
		msg.ContractAddress, // Use contract address as module name for emissions
		msg.Creator,
		msg.Description,
	)

	err = k.emissionsKeeper.SetRegisteredModule(ctx, registeredModule)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to register contract with emissions keeper")
	}

	// Emit comprehensive contract registration event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"contract_registered",
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("contract_address", msg.ContractAddress),
			sdk.NewAttribute("name", msg.Name),
			sdk.NewAttribute("description", msg.Description),
			sdk.NewAttribute("status", "REGISTERED"),
			sdk.NewAttribute("emissions_eligible", "true"),
		),
	)

	// Log successful registration
	sdkCtx.Logger().Info("✅ CONTRACT REGISTRATION SUCCESSFUL",
		"contract_address", msg.ContractAddress,
		"name", msg.Name,
		"creator", msg.Creator,
		"emissions_status", "REGISTERED")

	return &types.MsgRegisterContractResponse{}, nil
}

// validateContractExists checks if the contract was previously deployed via DeployContract
func (k msgServer) validateContractExists(ctx context.Context, contractAddress string) error {
	// Check if contract exists in our ContractInfo collection
	_, err := k.ContractInfo.Get(ctx, contractAddress)
	if err != nil {
		return errorsmod.Wrapf(types.ErrInvalidInput,
			"contract not found: %s (contract must be deployed before registration)",
			contractAddress)
	}

	return nil
}

// storeContractRegistration stores the contract registration information
func (k msgServer) storeContractRegistration(ctx context.Context, contractAddress, name, description, creator string) error {
	// Get existing contract info to update it with registration details
	contractInfo, err := k.ContractInfo.Get(ctx, contractAddress)
	if err != nil {
		return errorsmod.Wrap(err, "failed to get contract info")
	}

	// Update contract info with registration metadata
	// Note: In a more complete implementation, we might have a separate
	// ContractRegistration collection, but for V1 we'll store it in ContractInfo
	contractInfo.Label = name + " (Registered for Emissions)"

	// Store updated contract info
	err = k.ContractInfo.Set(ctx, contractAddress, contractInfo)
	if err != nil {
		return errorsmod.Wrap(err, "failed to update contract info with registration")
	}

	return nil
}
