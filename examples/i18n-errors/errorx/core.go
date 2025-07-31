package errorx

import (
	"github.com/epkgs/i18n"
	"github.com/epkgs/i18n/errorx"
)

func Definef[Args any](i18n *i18n.I18n, code int, format string, httpStatus int) *errorx.DefinitionF[*errorx.HttpError, Args] {
	return errorx.Definef[Args](i18n, format, errorx.WrapHttpError(code, httpStatus))
}

func Define(i18n *i18n.I18n, code int, format string, httpStatus int) *errorx.Definition[*errorx.HttpError] {
	return errorx.Define(i18n, format, errorx.WrapHttpError(code, httpStatus))
}
