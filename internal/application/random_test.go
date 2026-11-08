package application

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRandomUseCaseGenerateInRange(t *testing.T) {
	uc := NewRandomUseCase(WithWorkDelay(time.Millisecond))
	seen := make(map[int]bool)

	for i := 0; i < 1000; i++ {
		n, err := uc.Generate(context.Background())
		if err != nil {
			t.Fatalf("Generate() returned error: %v", err)
		}
		if n < minRandomNumber || n > maxRandomNumber {
			t.Fatalf("Generate() = %d, want in [%d, %d]", n, minRandomNumber, maxRandomNumber)
		}
		seen[n] = true
	}

	if len(seen) < 100 {
		t.Errorf("Generate() produced only %d distinct values over 1000 calls, want a well-spread range", len(seen))
	}
}

func TestRandomUseCaseGenerateWaitsForWorkDelay(t *testing.T) {
	const delay = 50 * time.Millisecond
	uc := NewRandomUseCase(WithWorkDelay(delay))

	start := time.Now()
	if _, err := uc.Generate(context.Background()); err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}
	if elapsed := time.Since(start); elapsed < delay {
		t.Errorf("Generate() returned after %v, want at least %v", elapsed, delay)
	}
}

func TestRandomUseCaseGenerateDefaultDelay(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real one-second wait in short mode")
	}
	uc := NewRandomUseCase()

	start := time.Now()
	if _, err := uc.Generate(context.Background()); err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}
	if elapsed := time.Since(start); elapsed < time.Second {
		t.Errorf("Generate() with default options returned after %v, want at least 1s", elapsed)
	}
}

func TestRandomUseCaseGenerateContextCancelled(t *testing.T) {
	uc := NewRandomUseCase(WithWorkDelay(time.Hour))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := uc.Generate(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Generate() error = %v, want context.Canceled", err)
	}
}
