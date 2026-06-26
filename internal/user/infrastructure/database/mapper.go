package database

import "github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"

func toDomain(user *Users) *domain.User {
	return &domain.User{}
}

func toDatabase(user *domain.User) *Users {
	return &Users{
		Uuid:     user.Uuid(),
		Username: user.Username(),
		Password: user.Password().Value(),
		Person: Person{
			Uuid:      user.Person().Uuid(),
			Name:      user.Person().Name(),
			Document:  user.Person().Document().Value(),
			BirthDate: user.Person().BirthDate(),
		},
	}
}
