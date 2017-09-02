package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"

	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/base/redis"
)

const (
	// DefaultTokenLength ...
	DefaultTokenLength = 32 // In bytes

	// DefaultTTL ttl in seconds
	DefaultTTL = 60 * 60 * 24 * 30

	// DefaultTokenPrefix is used when no generator name is provided
	DefaultTokenPrefix = "AT"
)

// ErrInvalid returns an error indicate that
// the token is invalid
var ErrInvalid = errors.New("Invalid token")

var ll = l.New()

// Validator interface
type Validator interface {
	Validate(tokenStr string) (Token, error)
}

// Token represents a token used in request/response
type Token struct {
	TokenStr  string
	SubjectID string
	UserID    string
	Value     string
}

// Store interface contains methods
// to perform token actions
type Store interface {
	Generate(userID string, ttl int) (Token, error)
	GenerateWithValue(userID string, value string, ttl int) (Token, error)
	Revoke(tokenStr string) error
	SetInfo(tokenStr string, value string) error
	GetInfo(tokenStr string) (string, error)
}

// Generator interface
type Generator interface {
	Validator
	Store
}

type generator struct {
	name       string
	redisStore redis.Store
}

// NewGenerator returns new token generator
func NewGenerator(name string, r redis.Store) Generator {
	if name == "" {
		name = DefaultTokenPrefix
	}
	return &generator{
		name:       name,
		redisStore: r,
	}
}

// toKey returns string that can be used as redis key
// from given token
func (g *generator) toKey(t Token) string {
	return t.SubjectID + ":" + t.TokenStr
}

// toValue returns string that can be used as redis value
// from given token
func (g *generator) toValue(t Token) string {
	return t.UserID + ":" + t.Value
}

// Generate creates token for given userID and TTL.
func (g *generator) Generate(userID string, ttl int) (Token, error) {
	t := Token{
		SubjectID: g.name,
		UserID:    userID,
	}

	return g.generate(t, ttl)
}

// GenerateWithValue creates token with given value for given userID and TTL.
func (g *generator) GenerateWithValue(userID, value string, ttl int) (Token, error) {
	t := Token{
		SubjectID: g.name,
		UserID:    userID,
		Value:     value,
	}

	return g.generate(t, ttl)
}

func (g *generator) generate(t Token, ttl int) (Token, error) {
	retry := 0
	for {
		token := RandomToken(DefaultTokenLength)
		t.TokenStr = token

		key := g.toKey(t)
		if g.redisStore.IsExist(key) {
			retry++
			if retry >= 3 {
				panic("Unable to generate token, retried 3 times!")
			}
			continue
		}

		value := g.toValue(t)
		err := g.redisStore.SetStringWithTTL(key, value, ttl)

		return t, err
	}
}

func (g *generator) Validate(token string) (Token, error) {
	t := Token{
		TokenStr:  token,
		SubjectID: g.name,
	}

	// Check if the token exist in database
	key := g.toKey(t)
	storedValue, err := g.redisStore.GetString(key)
	if err != nil || storedValue == "" {
		return t, ErrInvalid
	}

	s := strings.Split(storedValue, ":")
	t.UserID = s[0]
	t.Value = s[1]

	return t, nil
}

// Revoke deletes token from redis store.
func (g *generator) Revoke(tokenStr string) error {
	t := Token{
		TokenStr:  tokenStr,
		SubjectID: g.name,
	}
	key := g.toKey(t)
	err := g.redisStore.Del(key)
	if err != nil {
		ll.Error("Error revoking token", l.Error(err))
	}
	return err
}

// SetInfo set information for given token
func (g *generator) SetInfo(tokenStr string, value string) error {
	ttl, err := g.redisStore.GetTTL(tokenStr)
	if err != nil {
		return err
	}

	t := Token{TokenStr: tokenStr}
	key := g.toKey(t)
	err = g.redisStore.SetStringWithTTL(key, value, ttl)
	if err != nil {
		return err
	}

	return nil
}

// GetInfo return infomation for given token
func (g *generator) GetInfo(tokenStr string) (string, error) {
	storedValue, err := g.redisStore.GetString(tokenStr)
	if err != nil {
		return "", err
	}

	return storedValue, nil
}

// RandomToken generate new base64 string from random byte array with given length
func RandomToken(length int) string {
	token := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, token); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(token)
}
