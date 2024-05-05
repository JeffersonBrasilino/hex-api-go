package contracts

import "github.com/hex-api-go/internal/user/domain"

type UserRepository interface {
	Create(aggregate *domain.User)
}
