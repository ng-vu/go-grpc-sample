package store_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
)

func TestServiceProvider(t *testing.T) {
	var gProviderID model.ID
	var gProvider = &model.ServiceProvider{
		Codename: "expl",
		Name:     "Example Service Provider",
	}

	t.Run("Create", func(t *testing.T) {
		sp := gProvider
		err := providerStore.Create(sp)
		require.NoError(t, err)

		ll.Info("Created service provider", l.Object("sp", sp))
		require.NotEmpty(t, sp.ID)
		require.NotEmpty(t, sp.CreatedAt)
		require.NotEmpty(t, sp.UpdatedAt)
		ValidID(t, sp.ID, model.InfixServiceProvider)

		gProviderID = sp.ID
	})

	t.Run("GetByID", func(t *testing.T) {
		sp, err := providerStore.GetByID(gProviderID)
		require.NoError(t, err)
		require.Equal(t, gProvider, sp)
	})

	t.Run("GetAll", func(t *testing.T) {
		providers, err := providerStore.GetAll()
		require.NoError(t, err)
		require.Equal(t, []*model.ServiceProvider{gProvider}, providers)
	})
}
