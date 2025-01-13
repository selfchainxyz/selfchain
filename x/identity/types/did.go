package types

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidDIDID              = sdkerrors.Register(ModuleName, 1500, "invalid DID identifier")
	ErrInvalidDIDController      = sdkerrors.Register(ModuleName, 1501, "invalid DID controller")
	ErrInvalidVerificationMethod = sdkerrors.Register(ModuleName, 1502, "invalid verification method")
	ErrInvalidServiceEndpoint    = sdkerrors.Register(ModuleName, 1503, "invalid service endpoint")
	ErrUnauthorizedDIDUpdate     = sdkerrors.Register(ModuleName, 1506, "unauthorized DID update")
)

// ValidateBasic performs basic validation of DID document
func (d *DIDDocument) ValidateBasic() error {
	if d.Id == "" {
		return sdkerrors.Wrap(ErrInvalidDIDID, "DID identifier cannot be empty")
	}

	if len(d.Controller) == 0 {
		return sdkerrors.Wrap(ErrInvalidDIDController, "at least one controller is required")
	}

	// Validate verification methods
	for _, vm := range d.VerificationMethod {
		if err := vm.ValidateBasic(); err != nil {
			return err
		}
	}

	// Validate service endpoints
	for _, svc := range d.Service {
		if err := svc.ValidateBasic(); err != nil {
			return err
		}
	}

	// Check timestamps
	now := time.Date(2025, 1, 12, 15, 8, 51, 0, time.FixedZone("IST", 5*60*60+30*60))
	if d.Created == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "creation time must be set")
	}
	if d.Updated != nil && d.Updated.Before(*d.Created) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "update time cannot be before creation time")
	}
	if d.Updated != nil && d.Updated.After(now) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "update time cannot be in the future")
	}

	return nil
}

// ValidateBasic performs basic validation of verification method
func (v *VerificationMethod) ValidateBasic() error {
	if v.Id == "" {
		return sdkerrors.Wrap(ErrInvalidVerificationMethod, "verification method ID cannot be empty")
	}
	if v.Type == "" {
		return sdkerrors.Wrap(ErrInvalidVerificationMethod, "verification method type cannot be empty")
	}
	if v.Controller == "" {
		return sdkerrors.Wrap(ErrInvalidVerificationMethod, "controller cannot be empty")
	}
	if v.PublicKeyBase58 == "" {
		return sdkerrors.Wrap(ErrInvalidVerificationMethod, "public key cannot be empty")
	}
	return nil
}

// ValidateBasic performs basic validation of service endpoint
func (s *Service) ValidateBasic() error {
	if s.Id == "" {
		return sdkerrors.Wrap(ErrInvalidServiceEndpoint, "service ID cannot be empty")
	}
	if s.Type == "" {
		return sdkerrors.Wrap(ErrInvalidServiceEndpoint, "service type cannot be empty")
	}
	if s.ServiceEndpoint == "" {
		return sdkerrors.Wrap(ErrInvalidServiceEndpoint, "service endpoint cannot be empty")
	}
	return nil
}

// NewDIDDocument creates a new DID document
func NewDIDDocument(id string, controller []string) *DIDDocument {
	now := time.Date(2025, 1, 12, 15, 8, 51, 0, time.FixedZone("IST", 5*60*60+30*60))
	return &DIDDocument{
		Id:         id,
		Controller: controller,
		Created:    &now,
		Updated:    &now,
	}
}

// AddVerificationMethod adds a verification method to the DID document
func (d *DIDDocument) AddVerificationMethod(vm *VerificationMethod) error {
	if err := vm.ValidateBasic(); err != nil {
		return err
	}
	d.VerificationMethod = append(d.VerificationMethod, vm)
	return nil
}

// AddService adds a service endpoint to the DID document
func (d *DIDDocument) AddService(svc *Service) error {
	if err := svc.ValidateBasic(); err != nil {
		return err
	}
	d.Service = append(d.Service, svc)
	return nil
}

// HasController checks if the given DID is a controller of this DID document
func (d *DIDDocument) HasController(did string) bool {
	for _, controller := range d.Controller {
		if controller == did {
			return true
		}
	}
	return false
}

// HasVerificationMethod checks if the given verification method ID exists
func (d *DIDDocument) HasVerificationMethod(id string) bool {
	for _, vm := range d.VerificationMethod {
		if vm.Id == id {
			return true
		}
	}
	return false
}

// HasService checks if the given service ID exists
func (d *DIDDocument) HasService(id string) bool {
	for _, svc := range d.Service {
		if svc.Id == id {
			return true
		}
	}
	return false
}
