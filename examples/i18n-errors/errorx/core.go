package errorx

import (
	"github.com/epkgs/i18n"
	"github.com/epkgs/i18n/errorx"
)

func Definef[Args any](bundle *i18n.Bundle, code int, format string, httpStatus int) *errorx.DefinitionF[*errorx.HttpError, Args] {
	return errorx.Definef[Args](bundle, format, errorx.WrapHttpError(code, httpStatus))
}

func Define(bundle *i18n.Bundle, code int, format string, httpStatus int) *errorx.Definition[*errorx.HttpError] {
	return errorx.Define(bundle, format, errorx.WrapHttpError(code, httpStatus))
}
