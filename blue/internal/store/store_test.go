package store_test

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
	"github.com/ng-vu/go-grpc-sample/blue/internal/bluetest"
	. "github.com/ng-vu/go-grpc-sample/blue/internal/store"
)

var (
	db DB
	ll = l.New()

	orderStore    *OrderStore
	providerStore *ServiceProviderStore
)

func init() {
	db = bluetest.SetupForTesting()

	orderStore = NewOrderStore(db)
	providerStore = NewServiceProviderStore(db)
}

// ValidID ...
func ValidID(t *testing.T, ID model.ID, expectedInfix model.Infix) {
	if len(ID) != 26 {
		require.Fail(t, "Invalid ID length: id=`%v`", ID)
		return
	}

	infix := ID.Infix()
	if infix == model.InfixInvalid {
		require.Fail(t, "Invalid infix: id=`%v` infix=`%v`", ID, infix)
		return
	}

	if expectedInfix != model.Infix(0) && infix != expectedInfix {
		require.Fail(t, "Unexpected infix: id=`%v` infix=`%v` expectedInfix=`%v`", ID, infix, expectedInfix)
		return
	}
}
