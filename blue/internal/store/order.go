package store

import (
	"sync"

	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
)

// NewOrderStore ...
func NewOrderStore(db DB) *OrderStore {
	return &OrderStore{
		db: db,
	}
}

// OrderStore ...
type OrderStore struct {
	db DB

	model.Store
}

// Create ...
func (s *OrderStore) Create(order *model.Order) error {
	db := s.db.Create(order)
	return db.Error
}

// BulkCreate ...
func (s *OrderStore) BulkCreate(orders []*model.Order) []error {
	wg := sync.WaitGroup{}
	wg.Add(len(orders))
	errors := make([]error, len(orders))

	for i := range orders {
		order := orders[i]
		go func(i int) {
			defer wg.Done()

			ll.Info("Create", l.Int("i", i), l.Object("order", order))
			err := s.db.Create(order).Error
			errors[i] = err
			ll.Info("Created", l.Int("i", i), l.Object("order", order), l.Error(err))
		}(i)
	}
	wg.Wait()

	return errors
}

// GetByID ...
func (s *OrderStore) GetByID(ID model.ID) (*model.Order, error) {
	var order model.Order
	err := s.db.First(&order, "id = ?", ID).Error
	return &order, err
}

// GetByProviderCode ...
func (s *OrderStore) GetByProviderCode(code string) (*model.Order, error) {
	var order model.Order
	err := s.db.First(&order, "provider_order_code = ?", code).Error
	return &order, err
}

// GetByCustomerCode ...
func (s *OrderStore) GetByCustomerCode(code string) (*model.Order, error) {
	var order model.Order
	err := s.db.First(&order, "extra_order_code = ?", code).Error
	return &order, err
}

// GetByCustomerPhone ...
func (s *OrderStore) GetByCustomerPhone(phone string) (orders []*model.Order, err error) {
	err = s.db.Where(
		`customer_phone = ?`, phone,
	).Find(&orders).Error
	return
}
