package bluetest

import (
	"flag"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/config"
	"github.com/ng-vu/go-grpc-sample/blue/internal/store"
)

var (
	ll = l.New()

	isInit = false
	db     store.DB
)

// SetupForTesting ...
func SetupForTesting() store.DB {
	if flag.Lookup("test.v") == nil {
		ll.Fatal("This package should only run under `go test`")
	}

	if isInit {
		return db
	}
	isInit = true

	cfg := config.DefaultTest()
	return setupDB(cfg)
}

func setupDB(cfg *config.Config) store.DB {
	s := cfg.Postgres.ConnectionString()
	gormDB, err := gorm.Open("postgres", s)
	if err != nil {
		ll.Fatal("Unable to connect to postgres", l.Error(err), l.String("s", s))
	}

	_, err = gormDB.DB().Exec(fmt.Sprintf(`
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
		GRANT ALL ON SCHEMA public TO %v;
		GRANT ALL ON SCHEMA public TO public;
	`, cfg.Postgres.Username))
	if err != nil {
		ll.Warn("Drop all tables", l.Error(err))
	}

	ll.Info("Connected to postgres for testing. Dropped all data.", l.String("s", s))
	return store.SetupDB(gormDB)
}
