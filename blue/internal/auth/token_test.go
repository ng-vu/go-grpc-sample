package auth

import (
	"os"
	"testing"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/base/redis"
)

var (
	gFoo Generator
	gBar Generator
)

func TestMain(M *testing.M) {
	endpoint := "redis:6379"
	redisPool := &redigo.Pool{
		MaxIdle:     50,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", endpoint)
			if err != nil {
				ll.Fatal("Redis: unable to dial on %s", l.Error(err), l.String("endpoint", endpoint))
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			ll.Fatal("err", l.Error(err))
			return err
		},
	}
	rStore := redis.New(redisPool)
	gFoo = NewGenerator("foo", rStore)
	gBar = NewGenerator("bar", rStore)

	os.Exit(M.Run())
}

func TestGenerateValidateRevokeToken(t *testing.T) {
	id := uuid.NewV4().String()
	var tokenStr string

	t.Run("Generate", func(t *testing.T) {
		tok, err := gFoo.Generate(id, DefaultTTL)
		require.NoError(t, err)

		tokenStr = tok.TokenStr
	})

	t.Run("Validate", func(t *testing.T) {
		tok, err := gFoo.Validate(tokenStr)
		require.NoError(t, err)

		assert.Equal(t, id, tok.UserID)
	})

	t.Run("Revoke", func(t *testing.T) {
		err := gFoo.Revoke(tokenStr)
		require.NoError(t, err)

		_, err = gFoo.Validate(tokenStr)
		assert.EqualError(t, err, "Invalid token")
	})
}

func TestGenerateWithValueToken(T *testing.T) {
	id := uuid.NewV4().String()
	t, err := gFoo.GenerateWithValue(id, "foo", DefaultTTL)
	if err != nil {
		T.Fatal(err)
	}

	defer func() {
		if err := gFoo.Revoke(t.TokenStr); err != nil {
			panic(err)
		}
	}()

	got, err := gFoo.Validate(t.TokenStr)
	if err != nil {
		T.Fatal(err)
	}

	if got.Value != "foo" {
		T.Fatal("Token value was not set correctly")
	}
}

func TestTokenWithTTL(T *testing.T) {
	id := uuid.NewV4().String()
	t, err := gFoo.Generate(id, 1)
	if err != nil {
		T.Fatal(err)
	}

	// Make sure TTL passed
	time.Sleep(1005 * time.Millisecond)

	if _, err := gFoo.Validate(t.TokenStr); err == nil {
		defer func() {
			if err := gFoo.Revoke(t.TokenStr); err != nil {
				panic(err)
			}
		}()

		T.Fatal("Key does not expired")
	}
}
