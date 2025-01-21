package keeper

import (
	"context"
	"fmt"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"selfchain/x/identity/types"
)

// RateLimiter implements rate limiting for requests
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	window   time.Duration
	limit    int
}

func NewRateLimiter(window time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		window:   window,
		limit:    limit,
	}
}

func (k Keeper) ValidateBasic(ctx context.Context, msg sdk.Msg) error {
	switch msg := msg.(type) {
	case *types.MsgCreateDID:
		return k.validateCreateDID(ctx, msg)
	case *types.MsgUpdateDID:
		return k.validateUpdateDID(ctx, msg)
	case *types.MsgDeleteDID:
		return k.validateDeleteDID(ctx, msg)
	case *types.MsgLinkSocialIdentity:
		return k.validateLinkSocialIdentity(ctx, msg)
	case *types.MsgUnlinkSocialIdentity:
		return k.validateUnlinkSocialIdentity(ctx, msg)
	case *types.MsgAddMFA:
		return k.validateAddMFA(ctx, msg)
	case *types.MsgRemoveMFA:
		return k.validateRemoveMFA(ctx, msg)
	case *types.MsgVerifyMFA:
		return k.validateVerifyMFA(ctx, msg)
	default:
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized identity message type: %T", msg)
	}
}

func (k Keeper) validateCreateDID(ctx context.Context, msg *types.MsgCreateDID) error {
	if msg.Id == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "DID cannot be empty")
	}
	if len(msg.VerificationMethod) == 0 {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "DID must have at least one verification method")
	}
	return nil
}

func (k Keeper) validateUpdateDID(ctx context.Context, msg *types.MsgUpdateDID) error {
	if msg.Id == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "DID cannot be empty")
	}
	if len(msg.VerificationMethod) == 0 {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "DID must have at least one verification method")
	}
	return nil
}

func (k Keeper) validateDeleteDID(ctx context.Context, msg *types.MsgDeleteDID) error {
	if msg.Id == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "DID cannot be empty")
	}
	return nil
}

func (k Keeper) validateLinkSocialIdentity(ctx context.Context, msg *types.MsgLinkSocialIdentity) error {
	if msg.Creator == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "creator cannot be empty")
	}
	if msg.Provider == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "provider cannot be empty")
	}
	if msg.Token == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "token cannot be empty")
	}
	return nil
}

func (k Keeper) validateUnlinkSocialIdentity(ctx context.Context, msg *types.MsgUnlinkSocialIdentity) error {
	if msg.Creator == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "creator cannot be empty")
	}
	if msg.Provider == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "provider cannot be empty")
	}
	return nil
}

func (k Keeper) validateAddMFA(ctx context.Context, msg *types.MsgAddMFA) error {
	if msg.Creator == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "creator cannot be empty")
	}
	if msg.Did == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "DID cannot be empty")
	}
	if msg.Method == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "method cannot be empty")
	}
	if msg.Secret == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "secret cannot be empty")
	}
	return nil
}

func (k Keeper) validateRemoveMFA(ctx context.Context, msg *types.MsgRemoveMFA) error {
	if msg.Creator == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "creator cannot be empty")
	}
	if msg.Did == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "DID cannot be empty")
	}
	if msg.Method == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "method cannot be empty")
	}
	return nil
}

func (k Keeper) validateVerifyMFA(ctx context.Context, msg *types.MsgVerifyMFA) error {
	if msg.Creator == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "creator cannot be empty")
	}
	if msg.Did == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "DID cannot be empty")
	}
	if msg.Method == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "method cannot be empty")
	}
	if msg.Code == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "code cannot be empty")
	}
	return nil
}

func (k Keeper) RateLimitMiddleware(handler sdk.Handler) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		// Apply rate limiting based on message type and creator
		// TODO: Implement rate limiting logic
		return handler(ctx, msg)
	}
}

func (k Keeper) ValidateRequestMiddleware(handler sdk.Handler) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		if err := k.ValidateBasic(ctx, msg); err != nil {
			return nil, err
		}
		return handler(ctx, msg)
	}
}

func (k Keeper) LoggingMiddleware(next sdk.Handler) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		// Log the request
		k.Logger(ctx).Info(fmt.Sprintf("Processing message: %T", msg))

		// Call the next handler
		res, err := next(ctx, msg)

		// Log the result
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Error processing message: %v", err))
		} else {
			k.Logger(ctx).Info("Message processed successfully")
		}

		return res, err
	}
}

func (k Keeper) ChainMiddleware(middlewares ...func(sdk.Handler) sdk.Handler) func(sdk.Handler) sdk.Handler {
	return func(handler sdk.Handler) sdk.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}

func (k Keeper) SecurityMiddleware(ctx sdk.Context, action string, handler func() error) error {
	// Check if action is allowed
	if !k.isActionAllowed(ctx, action) {
		return sdkerrors.Wrap(types.ErrUnauthorized, "action not allowed")
	}

	// Apply rate limiting
	if err := k.checkRateLimit(ctx, action); err != nil {
		return err
	}

	// Log the action
	k.Logger(ctx).Info(fmt.Sprintf("Executing action: %s", action))

	// Execute the handler
	if err := handler(); err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("Error executing action: %v", err))
		return err
	}

	return nil
}

func (k Keeper) validateRequest(ctx sdk.Context, did string) error {
	// Check if DID exists
	if !k.HasDID(ctx, did) {
		return sdkerrors.Wrap(types.ErrDIDNotFound, "DID not found")
	}

	// Check if DID is active
	if !k.isDIDActive(ctx, did) {
		return sdkerrors.Wrap(types.ErrDIDInactive, "DID is inactive")
	}

	return nil
}

func (k Keeper) HasDID(ctx sdk.Context, did string) bool {
	return k.HasDIDDocument(ctx, did)
}

func (k Keeper) isDIDActive(ctx sdk.Context, did string) bool {
	// TODO: Implement DID status check
	return true
}

func (k Keeper) isActionAllowed(ctx sdk.Context, action string) bool {
	// TODO: Implement action authorization logic
	return true
}

func (k Keeper) checkRateLimit(ctx sdk.Context, action string) error {
	// TODO: Implement rate limiting logic
	return nil
}
