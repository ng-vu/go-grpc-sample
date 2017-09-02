package store

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
)

var (
	ll = l.New()
)

// DB wraps gorm.DB with extra functionality
type DB struct {
	*gorm.DB
}

type sqlTx interface {
	Commit() error
	Rollback() error
}

// SetupDB ...
func SetupDB(db *gorm.DB) DB {
	// Disable table name's pluralization globally
	db.SingularTable(true)

	// Connection pool
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Automatically migrate the schema, to keep it updated.
	db.AutoMigrate(
		&model.AgencyStaff{},
		&model.BlueStaff{},
		&model.UserInternal{},
		&model.Order{},
		&model.Service{},
		&model.ServiceProvider{},
		&model.OrderTransaction{},
		&model.MoneyTransaction{},
	)

	return DB{DB: db}
}

// GetByID ...
func (db DB) GetByID(ID, v interface{}) *gorm.DB {
	return db.DB.First(v, "id = ?", ID)
}

// Transact ...
func (db DB) Transact(txFn func(DB) error) (err error) {
	tx := db.DB.Begin()
	defer func() {
		e := recover()
		if e != nil {
			tx.Rollback()
			panic(e) // Throw after Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	return txFn(DB{DB: tx})
}

func init() {
	// Config time to work with PostgreSQL
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC().Truncate(1000 * time.Nanosecond)
	}
}
