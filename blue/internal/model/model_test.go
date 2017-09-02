package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateOrder(t *testing.T) {
	gOrder := &Order{
		Code:              "1234ABCD_CODE",
		ProviderID:        "SP001",
		ProviderOrderID:   "1234ABCD",
		ProviderOrderCode: "1234ABCD",
		ExtraOrderCode:    "1234ABCD",
		CustomerPhone:     "0123456789",
		CustomerName:      "Nguyễn Ngọc Minh",
	}

	t.Run("Valid", func(t *testing.T) {
		err := gOrder.Validate()
		require.NoError(t, err)
	})

	t.Run("Valid/NormalizeCustomerPhone", func(t *testing.T) {
		order := *gOrder
		order.CustomerPhone = "123-444-567"

		err := order.Validate()
		require.NoError(t, err)
		require.Equal(t, String("0123444567"), order.CustomerPhone)
	})

	t.Run("Invalid/ProviderCode", func(t *testing.T) {
		order := *gOrder
		order.ProviderOrderCode = "1234@ ABCD"

		err := order.Validate()
		require.EqualError(t, err, "ProviderOrderCode: 1234@ ABCD does not validate as code;;")
	})

	t.Run("Invalid/ProviderCode/Required", func(t *testing.T) {
		order := *gOrder
		order.ProviderOrderCode = ""

		err := order.Validate()
		require.EqualError(t, err, "ProviderOrderCode: non zero value required;")
	})
}
