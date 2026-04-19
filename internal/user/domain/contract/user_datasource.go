package contract

import "github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"

type UserDataSource interface {
	GetPerson() (*domain.Person, error)
	WithGateway(gateway string) UserDataSource
}
