package keeper

import (
	"context"
	"fmt"
	"strings"
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
	case *types.MsgConfigureMFA:
		return k.validateConfigureMFA(ctx, msg)
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

func (k Keeper) validateConfigureMFA(ctx context.Context, msg *types.MsgConfigureMFA) error {
	if msg.Did == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "DID cannot be empty")
	}
	if len(msg.Methods) == 0 {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "MFA methods cannot be empty")
	}
	return nil
}

func (k Keeper) validateVerifyMFA(ctx context.Context, msg *types.MsgVerifyMFA) error {
	if msg.Did == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "DID cannot be empty")
	}
	if msg.Method == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "MFA method cannot be empty")
	}
	if msg.Code == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "verification code cannot be empty")
	}
	return nil
}

func (k Keeper) RateLimitMiddleware(handler sdk.Handler) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		// Get the DID from the message
		did := k.GetDIDFromMsg(msg)
		if did == "" {
			return handler(ctx, msg)
		}

		// Check rate limit
		if err := k.CheckRateLimit(ctx, did); err != nil {
			return nil, err
		}

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
		start := time.Now()

		// Call the next handler
		result, err := next(ctx, msg)

		// Log request details
		k.Logger(ctx).Info("Request processed",
			"msg_type", fmt.Sprintf("%T", msg),
			"sender", msg.GetSigners()[0].String(),
			"duration_ms", time.Since(start).Milliseconds(),
			"error", err,
		)

		return result, err
	}
}

func (k Keeper) ChainMiddleware(middlewares ...func(sdk.Handler) sdk.Handler) func(sdk.Handler) sdk.Handler {
	return func(next sdk.Handler) sdk.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

func (k Keeper) SecurityMiddleware(ctx sdk.Context, action string, handler func() error) error {
	// Check rate limit for this action
	if err := k.CheckRateLimit(ctx, action); err != nil {
		return fmt.Errorf("rate limit exceeded for action %s: %v", action, err)
	}

	// Log the request
	start := time.Now()
	err := handler()
	duration := time.Since(start)

	// Log request details
	k.Logger(ctx).Info("Security middleware",
		"action", action,
		"duration_ms", duration.Milliseconds(),
		"error", err,
	)

	return err
}

func (k Keeper) validateRequest(ctx sdk.Context, did string) error {
	// Check DID format
	if !strings.HasPrefix(did, "did:selfchain:") {
		return fmt.Errorf("invalid DID format: %s", did)
	}

	// Check if DID exists
	if !k.HasDID(ctx, did) {
		return fmt.Errorf("DID not found: %s", did)
	}

	return nil
}

func (k Keeper) HasDID(ctx sdk.Context, did string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(did))
}
