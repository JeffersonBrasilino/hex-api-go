package facade

import (
	domaincontract "github.com/hex-api-go/internal/user/domain/contract"
	"github.com/hex-api-go/internal/user/domain/entity"
	"github.com/hex-api-go/internal/user/infrastructure/acl/contract"
	"github.com/hex-api-go/pkg/core"
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
	if f.gateways[f.activeGateway] == nil {
		return nil, &core.DependencyError{Message: "Invalid active gateway"}
	}
	result, err := f.gateways[f.activeGateway].GetPersonData()
	if err != nil {
		return nil, &core.DependencyError{Message: "Error getting person data"}
	}

	person := entity.NewPerson(result.Name, result.Age, result.BirthDate, result.Email, f.activeGateway)

	return person, nil
}

func (f *UserFacade) WithGateway(gateway string) domaincontract.UserDataSource {
	f.activeGateway = gateway
	return f
}
