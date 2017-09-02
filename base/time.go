package base

import (
	"time"
)

// Timestamp represents time as number of milliseconds from 1970.
type Timestamp int64

// ToTimestamp converts from Go time to timestamp (in millisecond).
func ToTimestamp(t time.Time) Timestamp {
	return Timestamp(t.UnixNano() / 1e6)
}

// Millis converts from Go time to timestamp (in millisecond).
func Millis(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixNano() / 1e6
}

// MillisP converts from Go time (with pointer) to timestamp (in millisecond).
func MillisP(t *time.Time) int64 {
	if t == nil || t.IsZero() {
		return 0
	}
	return t.UnixNano() / 1e6
}

// FromMillis converts millisecond to Go time
func FromMillis(t int64) time.Time {
	return Timestamp(t).ToTime()
}

// FromMillisP converts millisecond to Go time with pointer
func FromMillisP(t int64) *time.Time {
	if t == 0 {
		return nil
	}
	tt := Timestamp(t).ToTime()
	return &tt
}

// ToTime converts from timestamp to Go time.
func (t Timestamp) ToTime() time.Time {
	if t == 0 {
		return time.Time{}
	}
	return time.Unix(int64(t)/1e3, int64(t%1e3)*1e6).UTC()
}

// Unix converts Timestamp to seconds
func (t Timestamp) Unix() int64 {
	return int64(t) / 1e3
}

// UnixNano extracts nanoseconds from Timestamp
func (t Timestamp) UnixNano() int64 {
	return int64(t) * 1e6
}

// After reports whether the time instant t is after u
func (t Timestamp) After(u Timestamp) bool {
	return t > u
}

// Before reports whether the time instant t is before u
func (t Timestamp) Before(u Timestamp) bool {
	return t < u
}

// Add adds duration to timestamp
func (t Timestamp) Add(d time.Duration) Timestamp {
	return t + Timestamp(d/1e6)
}

// Sub subs timestamp
func (t Timestamp) Sub(u Timestamp) time.Duration {
	return time.Duration((t - u) * 1e6)
}

// AddDays add a number of days to timestamp
func (t Timestamp) AddDays(days int) Timestamp {
	return t + Timestamp(days*24*60*60*1e3)
}

func (t Timestamp) String() string {
	return t.ToTime().Format(time.RFC3339)
}

// Millis returns the timestamp as number of milliseconds from 1970
func (t Timestamp) Millis() int64 {
	return int64(t)
}

// IsZero reports whether timestamp is zero
func (t Timestamp) IsZero() bool {
	return t == 0
}

// SQLFormat returns time string for using inside SQL query
func (t Timestamp) SQLFormat() string {
	return t.ToTime().Format("2016-01-02")
}

var zeroTime = time.Unix(0, 0)

// IsZeroTime checks whether given time is zero.
// When transport time using grpc, empty time is marshalled to time.Unix(0, 0).
func IsZeroTime(t time.Time) bool {
	return t.IsZero() || t.Equal(zeroTime)
}

// Now returns current time.
func Now() Timestamp {
	return ToTimestamp(time.Now())
}

// Since returns the time elapsed since t. It is shorthand for common.Now().Sub(t).
func Since(t Timestamp) time.Duration {
	return time.Since(t.ToTime())
}
