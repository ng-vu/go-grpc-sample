package model

// Store serves as a workaround for dependency injection
type Store struct{}

func (s Store) _store() {}

// AgencyStaffStore ...
type AgencyStaffStore interface {
	Create(staff *AgencyStaff) (err error)
	GetByID(ID ID) (*AgencyStaff, error)
	UpdateInfo(ID ID, staff *AgencyStaff) error

	_store()
}

// UserInternalStore ...
type UserInternalStore interface {
	VerifyPassword(UserID ID, password string) (bool, error)
	UpdatePassword(UserID ID, newPassword string) error

	Disable(UserID ID) error
	Enable(UserID ID) error
	IsActive(UserID ID) (bool, error)

	_store()
}

// OrderStore ...
type OrderStore interface {
	GetByID(string) (*Order, error)
	GetByProviderCode(string) (*Order, error)
	GetByCustomerCode(string) (*Order, error)

	_store()
}

// ServiceStore ...
type ServiceStore interface {
	GetByID(ID ID) (*Service, error)

	_store()
}

// ServiceProviderStore ...
type ServiceProviderStore interface {
	GetByID(ID ID) (*ServiceProvider, error)

	_store()
}

// OrderTransactionStore ...
type OrderTransactionStore interface {
	GetByID(ID ID) (*OrderTransaction, error)
	GetByOrderID(OrderID ID) ([]*OrderTransaction, error)

	_store()
}
