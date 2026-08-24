// Domain-database mapper.
//
// Intent: act as the anti-corruption layer between the user domain aggregate and its GORM
// persistence models.
// Objective: convert a persisted Users record back into a domain.User aggregate (toDomain), and
// convert a domain.User aggregate into a Users record ready for persistence (toDatabase).
package database

import "github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"

// toDomain reconstitutes a domain.User aggregate from a persisted Users record.
//
// Intent: rehydrate the User aggregate (uuid, username, password hash, person) from storage so it
// can be used by use cases such as login that need to compare credentials.
// Parameters:
//   - user: the persisted Users record, including its associated Person.
//
// Returns: a *domain.User built from the persisted fields. Any construction error from
// domain.NewPerson or domain.NewUser is discarded (nil) since the persisted data is assumed
// already valid.
func toDomain(user *Users) *domain.User {
	person, _ := domain.NewPerson(&domain.PersonProps{
		UuId:      user.Person.Uuid,
		Name:      user.Person.Name,
		BirthDate: user.Person.BirthDate,
	})

	domainUser, _ := domain.NewUser(&domain.UserProps{
		UuId:     user.Uuid,
		Username: user.Username,
		Password: domain.NewPasswordFromHash(user.Password),
		Person:   person,
	})

	return domainUser
}

// toDatabase converts a domain.User aggregate into a Users persistence model.
//
// Intent: map each domain field to its corresponding persistence column via getter methods,
// respecting the aggregate's encapsulation.
// Parameters:
//   - user: the domain.User aggregate to convert.
//
// Returns: a *Users record ready to be persisted.
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
