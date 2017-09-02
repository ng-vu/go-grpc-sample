package grpcTransport

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"github.com/ng-vu/go-grpc-sample/base/idgen"
	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/internal/auth"
	"google.golang.org/grpc/codes"
)

var ll = l.New()
var reqInfix = idgen.CalcInfix("RQ")

// AuthFunc ...
type AuthFunc func(ctx context.Context, fullMethod string) (context.Context, error)

// Authentication ...
func Authentication(validator auth.Validator, magicToken string, exceptions []string) AuthFunc {
	return func(ctx context.Context, fullMethod string) (context.Context, error) {
		for _, exception := range exceptions {
			if exception == fullMethod {
				return ctx, nil
			}
		}

		tokenStr, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			ll.Warn("No authorization header", l.String("method", fullMethod), l.Error(err))
			return ctx, err
		}

		if magicToken != "" && tokenStr == magicToken {
			ll.Warn("DEVELOPMENT: Authenticated with magic token")
			token := auth.Token{
				TokenStr: "MAGIC",
				UserID:   "",
			}
			return auth.NewContext(ctx, &auth.Claim{Token: token}), nil
		}

		token, err := validator.Validate(tokenStr)
		if err != nil {
			ll.Warn("Invalid token", l.String("token", tokenStr), l.Error(err))
			return ctx, grpc.Errorf(codes.Unauthenticated, "Request unauthenticated")
		}

		return auth.NewContext(ctx, &auth.Claim{Token: token}), nil
	}
}

// AuthUnaryServerInterceptor ...
func AuthUnaryServerInterceptor(authFunc AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx, err := authFunc(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// LogUnaryServerInterceptor returns middleware for logging with zap
func LogUnaryServerInterceptor(logger l.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			e := recover()
			if e != nil {
				logger.Error("Panic (Recovered)", l.Error(err), l.Stack())
				err = grpc.Errorf(codes.Internal, "Internal Error (%v)", e)
			}

			if err == nil {
				logger.Info(info.FullMethod, l.Interface("\n→", req), l.Interface("\n⇐", resp))
				return
			}
			logger.Error(info.FullMethod, l.Interface("\n→", req), l.String("\n⇐ERROR", err.Error()))
		}()

		// Append correlation id
		const correlationID = "correlation-id"
		inMD, _ := metadata.FromIncomingContext(ctx)
		var reqID string
		if ids, ok := inMD[correlationID]; ok && len(ids) > 0 {
			reqID = ids[0]
		} else {
			reqID = idgen.Generate(reqInfix).String()
		}
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(correlationID, reqID))

		return handler(ctx, req)
	}
}
