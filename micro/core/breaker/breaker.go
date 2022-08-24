package breaker

type (
	Acceptable func(err error) bool

	Breaker interface {
		Name() string

		Counts() Counts

		State() State

		Do(req func() error) error

		DoWithAcceptable(req func() error, acceptable Acceptable) error

		DoWithFallback(req func() error, fallback func(err error) error) error

		DoWithFallbackAcceptable(req func() error, fallback func(err error) error, acceptable Acceptable) error
	}
)
