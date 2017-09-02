package store

import (
	"errors"
	"strings"

	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
)

// NewServiceProviderStore ...
func NewServiceProviderStore(db DB) *ServiceProviderStore {
	return &ServiceProviderStore{
		db: db,
	}
}

// ServiceProviderStore ...
type ServiceProviderStore struct {
	db DB

	model.Store
}

// GetByID ...
func (s *ServiceProviderStore) GetByID(id model.ID) (*model.ServiceProvider, error) {
	var sp model.ServiceProvider
	err := s.db.First(&sp, `id = ?`, id).Error
	return &sp, err
}

// GetAll ...
func (s *ServiceProviderStore) GetAll() ([]*model.ServiceProvider, error) {
	var providers []*model.ServiceProvider
	err := s.db.Find(&providers).Error
	return providers, err
}

// GetByCodename ...
func (s *ServiceProviderStore) GetByCodename(codename string) (*model.ServiceProvider, error) {
	codename = strings.ToLower(strings.TrimSpace(codename))
	var sp model.ServiceProvider
	err := s.db.First(&sp, `codename = ?`, codename).Error
	return &sp, err
}

// Create ...
func (s *ServiceProviderStore) Create(sp *model.ServiceProvider) error {
	err := s.db.Create(sp).Error
	return err
}

// UpdateSecret ...
func (s *ServiceProviderStore) UpdateSecret(id model.ID, secret string) error {
	if secret == "" {
		return errors.New("Empty secret")
	}

	// Update secret without touching UpdatedAt
	err := s.db.Model(&model.ServiceProvider{ID: id}).
		UpdateColumn("secret", secret).Error
	return err
}
