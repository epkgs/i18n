package errorx

import "errors"

type HttpError struct {
	error

	code       int // server error code
	httpStatus int // http status code
}

// WrapHttpError 是一个高阶函数，用于创建一个错误包装器，该包装器会将错误与指定的状态码关联。
// 这个函数主要解决了在处理错误时，需要统一设置自定义错误码和HTTP状态码的需求，通过返回一个闭包来实现。
// 参数:
//
//	code - 自定义的错误码，用于标识特定的错误情况。
//	httpStatus - HTTP状态码，用于标识HTTP请求的执行结果。
//
// 返回值:
//
//	Wrapper - 一个函数类型，接受一个错误作为输入，返回一个经过包装后的错误，该错误包含了原始错误、自定义错误码和HTTP状态码。
func WrapHttpError(code, httpStatus int) Wrapper[*HttpError] {
	return func(err *Error) *HttpError {
		return &HttpError{
			error:      err,
			code:       code,
			httpStatus: httpStatus,
		}
	}
}

// 返回错误状态码。
// 该方法属于withErrStatus类型，用于获取错误状态码。
// 主要用途是提供一种方式来获取错误的具体状态码，以便于错误处理和日志记录。
func (e *HttpError) Code() int {
	// 返回错误状态码。
	return e.code
}

// HttpStatus 返回错误状态码
//
// 该方法用于获取 HttpError 类型对象的 HTTP 状态码。
// 它返回一个整数类型的 HTTP 状态码，表示与错误相关的 HTTP 状态。
func (e *HttpError) HttpStatus() int {
	return e.httpStatus
}

// Unwrap 返回错误的底层原因（cause）。
// 此方法允许错误处理机制能够访问Error类型内部封装的实际错误。
// 参数: 无
// 返回值: error，代表错误的底层原因。
func (e *HttpError) Unwrap() error {
	return e.error
}

// Wrap 设置Error类型的cause字段
// 该方法用于将一个错误标记为另一个错误的直接原因，便于错误追踪和处理
// 参数:
//
//	err error: 导致当前错误的原始错误，不能为空
func (e *HttpError) Wrap(err error) error {
	e.error = err
	return e
}

func ErrorCode(err error) (int, bool) {
	var coder interface{ Code() int }
	if errors.As(err, &coder) {
		return coder.Code(), true
	}

	return 1, false
}

func HttpStatus(err error) (int, bool) {
	var coder interface{ HttpStatus() int }
	if errors.As(err, &coder) {
		return coder.HttpStatus(), true
	}

	return 200, false
}

type TraceableError struct {
	error

	traceID string
}

// WrapTraceableError 是一个返回 Wrapper 类型的函数。
// 它的作用是创建并返回一个闭包，该闭包用于将特定的 traceID 与错误信息关联。
// 参数 traceID 是一个字符串，代表了需要与错误信息关联的追踪ID。
func WrapTraceableError(traceID string) Wrapper[*TraceableError] {
	// 返回一个闭包，该闭包接收一个错误 err，并使用 WrapTraceableError 函数将 traceID 与 err 关联起来。
	// 这样做可以让错误信息携带特定的追踪ID，便于问题定位和追踪。
	return func(err *Error) *TraceableError {
		return &TraceableError{
			error:   err,
			traceID: traceID,
		}
	}
}

// TraceID 返回错误相关的追踪ID。
// 这个方法使得开发者可以在处理错误时获取到相关的追踪ID，进而进行问题定位和日志记录。
func (e *TraceableError) TraceID() string {
	return e.traceID
}

// Unwrap 返回错误的底层原因（cause）。
// 此方法允许错误处理机制能够访问Error类型内部封装的实际错误。
// 参数: 无
// 返回值: error，代表错误的底层原因。
func (e *TraceableError) Unwrap() error {
	return e.error
}

// Wrap 设置Error类型的cause字段
// 该方法用于将一个错误标记为另一个错误的直接原因，便于错误追踪和处理
// 参数:
//
//	err error: 导致当前错误的原始错误，不能为空
func (e *TraceableError) Wrap(err error) error {
	e.error = err
	return e
}

func GetTraceID(err error) (string, bool) {
	var tracer interface{ TraceID() string }
	if errors.As(err, &tracer) {
		return tracer.TraceID(), true
	}

	return "", false
}
