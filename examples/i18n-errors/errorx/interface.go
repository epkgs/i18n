package errorx

type cbuilder interface {
	Def() *Error // definition
}

type BuilderA0 interface {
	cbuilder
	New() *Error
}

type BuilderA1 interface {
	cbuilder
	NewA1(arg1 any) *Error
}

type BuilderA2 interface {
	cbuilder
	NewA2(arg1, arg2 any) *Error
}

type BuilderA3 interface {
	cbuilder
	NewA3(arg1, arg2, arg3 any) *Error
}

type BuilderA4 interface {
	cbuilder
	NewA4(arg1, arg2, arg3, arg4 any) *Error
}

type BuilderA5 interface {
	cbuilder
	NewA5(arg1, arg2, arg3, arg4, arg5 any) *Error
}

type BuilderAN interface {
	cbuilder
	NewAN(args ...any) *Error
}
