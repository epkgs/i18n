package errorx

type Error struct {
	httpStatus int

	msg string
}

func New(httpStatus int, msg string) *Error {
	return &Error{
		httpStatus: httpStatus,
		msg:        msg,
	}
}

func (e *Error) HttpStatus() int {
	return e.httpStatus
}

func (e *Error) WithMsg(msg string) *Error {
	return &Error{
		httpStatus: e.httpStatus,
		msg:        msg,
	}
}

func (e *Error) Error() string {
	return e.msg
}
