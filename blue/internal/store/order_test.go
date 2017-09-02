package store_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
)

func TestOrder(t *testing.T) {
	var tOrder = &model.Order{
		ProviderID:        "SP_001",
		Code:              "ABCD",
		ProviderOrderID:   "P_ABCD",
		ProviderOrderCode: "P_ABCD",
		ExtraOrderCode:    "C_ABCD",
		CustomerPhone:     "0123456789",
		CustomerName:      "Trần Ngọc Thu Nhàn",
		Note:              "Note for order01",
	}

	var gOrderID model.ID
	var gOrder *model.Order

	t.Run("Create", func(t *testing.T) {
		order := *tOrder
		err := orderStore.Create(&order)
		require.NoError(t, err)

		ll.Info("Created order", l.Object("order", order))
		require.NotEmpty(t, order.ID)
		require.NotEmpty(t, order.CreatedAt)
		require.NotEmpty(t, order.UpdatedAt)
		ValidID(t, order.ID, model.InfixOrder)

		gOrder = &order
		gOrderID = gOrder.ID
	})

	t.Run("Create/NoPhone", func(t *testing.T) {
		order := *tOrder
		order.CustomerPhone = "  "

		err := orderStore.Create(&order)
		require.EqualError(t, err, "CustomerPhone: non zero value required;")
	})

	t.Run("GetByID", func(t *testing.T) {
		order, err := orderStore.GetByID(gOrderID)
		require.NoError(t, err)
		require.Equal(t, gOrder, order)
	})

	t.Run("GetByID/NotFound", func(t *testing.T) {
		_, err := orderStore.GetByID("invalid_id")
		require.EqualError(t, err, model.ErrNotFound.Error())
	})

	t.Run("GetByProviderCode", func(t *testing.T) {
		order, err := orderStore.GetByProviderCode("P_ABCD")
		require.NoError(t, err)
		require.Equal(t, gOrder, order)
	})

	t.Run("GetByProviderCode/NotFound", func(t *testing.T) {
		_, err := orderStore.GetByProviderCode("invalid")
		require.EqualError(t, err, model.ErrNotFound.Error())
	})

	t.Run("GetByCustomerCode", func(t *testing.T) {
		order, err := orderStore.GetByCustomerCode("C_ABCD")
		require.NoError(t, err)
		require.Equal(t, gOrder, order)
	})

	t.Run("GetByCustomerCode/NotFound", func(t *testing.T) {
		_, err := orderStore.GetByCustomerCode("invalid")
		require.EqualError(t, err, model.ErrNotFound.Error())
	})

	t.Run("GetByCustomerPhone", func(t *testing.T) {
		orders, err := orderStore.GetByCustomerPhone("0123456789")
		require.NoError(t, err)
		require.Equal(t, []*model.Order{gOrder}, orders)
	})
}
