package store

import (
	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
)

func NewAgencyStaffStore(db DB) *AgencyStaffStore {
	return &AgencyStaffStore{
		db: db,
	}
}

type AgencyStaffStore struct {
	db DB

	model.Store
}

func (s *AgencyStaffStore) Create(staff *model.AgencyStaff, password string) error {
	return s.db.Transact(func(db DB) error {
		staff.ID = model.EmptyID
		err := db.Create(staff).Error
		if err != nil {
			return err
		}

		userID := staff.ID
		if userID == model.EmptyID {
			ll.Error("Unexpected empty id", l.Object("staff", staff))
			return model.ErrUnexpected
		}

		userInternalStore := NewUserInternalStore(db)
		err = userInternalStore._create(userID, password)
		return err
	})
}

func (s *AgencyStaffStore) GetByID(ID model.ID) (staff *model.AgencyStaff, err error) {
	panic("TODO")
}

func (s *AgencyStaffStore) GetByPhone(phone string) (*model.AgencyStaff, error) {
	var staff model.AgencyStaff
	err := s.db.First(&staff, "phone = ?", phone).Error
	return &staff, err
}

func (s *AgencyStaffStore) UpdateInfo(ID model.ID, staff *model.AgencyStaff) error {
	panic("TODO")
}

func (s *AgencyStaffStore) VerifyPasswordWithPhone(phone string, password string) error {
	panic("TODO")
}

func (s *AgencyStaffStore) UpdatePassword() {
	panic("TODO")
}

func (s *AgencyStaffStore) Disable() {
	panic("TODO")
}
