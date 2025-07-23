package i18n

type ErrorWrapper interface {
	error
	Unwrap() error
	Wrap(error) error
	Cause() error
}
