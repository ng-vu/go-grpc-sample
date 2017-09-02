package model

import (
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/ng-vu/go-grpc-sample/base/validate"
)

// Common errors
var (
	ErrNotFound = gorm.ErrRecordNotFound

	ErrEmptyID        = errors.New("Empty ID")
	ErrInvalidID      = errors.New("Invalid ID")
	ErrUnexpected     = errors.New("Internal Error (unexpected)")
	ErrNoProviderCode = errors.New("No provider code")
)

// StatusDone ...
type StatusDone string

// StatusDone constants
const (
	StatusInit = StatusDone("")
	StatusC    = StatusDone("C")
	StatusD    = StatusDone("D")
)

// BaseModel includes common fields
type BaseModel struct {
	ID        ID         `gorm:"primary_key"`
	CreatedAt time.Time  `gorm:"index"`
	UpdatedAt time.Time  `gorm:"index"`
	DeletedAt *time.Time `gorm:"index"`
}

// AgencyStaff ...
type AgencyStaff struct {
	ID        ID         `gorm:"primary_key"`
	CreatedAt time.Time  `gorm:"index"`
	UpdatedAt time.Time  `gorm:"index"`
	DeletedAt *time.Time `gorm:"index"`

	Name    String `gorm:"not null" valid:"name,required"`
	Email   String `gorm:"unique_index" valid:"email"`
	Phone   String `gorm:"not null;unique_index" valid:"phone,required"`
	Address String `gorm:"not null" valid:"required"`

	UserInternal *UserInternal `gorm:"polymorphic:User"`
}

// Validate ...
func (m *AgencyStaff) Validate() error {
	phone, _ := validate.NormalizePhone(string(m.Phone))
	name, _ := validate.NormalizeName(string(m.Name))

	m.Phone = String(phone)
	m.Name = String(name)
	return validate.Check(m)
}

// BeforeCreate implements callback for auto generating ID in gorm
func (m *AgencyStaff) BeforeCreate(scope *gorm.Scope) error {
	if err := m.Validate(); err != nil {
		return err
	}

	m.ID = NewID(InfixAgencyStaff)
	return nil
}

// BlueStaff ...
type BlueStaff struct {
	ID        ID         `gorm:"primary_key"`
	CreatedAt time.Time  `gorm:"index"`
	UpdatedAt time.Time  `gorm:"index"`
	DeletedAt *time.Time `gorm:"index"`

	Name  string `gorm:"not null" valid:"name,required"`
	Email string `gorm:"unique_index" valid:"email,required"`
	Phone string `gorm:"unique_index" valid:"phone,required"`

	UserInternal *UserInternal `gorm:"polymorphic:User"`
}

// BeforeCreate implements callback for auto generating ID in gorm
func (m *BlueStaff) BeforeCreate(scope *gorm.Scope) error {
	m.ID = NewID(InfixBlueStaff)
	return nil
}

// UserInternal ...
type UserInternal struct {
	UserID    ID         `gorm:"primary_key"`
	CreatedAt time.Time  `gorm:"index"`
	UpdatedAt time.Time  `gorm:"index"`
	DeletedAt *time.Time `gorm:"index"`

	UserType string `gorm:"index"`
	HashPwd  string
}

// BeforeCreate implements callback for auto generating ID in gorm
func (m *UserInternal) BeforeCreate(scope *gorm.Scope) error {
	if m.UserID == EmptyID {
		return ErrEmptyID
	}
	return nil
}

// Order ...
type Order struct {
	ID           ID         `gorm:"primary_key"`
	CreatedAt    time.Time  `gorm:"index"`
	UpdatedAt    time.Time  `gorm:"index"`
	DeletedAt    *time.Time `gorm:"index"`
	OrderTime    time.Time
	ExpectedTime *time.Time

	Code String `gorm:"unique_index" valid:"code"`
	Type String // Delivery, etc

	ServiceID ID `gorm:"index"`

	ProviderID      ID     `gorm:"unique_index:cudx_order_provider_order_id" valid:"code,required"`
	ProviderOrderID String `gorm:"unique_index:cudx_order_provider_order_id" valid:"code,required"`

	ProviderOrderCode String `gorm:"index" valid:"code,required"`
	ExtraOrderCode    String `gorm:"index" valid:"code"`

	LocationID      ID     `gorm:"index"`
	LocationAddress String `gorm:"address"`

	CustomerID          ID     `gorm:"index" valid:"alphanum"`
	CustomerReferenceID String `valid:"code"`
	CustomerName        String `valid:"name,required"`
	CustomerPhone       String `gorm:"index" valid:"phone,required"`
	CustomerEmail       String
	CustomerAddress     String

	Note     String
	InfoJSON String `gorm:"type:JSONB DEFAULT '{}'::JSONB"`

	TotalFee    int32 `gorm:"type:int"`
	TotalAmount int32 `gorm:"type:int"`

	Status         String
	StatusProvider String
	StatusMessage  String
	StatusDone     int32 `gorm:"type:int;index"` // -1, 0, 1
}

