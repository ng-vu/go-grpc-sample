package auth

import (
	"context"

	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
)

// Claim contains information for current user
type Claim struct {
	Token Token
}

type keyClaim struct{}

type providerClaim struct{}

// FromContext returns claim from context
func FromContext(ctx context.Context) (claim *Claim, ok bool) {
	claim, ok = ctx.Value(keyClaim{}).(*Claim)
	return
}

// NewContext ...
func NewContext(ctx context.Context, claim *Claim) context.Context {
	return context.WithValue(ctx, keyClaim{}, claim)
}

// ServiceProviderClaim ...
type ServiceProviderClaim struct {
	ID       model.ID
	Codename string
	Name     string
}

// ProviderFromContext returns claim from context
func ProviderFromContext(ctx context.Context) (claim ServiceProviderClaim, ok bool) {
	claim, ok = ctx.Value(providerClaim{}).(ServiceProviderClaim)
	return
}

// NewContextWithProvider ...
func NewContextWithProvider(ctx context.Context, claim ServiceProviderClaim) context.Context {
	if claim.ID == "" {
		ll.Error("Unexpected: empty ServiceProviderClaim", l.Stack())
		return ctx
	}
	return context.WithValue(ctx, providerClaim{}, claim)
}
