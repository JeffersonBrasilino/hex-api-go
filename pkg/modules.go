package pkg

import (
	"context"
)

type Module interface {
	Register(ctx context.Context) error
}