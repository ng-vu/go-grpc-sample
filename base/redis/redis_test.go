package redis_test

import (
	"os"
	"testing"

	"github.com/ng-vu/go-grpc-sample/base/l"
	. "github.com/ng-vu/go-grpc-sample/base/redis"

	REQUIRE "github.com/stretchr/testify/require"
)

var (
	store Store
	ll    = l.New()
)

func init() {
	redisAddress := "redis://redis:6379"
	if os.Getenv("USE_DOCKER_HOST") == "1" {
		redisAddress = "redis://dockerhost:6379"
	}
	store = NewWithPool(redisAddress)
}

func TestGetSetInterface(T *testing.T) {
	type Foo struct {
		Bar int
		Baz string
	}

	T.Run("Test set interface", func(t *testing.T) {
		err := store.Set("foo", &Foo{10, "sample"})
		REQUIRE.Nil(t, err)
	})

	T.Run("Test get interface", func(t *testing.T) {
		var foo Foo
		err := store.Get("foo", &foo)
		REQUIRE.Nil(t, err)
		REQUIRE.Equal(t, 10, foo.Bar)
		REQUIRE.Equal(t, "sample", foo.Baz)
	})

	T.Run("Test set interface with ttl", func(t *testing.T) {
		err := store.SetWithTTL("foo1", &Foo{Bar: 10, Baz: "sample"}, 2)
		REQUIRE.Nil(t, err)
	})

	T.Run("Test get interface ttl", func(t *testing.T) {
		ttl, err := store.GetTTL("foo1")
		REQUIRE.Nil(t, err)

		t.Log("ttl: ", ttl)
		REQUIRE.True(t, (0 <= ttl) && (ttl <= 2) || (ttl == -2))
	})
}

func TestGetSetString(T *testing.T) {
	T.Run("Test set string", func(t *testing.T) {
		err := store.SetString("foo", "bar")
		REQUIRE.Nil(t, err)
	})

	T.Run("Test get string", func(t *testing.T) {
		foo, err := store.GetString("foo")
		REQUIRE.Nil(t, err)
		REQUIRE.Equal(t, "bar", foo)
	})

	T.Run("Test string is existed", func(t *testing.T) {
		exist := store.IsExist("foo")
		REQUIRE.True(t, exist)
	})

	T.Run("Test delete string", func(t *testing.T) {
		err := store.Del("foo")
		REQUIRE.Nil(t, err)

		exist := store.IsExist("foo")
		REQUIRE.False(t, exist)
	})

	T.Run("Test set string with ttl", func(t *testing.T) {
		err := store.SetStringWithTTL("string1", "string1", 2)
		REQUIRE.Nil(t, err)
	})
}

func TestGetStrings(t *testing.T) {
	err := store.SetString("t:foo", "bar")
	REQUIRE.Nil(t, err)
	err = store.SetString("t:bar", "baz")
	REQUIRE.Nil(t, err)

	values, err := store.GetStrings("t:*")
	REQUIRE.Nil(t, err)
	REQUIRE.NotEmpty(t, values)
}

func TestGetSetUint64(T *testing.T) {
	T.Run("Test set uint64", func(t *testing.T) {
		err := store.SetUint64("ten", 10)
		REQUIRE.Nil(t, err)
	})

	T.Run("Test get uint64", func(t *testing.T) {
		ten, err := store.GetUint64("ten")
		REQUIRE.Nil(t, err)
		REQUIRE.Equal(t, uint64(10), ten)
	})

	T.Run("Test set uint64 with ttl", func(t *testing.T) {
		err := store.SetUint64WithTTL("six", 6, 6)
		REQUIRE.Nil(t, err)
	})
}

func TestDel(t *testing.T) {
	values, err := store.GetStrings("t:*")
	REQUIRE.NoError(t, err)
	REQUIRE.NotEmpty(t, values)

	err = store.Del(values...)
	REQUIRE.NoError(t, err)

	values, err = store.GetStrings("t:*")
	REQUIRE.NoError(t, err)
	REQUIRE.Empty(t, values)

}
