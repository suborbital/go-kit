package observability

import (
	"time"
)

type timer struct {
	start time.Time
}

// NewTimer returns a pointer to a timer instance with the start time set as the creation date.
func NewTimer() *timer {
	return &timer{
		start: time.Now(),
	}
}

// Observe returns the number of milliseconds between when the timer was constructed, and when this function is called
// as an int64.
func (t *timer) Observe() int64 {
	return time.Since(t.start).Milliseconds()
}
