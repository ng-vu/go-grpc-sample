package model

import (
	"crypto/rand"
	"database/sql"
	"database/sql/driver"

	"github.com/ng-vu/go-grpc-sample/base/idgen"
	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/oklog/ulid"
)

// Infix for ids
type Infix uint16

// InfixInvalid ...
const InfixInvalid = Infix(0)

func (i Infix) String() string {
	if s, ok := infix2str[i]; ok {
		return s
	}
	if i == InfixInvalid {
		return "[Invalid InfixID]"
	}
	return "[Unknown InfixID]"
}

const alphabet = ulid.Encoding

// ID represents an ID
type ID string

// EmptyID represents an empty id
const EmptyID = ID("")

var (
	ll = l.New()

	entropy = rand.Reader
)

func init() {
	// Test generating new ID
	_ = NewID(InfixOrder)
}

// NewID returns a new ID with infix
func NewID(infix Infix) ID {
	return ID(idgen.Generate(uint16(infix)).String())
}

// InfixString returns infix id
func (id ID) InfixString() string {
	if len(id) < 2 {
		return ""
	}
	return string(id[10:12])
}

// Infix returns infix id
func (id ID) Infix() Infix {
	p := id.InfixString()
	if infix, ok := str2infix[p]; ok {
		return infix
	}
	return InfixInvalid
}

// String is short-hand for sql.NullString.
type String string

// Scan implements the Scanner interface.
func (s *String) Scan(value interface{}) (err error) {
	var ns sql.NullString
	err = ns.Scan(value)
	if err == nil && ns.Valid {
		*s = String(ns.String)
	}
	return
}

// Value implements the driver Valuer interface.
func (s String) Value() (driver.Value, error) {
	if s == "" {
		return nil, nil
	}
	return string(s), nil
}
