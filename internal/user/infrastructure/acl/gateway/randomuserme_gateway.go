package gateway

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/acl/translator"
)

type RandonUserMeGateway struct{}

func NewRandonUserMeGateway() *RandonUserMeGateway {
	return &RandonUserMeGateway{}
}

func (a *RandonUserMeGateway) GetPersonData() (*translator.PersonDto, error) {
	request := fiber.Get("https://randomuser.me/api/")
	_, body, errs := request.Bytes()

	if len(errs) > 0 {
		return nil, nil //&domain.InternalError{Message: "error"}
	}

	var p map[string]any
	err := json.Unmarshal(body, &p)

	result := p["results"].([]any)[0].(map[string]any)

	name := result["name"].(map[string]any)

	if err != nil {
		return nil, nil //&domain.InternalError{Message: "unmarshal error"}
	}

	return &translator.PersonDto{
		Name:      name["first"].(string),
		BirthDate: "12-12-1980",
		Age:       34,
		Email:     result["email"].(string),
	}, nil
}
