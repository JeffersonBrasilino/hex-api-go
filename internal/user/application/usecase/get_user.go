package usecase

import (
	"github.com/hex-api-go/internal/user/domain/contract"
	"github.com/hex-api-go/pkg/core"
)

type GetUserUseCase struct {
	dataSource contract.UserDataSource
}

func NewGetUserUseCase(dataSource contract.UserDataSource) *GetUserUseCase {
	return &GetUserUseCase{dataSource}
}

func (u *GetUserUseCase) Execute(dataSource string) *core.Result {
	result, err := u.dataSource.WithGateway(dataSource).GetPerson()

	if err != nil {
		return core.ResultFailure(err)
	}

	response := struct {
		Id         string
		Name       string
		Age        int
		BirthDate  string
		Email      string
		DataSource string
	}{
		result.GetUuid(),
		result.GetName(),
		result.GetAge(),
		result.GetBirthDate(),
		result.GetEmail(),
		result.GetDataSource(),
	}

	return core.ResultSuccess(response)
}
