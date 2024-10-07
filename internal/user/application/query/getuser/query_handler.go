package getuser

import (
	"fmt"

	"github.com/hex-api-go/internal/user/domain/contract"
)

type QueryHandler struct {
	dataSource contract.UserDataSource
}

func NewQueryHandler(dataSource contract.UserDataSource) *QueryHandler {
	return &QueryHandler{dataSource}
}

func (h *QueryHandler) Handle(data any) (any, error) {
	fmt.Println("get user > handle", data)
	/* res, err := h.dataSource.WithGateway(data.DataSource).GetPerson()

	if err != nil {
		return nil, err
	}

	response := struct {
		Id         string
		Name       string
		Age        int
		BirthDate  string
		Email      string
		DataSource string
	}{
		res.Uuid(),
		res.GetName(),
		res.GetAge(),
		res.GetBirthDate(),
		res.GetEmail(),
		res.GetDataSource(),
	} */
	return nil, nil
}
