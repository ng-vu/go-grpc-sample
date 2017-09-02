package ctrl

import (
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/base/redis"
	"github.com/ng-vu/go-grpc-sample/blue/internal/auth"
	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
	"github.com/ng-vu/go-grpc-sample/blue/internal/store"
)

// ProviderCtrl ...
type ProviderCtrl struct {
	ServiceStore         *store.ServiceStore         `inject:""`
	ServiceProviderStore *store.ServiceProviderStore `inject:""`

	redis redis.Store

	loaded       bool
	mapProviders map[string]*model.ServiceProvider
}

// LoadAll ...
func (c *ProviderCtrl) LoadAll() error {
	if c.loaded {
		ll.Warn("Already initialized")
		return nil
	}

	providers, err := c.ServiceProviderStore.GetAll()
	if err != nil {
		return err
	}

	ids := make([]string, len(providers))
	mapProviders := make(map[string]*model.ServiceProvider)
	for i, sp := range providers {
		mapProviders[string(sp.ID)] = sp
		ids[i] = string(sp.ID)
	}
	c.mapProviders = mapProviders
	ll.Info("Loaded all providers: " + strings.Join(ids, ", "))
	return nil
}

// RefreshCache ...
func (c *ProviderCtrl) RefreshCache() {
	c.loaded = false
	if err := c.LoadAll(); err != nil {
		ll.Error("Unable to reload service providers", l.Error(err))
	}
}

// ValidateAPIKey ...
func (c ProviderCtrl) ValidateAPIKey(apikey string) (auth.ServiceProviderClaim, error) {
	parts := strings.Split(apikey, ":")
	if len(parts) != 2 {
		ll.Error("Invalid apikey", l.String("apikey", apikey), l.Object("parts", parts))
		return auth.ServiceProviderClaim{}, grpc.Errorf(codes.Unauthenticated, "Request unauthenticated")
	}
	id := parts[0]
	secret := parts[1]

	sp, ok := c.mapProviders[id]
	if !ok {
		ll.Error("Provider not found", l.String("apikey", apikey), l.String("id", id), l.Object("map", c.mapProviders))
		return auth.ServiceProviderClaim{}, grpc.Errorf(codes.Unauthenticated, "Request unauthenticated")
	}

	if secret == string(sp.Secret) {
		return auth.ServiceProviderClaim{
			ID:       sp.ID,
			Codename: string(sp.Codename),
			Name:     string(sp.Name),
		}, nil
	}
	return auth.ServiceProviderClaim{}, grpc.Errorf(codes.Unauthenticated, "Request unauthenticated")
}

// GetByID ...
func (c ProviderCtrl) GetByID(id model.ID) (*model.ServiceProvider, error) {
	sp, err := c.ServiceProviderStore.GetByID(id)
	if err != nil {
		return nil, err
	}

	sp.Secret = ""
	return sp, nil
}

// Create ...
func (c ProviderCtrl) Create(sp *model.ServiceProvider) (id model.ID, apikey string, err error) {
	secret := auth.RandomToken(auth.DefaultTokenLength)
	sp.Secret = model.String(secret)
	err = c.ServiceProviderStore.Create(sp)
	if err != nil {
		return "", "", err
	}
	c.RefreshCache()

	apikey = string(sp.ID) + ":" + secret
	return sp.ID, apikey, err
}
