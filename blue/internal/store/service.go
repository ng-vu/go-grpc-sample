package store

import (
	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
)

func NewServiceStore(db DB) *ServiceStore {
	return &ServiceStore{
		db: db,
	}
}

type ServiceStore struct {
	db DB

	model.Store
}

func (s *ServiceStore) GetByID(ID model.ID) (*model.Service, error) {
	panic("TODO")
}
