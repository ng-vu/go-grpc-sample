// Package redis abstracts Redis or a key-value database implementation
package redis

import (
	"encoding/json"
	"time"

	"github.com/ng-vu/go-grpc-sample/base/l"

	"github.com/garyburd/redigo/redis"
)

var ll = l.New()

// Store ...
type Store interface {
	Set(k string, v interface{}) error
	SetWithTTL(k string, v interface{}, ttl int) error
	Get(k string, v interface{}) error
	SetString(k string, v string) error
	SetStringWithTTL(k string, v string, ttl int) error
	GetString(k string) (string, error)
	GetStrings(p string) ([]string, error)
	SetUint64(k string, v uint64) error
	SetUint64WithTTL(k string, v uint64, ttl int) error
	GetUint64(k string) (uint64, error)
	GetTTL(k string) (int, error)
	IsExist(k string) bool
	Del(keys ...string) error
}

type redisStore struct {
	pool *redis.Pool
}

// New returns new Store
func New(pool *redis.Pool) Store {
	return &redisStore{pool: pool}
}

// NewWithPool returns new Redis Store with default pool config
func NewWithPool(address string) Store {
	redisPool := &redis.Pool{
		MaxIdle:     50,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(address)
			if err != nil {
				ll.Fatal("Unable to connect to Redis", l.Error(err), l.String("endpoint", address))
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			if err != nil {
				ll.Fatal("err", l.Error(err))
			}
			return nil
		},
	}
	return New(redisPool)
}

func (r redisStore) Set(k string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	c := r.pool.Get()
	defer c.Close()

	_, err = c.Do("SET", k, data)
	return err
}

// ttl: time in second
func (r redisStore) SetWithTTL(k string, v interface{}, ttl int) error {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	c := r.pool.Get()
	defer c.Close()

	_, err = c.Do("SETEX", k, ttl, data)
	return err
}

func (r redisStore) Get(k string, v interface{}) error {
	c := r.pool.Get()
	defer c.Close()

	data, err := redis.Bytes(c.Do("GET", k))
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func (r redisStore) SetString(k string, v string) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("SET", k, v)
	return err
}

// ttl: time in second
func (r redisStore) SetStringWithTTL(k string, v string, ttl int) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("SETEX", k, ttl, v)
	return err
}

func (r redisStore) GetString(k string) (string, error) {
	c := r.pool.Get()
	defer c.Close()

	s, err := redis.String(c.Do("GET", k))
	if err == redis.ErrNil {
		return "", nil
	}
	return s, err
}

func (r redisStore) GetStrings(p string) ([]string, error) {
	c := r.pool.Get()
	defer c.Close()

	values, err := redis.Strings(c.Do("KEYS", p))
	return values, err
}

func (r redisStore) SetUint64(k string, v uint64) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("SET", k, v)
	return err
}

// ttl: time in second
func (r redisStore) SetUint64WithTTL(k string, v uint64, ttl int) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("SETEX", k, ttl, v)
	return err
}

func (r redisStore) GetUint64(k string) (uint64, error) {
	c := r.pool.Get()
	defer c.Close()

	result, err := redis.Uint64(c.Do("GET", k))
	if err == redis.ErrNil {
		return 0, nil
	}
	return result, err
}

func (r redisStore) GetTTL(k string) (int, error) {
	c := r.pool.Get()
	defer c.Close()

	result, err := redis.Int(c.Do("TTL", k))
	if err == redis.ErrNil {
		return 0, nil
	}
	return result, err
}

func (r redisStore) IsExist(k string) bool {
	s, _ := r.GetString(k)
	return s != ""
}

func (r redisStore) Del(keys ...string) error {
	ks := make([]interface{}, len(keys))
	for i := range keys {
		ks[i] = keys[i]
	}

	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("DEL", ks...)

	return err
}
