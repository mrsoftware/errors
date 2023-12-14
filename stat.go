package errors

// DefaultStat is used by Error and can be customized.
var DefaultStat Stater = &nopStater{} // nolint: gochecknoglobals

// SetDefaultStat is used to set Default Error Stater.
func SetDefaultStat(stater Stater) {
	DefaultStat = stater
}

// Stater is used to count errors.
type Stater interface {
	OnError(err error)
}

var _ Stater = &nopStater{}

type nopStater struct{}

// OnError do nothing.
func (n nopStater) OnError(err error) {}
