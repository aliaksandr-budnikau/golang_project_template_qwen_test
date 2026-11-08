package application

import (
	"context"
	"math/rand/v2"
	"time"
)

// Bounds of the random number produced by RandomUseCase.Generate.
const (
	minRandomNumber = 1
	maxRandomNumber = 154
)

// RandomUseCase imitates a short unit of work (a configurable delay, one
// second by default) and then produces a random number.
type RandomUseCase struct {
	workDelay time.Duration
}

// RandomOption configures a RandomUseCase at construction time.
type RandomOption func(*RandomUseCase)

// WithWorkDelay sets how long Generate waits before answering. The default is
// one second; a non-positive delay means Generate does not wait at all.
func WithWorkDelay(d time.Duration) RandomOption {
	return func(uc *RandomUseCase) { uc.workDelay = d }
}

// NewRandomUseCase constructs a RandomUseCase.
func NewRandomUseCase(opts ...RandomOption) *RandomUseCase {
	uc := &RandomUseCase{workDelay: time.Second}
	for _, opt := range opts {
		opt(uc)
	}
	return uc
}

// Generate waits for the configured work delay — honouring ctx cancellation —
// and returns a random number in [1, 154], both bounds inclusive. If ctx is
// cancelled before the delay elapses, Generate returns the context error.
func (uc *RandomUseCase) Generate(ctx context.Context) (int, error) {
	timer := time.NewTimer(uc.workDelay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-timer.C:
	}

	// gosec G404: math/rand is intentional — this is a smoke-test number, not a
	// security-sensitive value, so crypto/rand would be overkill.
	return rand.IntN(maxRandomNumber-minRandomNumber+1) + minRandomNumber, nil //nolint:gosec // non-security random value
}
