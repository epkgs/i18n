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

// WithStatus 为错误定义添加状态码和HTTP状态码。
// 该方法接收两个整型参数：code 表示错误状态码，httpStatus 表示HTTP状态码。
// 返回值是指向 errorDefinition 类型的指针，允许进行链式调用。
func (d *errorDefinition) WithStatus(code, httpStatus int) *errorDefinition {
	// 调用 With 方法，并传入 ErrStatusWrapper 类型的实例，实现状态码和HTTP状态码的设置。
	return d.With(WithErrStatusWrapper(code, httpStatus))
}

// WithErrStatusWrapper 是一个高阶函数，用于创建一个错误包装器，该包装器会将错误与指定的状态码关联。
// 这个函数主要解决了在处理错误时，需要统一设置自定义错误码和HTTP状态码的需求，通过返回一个闭包来实现。
// 参数:
//
//	code - 自定义的错误码，用于标识特定的错误情况。
//	httpStatus - HTTP状态码，用于标识HTTP请求的执行结果。
//
// 返回值:
//
//	ErrorWrapper - 一个函数类型，接受一个错误作为输入，返回一个经过包装后的错误，该错误包含了原始错误、自定义错误码和HTTP状态码。
func WithErrStatusWrapper(code, httpStatus int) ErrorWrapper {
	return func(err error) error {
		return WithErrStatus(err, code, httpStatus)
	}
}

// WithErrStatus 是一个用于包装错误信息的函数，它允许错误信息携带额外的状态码和HTTP状态码。
// 这在需要对外提供API错误响应时特别有用，因为它可以提供更丰富的错误上下文信息。
//
// 参数:
// - err: error类型，表示发生的错误。
// - code: int类型，表示错误的自定义状态码，用于在程序内部进一步区分错误类型。
// - httpStatus: int类型，表示与错误相对应的HTTP状态码，用于在HTTP响应中返回。
//
// 返回值:
// 返回error类型，实际上是withErrStatus类型的错误，它包含了原始错误、自定义状态码和HTTP状态码。
// 这样的设计允许错误处理程序根据不同的错误类型和状态码做出更精细化的响应。
func WithErrStatus(err error, code int, httpStatus int) error {
	return &withErrStatus{
		error:      err,
		code:       code,
		httpStatus: httpStatus,
	}
}

// withErrStatusCode方法返回错误状态码。
// 该方法属于withErrStatus类型，用于获取错误状态码。
// 主要用途是提供一种方式来获取错误的具体状态码，以便于错误处理和日志记录。
func (e *withErrStatus) Code() int {
	// 返回错误状态码。
	return e.code
}

// HttpStatus 返回错误状态码
//
// 该方法用于获取 withErrStatus 类型对象的 HTTP 状态码。
// 它返回一个整数类型的 HTTP 状态码，表示与错误相关的 HTTP 状态。
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

// WithTraceID 为错误定义添加追踪ID。
// 这个方法允许将特定的追踪ID与错误信息关联起来，便于在错误日志中进行追踪和调试。
// 参数:
//
//	traceID (string): 追踪ID，用于标识和追踪错误发生的上下文。
//
// 返回值:
//
//	*errorDefinition: 返回错误定义自身，使得可以进行链式调用。
func (d *errorDefinition) WithTraceID(traceID string) *errorDefinition {
	return d.With(WithErrTraceIDWrapper(traceID))
}

// WithErrTraceIDWrapper 是一个返回 ErrorWrapper 类型的函数。
// 它的作用是创建并返回一个闭包，该闭包用于将特定的 traceID 与错误信息关联。
// 参数 traceID 是一个字符串，代表了需要与错误信息关联的追踪ID。
func WithErrTraceIDWrapper(traceID string) ErrorWrapper {
	// 返回一个闭包，该闭包接收一个错误 err，并使用 WithErrTraceID 函数将 traceID 与 err 关联起来。
	// 这样做可以让错误信息携带特定的追踪ID，便于问题定位和追踪。
	return func(err error) error {
		return WithErrTraceID(err, traceID)
	}
}

// WithErrTraceID is a function that wraps an error with a trace ID.
// This function is primarily used for error handling, where it adds a trace ID to the error for better tracking and debugging.
// Parameters:
//
//	err     - The original error to wrap.
//	traceID - The unique identifier for the error trace, used to track the error throughout the system.
//
// Return value:
//
//	Returns a new error that includes the original error and the trace ID.
func WithErrTraceID(err error, traceID string) error {
	return &withErrTraceID{
		error:   err,
		traceID: traceID,
	}
}

// TraceID 返回错误相关的追踪ID。
// 这个方法使得开发者可以在处理错误时获取到相关的追踪ID，进而进行问题定位和日志记录。
func (e *withErrTraceID) TraceID() string {
	return e.traceID
}
