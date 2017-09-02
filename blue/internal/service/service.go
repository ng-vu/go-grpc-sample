package service

import (
	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/internal/controller"
)

// Common errors
var (
	ErrTODO = ctrl.ErrTODO

	ll = l.New()
)

// InnerService encapsulates middle layers for logic
type InnerService struct {
	OrderCtrl    *ctrl.OrderCtrl    `inject:""`
	ProviderCtrl *ctrl.ProviderCtrl `inject:""`
	UserCtrl     *ctrl.UserCtrl     `inject:""`
}

// NewInnerService returns middle layers for logic
func NewInnerService() InnerService {
	return InnerService{}
}

// SetupServices ...
func (s InnerService) SetupServices() error {
	err := s.ProviderCtrl.LoadAll()
	return err
}
