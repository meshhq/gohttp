package gohttp

import (
	"time"

	"github.com/meshhq/gohttp/Godeps/_workspace/src/github.com/cenk/backoff"
)

// GoHTTP Default backoff parameters.
const (
	DefaultInitialInterval     = 100 * time.Millisecond
	DefaultRandomizationFactor = 0.75
	DefaultMultiplier          = 2
	DefaultMaxInterval         = 3 * time.Second
	DefaultMaxElapsedTime      = 5 * time.Second
)

// Backoff returns a backoff.ExponentialBackOff algorithm with the default
// GoHTTP backoff policy.
func Backoff() *backoff.ExponentialBackOff {
	b := &backoff.ExponentialBackOff{
		InitialInterval:     DefaultInitialInterval,
		RandomizationFactor: DefaultRandomizationFactor,
		Multiplier:          DefaultMultiplier,
		MaxInterval:         DefaultMaxInterval,
		MaxElapsedTime:      DefaultMaxElapsedTime,
		Clock:               backoff.SystemClock,
	}
	if b.RandomizationFactor < 0 {
		b.RandomizationFactor = 0
	} else if b.RandomizationFactor > 1 {
		b.RandomizationFactor = 1
	}
	b.Reset()
	return b
}
