package errorx

type Builder struct {
	def *Error // definition
}

// definition
func (b *Builder) Def() *Error {
	return b.def
}

func (b *Builder) new(args ...any) *Error {
	return &Error{
		args: args,

		Item:       b.def.Item,
		code:       b.def.code,
		httpStatus: b.def.httpStatus,
	}
}

func (b *Builder) NewA0() *Error {
	return b.new()
}

func (b *Builder) NewA1(arg1 any) *Error {
	return b.new(arg1)
}

func (b *Builder) NewA2(arg1, arg2 any) *Error {
	return b.new(arg1, arg2)
}

func (b *Builder) NewA3(arg1, arg2, arg3 any) *Error {
	return b.new(arg1, arg2, arg3)
}

func (b *Builder) NewA4(arg1, arg2, arg3, arg4 any) *Error {
	return b.new(arg1, arg2, arg3, arg4)
}

func (b *Builder) NewA5(arg1, arg2, arg3, arg4, arg5 any) *Error {
	return b.new(arg1, arg2, arg3, arg4, arg5)
}

func (b *Builder) NewAN(args ...any) *Error {
	return b.new(args...)
}
