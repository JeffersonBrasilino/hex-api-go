package contract

import "github.com/hex-api-go/internal/user/infrastructure/acl/translator"

type PersonGateway interface {
	GetPersonData() (*translator.PersonDto, error)
}
