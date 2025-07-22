package i18n

type HTTPError interface {
	error
	Code() int
	HttpStatus() int
}

type withErrStatus struct {
	error

	code       int // server error code
	httpStatus int // http status code
}

var _ HTTPError = (*withErrStatus)(nil)

func (d *errorDefinition) WithStatus(code, httpStatus int) *errorDefinition {
	return d.With(WithErrStatusWrapper(code, httpStatus))
}

func WithErrStatusWrapper(code, httpStatus int) ErrorWrapper {
	return func(err error) error {
		return WithErrStatus(err, code, httpStatus)
	}
}

func WithErrStatus(err error, code int, httpStatus int) error {
	return &withErrStatus{
		error:      err,
		code:       code,
		httpStatus: httpStatus,
	}
}

func (e *withErrStatus) Code() int {
	return e.code
}

func (e *withErrStatus) HttpStatus() int {
	return e.httpStatus
}

type TracedError interface {
	error
	TraceID() string
}

var _ TracedError = (*withErrTraceID)(nil)

type withErrTraceID struct {
	error

	traceID string
}

func (d *errorDefinition) WithTraceID(traceID string) *errorDefinition {
	return d.With(WithErrTraceIDWrapper(traceID))
}

func WithErrTraceIDWrapper(traceID string) ErrorWrapper {
	return func(err error) error {
		return WithErrTraceID(err, traceID)
	}
}

func WithErrTraceID(err error, traceID string) error {
	return &withErrTraceID{
		error:   err,
		traceID: traceID,
	}
}

func (e *withErrTraceID) TraceID() string {
	return e.traceID
}
