package facade

import (
	domaincontract "github.com/jeffersonbrasilino/hex-api-go/internal/user/domain/contract"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain/entity"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/acl/contract"
)

type facadeGateways map[string]contract.PersonGateway
type UserFacade struct {
	gateways      facadeGateways
	activeGateway string
}

func NewUserFacade(gateways facadeGateways) *UserFacade {
	return &UserFacade{gateways: gateways}
}

func (f *UserFacade) GetPerson() (*entity.Person, error) {
	/* if f.gateways[f.activeGateway] == nil {
		a := domain.NewError[*domain.DependencyError]("invalid active gateway")
		return nil, a
	} */
	result, err := f.gateways[f.activeGateway].GetPersonData()
	if err != nil {
		return nil, nil
	}

	person := entity.NewPerson(result.Name, result.Age, result.BirthDate, result.Email, f.activeGateway)

	return person, nil
}

func (f *UserFacade) WithGateway(gateway string) domaincontract.UserDataSource {
	f.activeGateway = gateway
	return f
}
