package contract

import "github.com/hex-api-go/internal/user/domain/entity"

type UserDataSource interface {
	GetPerson() (*entity.Person, error)
	WithGateway(gateway string) UserDataSource
}
