package base

import (
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// Authorization definition
const (
	AuthorizationHeader = "Authorization"
	AuthorizationType   = "Bearer"
)

// AppendAccessToken ...
func AppendAccessToken(ctx context.Context, accessToken string) context.Context {
	return metadata.NewOutgoingContext(ctx,
		metadata.Pairs(
			AuthorizationHeader,
			strings.Join([]string{
				AuthorizationType,
				accessToken,
			}, " "),
		),
	)
}

// AppendMetadata ...
func AppendMetadata(ctx context.Context, pairs ...string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(pairs...))
}
