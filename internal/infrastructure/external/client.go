// Package external contains adapters for calling external systems (third-party
// APIs, message brokers, etc.).
package external

import "context"

// Greeter is an example port for an external system.
type Greeter interface {
	// Greet returns a greeting for the given name.
	Greet(ctx context.Context, name string) (string, error)
}

// StaticGreeter is a no-op Greeter for local development and tests.
type StaticGreeter struct{}

// NewStaticGreeter constructs a StaticGreeter.
func NewStaticGreeter() *StaticGreeter {
	return &StaticGreeter{}
}

// Greet returns a canned greeting.
func (g *StaticGreeter) Greet(_ context.Context, name string) (string, error) {
	return "hello, " + name, nil
}

// Compile-time assertion that *StaticGreeter satisfies Greeter.
var _ Greeter = (*StaticGreeter)(nil)
