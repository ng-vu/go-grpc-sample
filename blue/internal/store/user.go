package store

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"io"

	"github.com/jinzhu/gorm"
	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
)

func NewUserInternalStore(db DB) *UserInternalStore {
	return &UserInternalStore{
		db: db,
	}
}

type UserInternalStore struct {
	db DB

	model.Store
}

// This function for internal usage only
func (s *UserInternalStore) _create(userID model.ID, password string) error {
	user := &model.UserInternal{
		UserID:  userID,
		HashPwd: encodePassword(password),
	}
	return s.db.Create(user).Error
}

func (s *UserInternalStore) GetByID(UserID model.ID) (*model.UserInternal, error) {
	var user model.UserInternal
	err := s.db.First(&user, "user_id = ?", UserID).Error
	return &user, err
}

func (s *UserInternalStore) VerifyPassword(UserID model.ID, password string) (bool, error) {
	var user model.UserInternal
	err := s.db.First(&user, "user_id = ?", UserID).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return verifyPassword(user.HashPwd, password), nil
}

func (s *UserInternalStore) UpdatePassword(UserID model.ID, password string) error {
	err := s.db.Model(&model.UserInternal{UserID: UserID}).
		Update("hash_pwd", encodePassword(password)).Error
	return err
}

func (s *UserInternalStore) Disable(UserID model.ID) error {
	panic("TODO")
}

func (s *UserInternalStore) Enable(UserID model.ID) error {
	panic("TODO")
}

func (s *UserInternalStore) IsActive(UserID model.ID) (bool, error) {
	panic("TODO")
}

// SaltSize is salt size in bytes.
const SaltSize = 16

func encodePassword(password string) string {
	return hexa(saltedHashPassword([]byte(password)))
}

func verifyPassword(hashpwd, password string) bool {
	return isPasswordMatch(dehexa(hashpwd), []byte(password))
}

func saltedHashPassword(secret []byte) []byte {
	buf := make([]byte, SaltSize, SaltSize+sha1.Size)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		ll.Panic("Unable to read from rand.Reader", l.Error(err))
		panic(err)
	}

	h := sha1.New()
	_, err = h.Write(buf)
	if err != nil {
		ll.Error("Write to buffer", l.Error(err))
	}

	_, err = h.Write(secret)
	if err != nil {
		ll.Error("Write to buffer", l.Error(err))
	}

	return h.Sum(buf)
}

func isPasswordMatch(data, secret []byte) bool {
	if len(data) != SaltSize+sha1.Size {
		panic("wrong length of data")
	}

	h := sha1.New()
	_, err := h.Write(data[:SaltSize])
	if err != nil {
		ll.Error("Write to buffer", l.Error(err))
	}

	_, err = h.Write(secret)
	if err != nil {
		ll.Error("Write to buffer", l.Error(err))
	}

	return bytes.Equal(h.Sum(nil), data[SaltSize:])
}

func hexa(data []byte) string {
	return hex.EncodeToString(data)
}

func dehexa(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}
