package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Return int into duration in minutes
func d(n int) time.Duration {
	return time.Duration(n) * time.Minute
}

func modifyHour(h int) time.Time {
	return time.Date(2022, 12, 2, h, 0, 0, 0, time.Now().Location())
}

func TestSleep(t *testing.T) {
	sevenMinAndHalve := d(7) + 30*time.Second
	assert.Equal(t, sevenMinAndHalve, SleepTime(modifyHour(12)))
	assert.Equal(t, d(33), SleepTime(modifyHour(2)))
	assert.Equal(t, d(22), SleepTime(modifyHour(17)))
	assert.Equal(t, d(15), SleepTime(modifyHour(6)))
}
