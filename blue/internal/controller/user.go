package ctrl

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/base/redis"
	"github.com/ng-vu/go-grpc-sample/base/validate"
	"github.com/ng-vu/go-grpc-sample/blue/internal/auth"
	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
	"github.com/ng-vu/go-grpc-sample/blue/internal/store"
)

type UserCtrl struct {
	AgencyStaffStore  *store.AgencyStaffStore  `inject:""`
	UserInternalStore *store.UserInternalStore `inject:""`

	ProviderCtrl *ProviderCtrl `inject:""`

	redis          redis.Store
	tokenGenerator auth.Generator
}

func NewUserCtrl(redisStore redis.Store, tokenGenerator auth.Generator) *UserCtrl {
	return &UserCtrl{
		redis:          redisStore,
		tokenGenerator: tokenGenerator,
	}
}

func (c *UserCtrl) CreateStaff(staff *model.AgencyStaff, password string) (userID model.ID, err error) {
	if password == "" {
		return model.EmptyID, grpc.Errorf(codes.InvalidArgument, "Missing password")
	}
	err = c.AgencyStaffStore.Create(staff, password)
	return staff.ID, err
}

func (c *UserCtrl) StaffLoginByPhone(ctx context.Context, phone, password string) (*model.AgencyStaff, *auth.Token, error) {
	phone, phoneOk := validate.NormalizePhone(phone)
	if !phoneOk {
		ll.Error("Invalid phone number", l.String("phone", phone))
		return nil, nil, grpc.Errorf(codes.InvalidArgument, "Invalid phone number")
	}

	staff, err := c.AgencyStaffStore.GetByPhone(phone)
	if err == model.ErrNotFound {
		ll.Error("LoginByPhone NotFound", l.String("phone", phone))
		return nil, nil, grpc.Errorf(codes.Unauthenticated, "Incorrect phone number or password")
	}
	if err != nil {
		ll.Error("Error", l.Error(err))
		return nil, nil, ErrInternal
	}

	userID := staff.ID
	ok, err := c.UserInternalStore.VerifyPassword(userID, password)
	if err != nil {
		ll.Error("Error", l.Error(err))
		return nil, nil, ErrInternal
	}

	if !ok {
		return nil, nil, grpc.Errorf(codes.Unauthenticated, "Incorrect phone number or password")
	}

	token, err := c.tokenGenerator.Generate(string(userID), TTLAccessToken)
	if err != nil {
		ll.Error("Error", l.Error(err))
		return nil, nil, ErrInternal
	}

	return staff, &token, nil
}

func (c *UserCtrl) Revoke(tokenStr string) error {
	return c.tokenGenerator.Revoke(tokenStr)
}
