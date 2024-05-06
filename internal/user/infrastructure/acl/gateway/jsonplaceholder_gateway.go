package gateway

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user/infrastructure/acl/translator"
	"github.com/hex-api-go/pkg/core"
)

type JsonPlaceholderGateway struct{}

func NewJsonPlaceholderGateway() *JsonPlaceholderGateway {
	return &JsonPlaceholderGateway{}
}

func (a *JsonPlaceholderGateway) GetPersonData() (*translator.PersonDto, error) {
	request := fiber.Get("https://jsonplaceholder.typicode.com/users")
	_, body, errs := request.Bytes()

	if len(errs) > 0 {
		return nil, &core.InternalError{Message: "error"}
	}

	var p []any
	err := json.Unmarshal(body, &p)
	if err != nil {
		return nil, &core.InternalError{Message: "unmarshal error"}
	}
	result := p[0].(map[string]any)
	return &translator.PersonDto{Name: result["name"].(string), Age: 30, BirthDate: "01-01-1990"}, nil
}
