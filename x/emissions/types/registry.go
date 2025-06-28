package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ModuleStatus represents the status of a registered module
type ModuleStatus string

const (
	// REGISTERED - Module is registered and eligible for governance funding
	ModuleStatusRegistered ModuleStatus = "REGISTERED"
	// PENDING - Module registration is pending validation (future use)
	ModuleStatusPending ModuleStatus = "PENDING"
	// DISABLED - Module is disabled and cannot receive funding
	ModuleStatusDisabled ModuleStatus = "DISABLED"
)

// RegisteredModule represents a module registered in the emissions registry
type RegisteredModule struct {
	// ModuleName is the unique identifier for the module
	ModuleName string `json:"module_name"`
	// Creator is the address that registered the module
	Creator string `json:"creator"`
	// Description provides context about the module's purpose
	Description string `json:"description"`
	// Status indicates the current state of the module registration
	Status ModuleStatus `json:"status"`
	// CreatedAt tracks when the module was registered
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt tracks the last modification time
	UpdatedAt time.Time `json:"updated_at"`
}

// NewRegisteredModule creates a new RegisteredModule with REGISTERED status
func NewRegisteredModule(moduleName, creator, description string) RegisteredModule {
	now := time.Now()
	return RegisteredModule{
		ModuleName:  moduleName,
		Creator:     creator,
		Description: description,
		Status:      ModuleStatusRegistered,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Validate performs basic validation on the RegisteredModule
func (rm RegisteredModule) Validate() error {
	if rm.ModuleName == "" {
		return ErrInvalidDestination.Wrap("module name cannot be empty")
	}

	if rm.Creator == "" {
		return ErrInvalidDestination.Wrap("creator cannot be empty")
	}

	if _, err := sdk.AccAddressFromBech32(rm.Creator); err != nil {
		return ErrInvalidDestination.Wrapf("invalid creator address: %s", err.Error())
	}

	if rm.Description == "" {
		return ErrInvalidDestination.Wrap("description cannot be empty")
	}

	if len(rm.Description) > 500 {
		return ErrInvalidDestination.Wrap("description too long (max 500 characters)")
	}

	// Validate status
	switch rm.Status {
	case ModuleStatusRegistered, ModuleStatusPending, ModuleStatusDisabled:
		// Valid status
	default:
		return ErrInvalidDestination.Wrapf("invalid module status: %s", rm.Status)
	}

	return nil
}

// IsActive returns true if the module is in an active state for funding eligibility
func (rm RegisteredModule) IsActive() bool {
	return rm.Status == ModuleStatusRegistered
}

// ValidateModuleName performs validation on a module name for registration
func ValidateModuleName(moduleName string) error {
	if moduleName == "" {
		return ErrInvalidDestination.Wrap("module name cannot be empty")
	}

	if len(moduleName) < 3 {
		return ErrInvalidDestination.Wrap("module name too short (minimum 3 characters)")
	}

	if len(moduleName) > 64 {
		return ErrInvalidDestination.Wrap("module name too long (maximum 64 characters)")
	}

	// Check for valid characters (alphanumeric, underscore, hyphen)
	for _, char := range moduleName {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return ErrInvalidDestination.Wrapf("invalid character in module name: %c", char)
		}
	}

	return nil
}
