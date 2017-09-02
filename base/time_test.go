package base

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp(t *testing.T) {
	// We set nanosecond to 321000000 because Timestamp only support milisecond
	now := time.Date(2020, 10, 20, 10, 20, 30, 321*1e6, time.UTC)
	tsNow := ToTimestamp(now)

	t.Run("ToTimestamp / ToTime", func(t *testing.T) {
		assert.Equal(t, Timestamp(1603189230321), tsNow)

		_now := tsNow.ToTime()
		assert.Equal(t, now, _now)
	})

	t.Run("Millis", func(t *testing.T) {
		assert.Equal(t, Timestamp(1603189230321), tsNow)
		assert.Equal(t, int64(1603189230321), Millis(now))
	})

	t.Run("Unix", func(t *testing.T) {
		assert.Equal(t, now.Unix(), tsNow.Unix())
	})

	t.Run("UnixNano", func(t *testing.T) {
		assert.Equal(t, now.UnixNano(), tsNow.UnixNano())
	})

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, now.Format(time.RFC3339), tsNow.String())
	})

	t.Run("After & Before", func(t *testing.T) {
		assert.True(t, tsNow.After(tsNow-10))
		assert.False(t, tsNow.After(tsNow))
		assert.False(t, tsNow.After(tsNow+10))

		assert.False(t, tsNow.Before(tsNow-10))
		assert.False(t, tsNow.Before(tsNow))
		assert.True(t, tsNow.Before(tsNow+10))
	})

	t.Run("Add & Sub", func(t *testing.T) {
		assert.Equal(t, now.Add(10*time.Hour), tsNow.Add(10*time.Hour).ToTime())
		assert.Equal(t, tsNow.Sub(tsNow.Add(10*time.Hour)), -10*time.Hour)
	})

	t.Run("AddDays", func(t *testing.T) {
		assert.Equal(t, now.AddDate(0, 0, 12), tsNow.AddDays(12).ToTime())
	})

	t.Run("IsZeroTime", func(t *testing.T) {
		assert.Equal(t, true, IsZeroTime(time.Time{}))
		assert.Equal(t, true, IsZeroTime(time.Unix(0, 0)))
		assert.Equal(t, false, IsZeroTime(now))

		assert.True(t, Timestamp(0).IsZero())
	})
}
