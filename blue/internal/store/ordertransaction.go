package store

import (
	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
)

func NewOrderTransactionStore(db DB) *OrderTransactionStore {
	return &OrderTransactionStore{
		db: db,
	}
}

// OrderTransactionStore ...
type OrderTransactionStore struct {
	db DB

	model.Store
}

func (s *OrderTransactionStore) GetByID(ID model.ID) (m *model.OrderTransaction, err error) {
	err = s.db.GetByID(ID, &m).Error
	return
}

func (s *OrderTransactionStore) GetByOrderID(OrderID model.ID) ([]*model.OrderTransaction, error) {
	panic("TODO")
}