// Validate ...
func (m *Order) Validate() error {
	phone, _ := validate.NormalizePhone(string(m.CustomerPhone))

	m.CustomerPhone = String(phone)
	return validate.Check(m)
}

// BeforeCreate implements callback for auto generating ID in gorm
func (m *Order) BeforeCreate(scope *gorm.Scope) error {
	if err := m.Validate(); err != nil {
		return err
	}

	m.ID = NewID(InfixOrder)
	return nil
}

// Service ...
type Service struct {
	ID        ID         `gorm:"primary_key"`
	CreatedAt time.Time  `gorm:"index"`
	UpdatedAt time.Time  `gorm:"index"`
	DeletedAt *time.Time `gorm:"index"`

	Codename string `gorm:"unique_index"`
	Name     string

	ServiceProvider   ServiceProvider `gorm:"ForeignKey:ServiceProviderID"`
	ServiceProviderID ID
}

// ServiceProviderClaim ...
type ServiceProviderClaim struct {
	ID       ID
	Codename string
	Name     string
}

// ServiceProvider ...
type ServiceProvider struct {
	ID        ID         `gorm:"primary_key"`
	CreatedAt time.Time  `gorm:"index"`
	UpdatedAt time.Time  `gorm:"index"`
	DeletedAt *time.Time `gorm:"index"`

	Codename String `gorm:"unique_index" valid:"code,required"`
	Name     String `valid:"name,required"`
	Secret   String `valid:"code,required"`
}

// BeforeCreate implements callback for auto generating ID in gorm
func (m *ServiceProvider) BeforeCreate(scope *gorm.Scope) error {
	m.Codename = String(strings.ToLower(string(m.Codename)))
	m.ID = NewID(InfixServiceProvider)
	return nil
}

// OrderTransaction ...
type OrderTransaction struct {
	ID              ID         `gorm:"primary_key"`
	CreatedAt       time.Time  `gorm:"index"`
	UpdatedAt       time.Time  `gorm:"index"`
	DeletedAt       *time.Time `gorm:"index"`
	TransactionTime time.Time  `gorm:"index"`

	Order   *Order `gorm:"ForeignKey:OrderID"`
	OrderID ID     `gorm:"index"`

	AgencyStaff   *AgencyStaff `gorm:"ForeignKey:AgencyStaffID"`
	AgencyStaffID string       `gorm:"index"`

	FromID ID `gorm:"index"`
	ToID   ID `gorm:"index"`
}

// BeforeCreate implements callback for auto generating ID in gorm
func (m *OrderTransaction) BeforeCreate(scope *gorm.Scope) error {
	m.ID = NewID(InfixOrderTransaction)
	return nil
}

// MoneyTransaction ...
type MoneyTransaction struct {
	ID              ID         `gorm:"primary_key"`
	CreatedAt       time.Time  `gorm:"index"`
	UpdatedAt       time.Time  `gorm:"index"`
	DeletedAt       *time.Time `gorm:"index"`
	TransactionTime time.Time  `gorm:"index"`

	Order   *Order `gorm:"ForeignKey:OrderID"`
	OrderID string `gorm:"index"`

	TransactionType string
	FromID          ID `gorm:"index"`
	FromWallet      string
	ToID            ID `gorm:"index"`
	ToWallet        string
}

// BeforeCreate implements callback for auto generating ID in gorm
func (m *MoneyTransaction) BeforeCreate(scope *gorm.Scope) error {
	m.ID = NewID(InfixMoneyTransaction)
	return nil
}
