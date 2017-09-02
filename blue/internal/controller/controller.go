package ctrl

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/ng-vu/go-grpc-sample/base/l"
)

// Time to live
const (
	TTLAccessToken = 7 * 24 * 60 * 60 // 7 days
)

// Common errors
var (
	ErrInternal = grpc.Errorf(codes.Internal, "Internal Error")
	ErrTODO     = grpc.Errorf(codes.Unimplemented, "TODO")

	ll = l.New()
)
