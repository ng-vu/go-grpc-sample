package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"

	"github.com/ng-vu/go-grpc-sample/base"
	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/base/redis"
	"github.com/ng-vu/go-grpc-sample/blue/config"
	"github.com/ng-vu/go-grpc-sample/blue/internal/store"
)

var (
	flConfigFile = flag.String("config-file", "", "Path to config file")

	ll  = l.New()
	cfg *config.Config
	ctx context.Context

	ctxCancel context.CancelFunc

	grpc1, grpc2 *grpc.Server
	http1, http2 *http.Server
)

func loadConfig() *config.Config {
	flag.Parse()
	cfg, err := config.Load(*flConfigFile)
	if err != nil {
		ll.Fatal("Unable to load config", l.Error(err))
	}
	return cfg
}

func connectRedis(cfg *config.Config) redis.Store {
	s := cfg.Redis.ConnectionString()
	redisStore := redis.NewWithPool(s)
	_, err := redisStore.GetString("_test_")
	if err != nil {
		ll.Fatal("Unable to connect to Redis", l.Error(err), l.String("ConnectionString", s))
	}
	return redisStore
}

func connectMongo(cfg *config.Config) *mgo.Database {
	s := cfg.Mongo.ConnectionString()
	session, err := mgo.Dial(s)
	if err != nil {
		ll.Fatal("Unable to connect to MongoDB", l.Error(err), l.Object("ConnectionString", s))
	}
	return session.DB(cfg.Mongo.Database)
}

func connectPostgres(cfg *config.Config) store.DB {
	s := cfg.Postgres.ConnectionString()
	db, err := gorm.Open("postgres", s)
	if err != nil {
		ll.Fatal("Unable to connect to PostgreSQL", l.Error(err), l.String("ConnectionString", s))
	}

	return store.SetupDB(db)
}

func main() {
	cfg = loadConfig()
	ll.Info("App started with config", l.Object("\nconfig", cfg))
	if cfg.Development {
		ll.Warn("DEVELOPMENT MODE ENABLED")
		base.SetDevelopmentMode(true)
	}

	ctx, ctxCancel = context.WithCancel(context.Background())
	go func() {
		defer ctxCancel()

		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
		ll.Info("Received OS signal", l.Stringer("signal", <-osSignal))
	}()

	startServers()

	// Wait for OS signal or any error from services
	<-ctx.Done()
	ll.Info("Waiting for all requests to finish")

	// Wait for maximum 15s
	go func() {
		timer := time.NewTimer(15 * time.Second)
		<-timer.C
		ll.Fatal("Force shutdown due to timeout!")
	}()

	// Graceful stop
	grpc1.GracefulStop()
	grpc2.GracefulStop()
	http1.Shutdown(context.Background())
	http2.Shutdown(context.Background())
	ll.Info("Gracefully stopped!")
}
