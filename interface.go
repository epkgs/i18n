package i18n

import (
	"context"
)

type Translable interface {
	Translate(ctx context.Context) string
}
