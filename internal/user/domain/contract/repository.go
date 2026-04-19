package contract

import (
	"context"

	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
)

type UserRepository interface {
	Create(ctx context.Context, aggregate *domain.User) error
}
