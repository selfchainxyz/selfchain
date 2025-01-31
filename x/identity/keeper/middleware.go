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

// validateCreateDID validates MsgCreateDID
func (k Keeper) validateCreateDID(ctx context.Context, msg *types.MsgCreateDID) error {
	if msg.Creator == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "creator cannot be empty")
	}
	return nil
}

// validateUpdateDID validates MsgUpdateDID
func (k Keeper) validateUpdateDID(ctx context.Context, msg *types.MsgUpdateDID) error {
	if msg.Creator == "" || msg.Id == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "creator and DID cannot be empty")
	}
	return nil
}

// validateDeleteDID validates MsgDeleteDID
func (k Keeper) validateDeleteDID(ctx context.Context, msg *types.MsgDeleteDID) error {
	if msg.Creator == "" || msg.Id == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "creator and DID cannot be empty")
	}
	return nil
}

// validateLinkSocialIdentity validates MsgLinkSocialIdentity
func (k Keeper) validateLinkSocialIdentity(ctx context.Context, msg *types.MsgLinkSocialIdentity) error {
	if msg.Creator == "" || msg.Provider == "" || msg.Token == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "creator, provider and token cannot be empty")
	}
	return nil
}

// validateUnlinkSocialIdentity validates MsgUnlinkSocialIdentity
func (k Keeper) validateUnlinkSocialIdentity(ctx context.Context, msg *types.MsgUnlinkSocialIdentity) error {
	if msg.Creator == "" || msg.Provider == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "creator and provider cannot be empty")
	}
	return nil
}

// validateAddMFA validates MsgAddMFA
func (k Keeper) validateAddMFA(ctx context.Context, msg *types.MsgAddMFA) error {
	if msg.Creator == "" || msg.Type() == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "creator and MFA type cannot be empty")
	}
	return nil
}

// validateRemoveMFA validates MsgRemoveMFA
func (k Keeper) validateRemoveMFA(ctx context.Context, msg *types.MsgRemoveMFA) error {
	if msg.Creator == "" || msg.Type() == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "creator and MFA type cannot be empty")
	}
	return nil
}

// validateVerifyMFA validates MsgVerifyMFA
func (k Keeper) validateVerifyMFA(ctx context.Context, msg *types.MsgVerifyMFA) error {
	if msg.Creator == "" || msg.Code == "" {
		return sdkerrors.Wrap(types.ErrInvalidRequest, "creator and code cannot be empty")
	}
	return nil
}

// RateLimitMiddleware implements rate limiting middleware
func (k Keeper) RateLimitMiddleware(ctx sdk.Context, msg sdk.Msg) error {
	var did string

	// Extract DID from message
	switch msg := msg.(type) {
	case *types.MsgCreateDID:
		did = msg.Creator
	case *types.MsgUpdateDID:
		did = msg.Creator
	case *types.MsgDeleteDID:
		did = msg.Creator
	case *types.MsgLinkSocialIdentity:
		did = msg.Creator
	case *types.MsgUnlinkSocialIdentity:
		did = msg.Creator
	case *types.MsgAddMFA:
		did = msg.Creator
	case *types.MsgRemoveMFA:
		did = msg.Creator
	case *types.MsgVerifyMFA:
		did = msg.Creator
	default:
		return nil // No rate limiting for other message types
	}

	operation := sdk.MsgTypeURL(msg)
	return k.CheckRateLimit(ctx, did, operation)
}

// ValidateRequestMiddleware implements request validation middleware
func (k Keeper) ValidateRequestMiddleware(handler sdk.Handler) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		if err := k.ValidateBasic(ctx, msg); err != nil {
			return nil, err
		}
		return handler(ctx, msg)
	}
}

// LoggingMiddleware implements logging middleware
func (k Keeper) LoggingMiddleware(next sdk.Handler) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		// Log request
		k.Logger(ctx).Info("processing request",
			"module", types.ModuleName,
			"msg_type", sdk.MsgTypeURL(msg),
		)

		// Process request
		result, err := next(ctx, msg)

		// Log result
		if err != nil {
			k.Logger(ctx).Error("request failed",
				"module", types.ModuleName,
				"msg_type", sdk.MsgTypeURL(msg),
				"error", err,
			)
		}

		return result, err
	}
}

// ChainMiddleware chains multiple middleware functions together
func (k Keeper) ChainMiddleware(middlewares ...func(sdk.Handler) sdk.Handler) func(sdk.Handler) sdk.Handler {
	return func(handler sdk.Handler) sdk.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}

// SecurityMiddleware implements security checks
func (k Keeper) SecurityMiddleware(ctx sdk.Context, action string, handler func() error) error {
	// Check if DID exists and is active
	if err := k.validateRequest(ctx, action); err != nil {
		return err
	}

	// Check rate limits
	if err := k.CheckRateLimit(ctx, action, "security"); err != nil {
		return err
	}

	// Execute handler
	if err := handler(); err != nil {
		return err
	}

	// Log security event
	k.LogAuditEvent(ctx, &types.AuditEvent{
		EventType: "security_check",
		Did:       action,
		Success:   true,
	})

	return nil
}

// validateRequest validates a request
func (k Keeper) validateRequest(ctx sdk.Context, did string) error {
	if !k.HasDID(ctx, did) {
		return sdkerrors.Wrap(types.ErrDIDNotFound, did)
	}

	if !k.isDIDActive(ctx, did) {
		return sdkerrors.Wrap(types.ErrDIDInactive, did)
	}

	return nil
}

// HasDID checks if a DID exists
func (k Keeper) HasDID(ctx sdk.Context, did string) bool {
	doc, found := k.GetDIDDocument(ctx, did)
	return found && doc.Status == types.Status_STATUS_ACTIVE
}

// isDIDActive checks if a DID is active
func (k Keeper) isDIDActive(ctx sdk.Context, did string) bool {
	doc, found := k.GetDIDDocument(ctx, did)
	return found && doc.Status == types.Status_STATUS_ACTIVE
}

// isActionAllowed checks if an action is allowed
func (k Keeper) isActionAllowed(ctx sdk.Context, action string) bool {
	return true
}

// getRateLimitKey constructs a key for rate limit storage
func (k Keeper) getRateLimitKey(did string, operation string) []byte {
	return []byte(fmt.Sprintf("%s%s/%s", types.RateLimitPrefix, did, operation))
}
