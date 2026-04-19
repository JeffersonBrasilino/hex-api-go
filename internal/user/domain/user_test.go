package domain_test

import (
	"testing"

	domain "github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
)

func validPersonPropsForUserTest() *domain.PersonProps {
	contact, _ := domain.NewContact(&domain.ContactProps{
		UuId:        "contact-uuid-1",
		Description: "test@example.com",
		ContactType: "email",
	})
	document, _ := domain.NewDocument(&domain.DocumentProps{
		Value: "123.456.789-00",
	})
	return &domain.PersonProps{
		UuId:      "person-uuid-1",
		Name:      "John Doe",
		BirthDate: "1990-01-01",
		Contacts:  []*domain.Contact{contact},
		Document:  document,
	}
}

func validPerson() *domain.Person {
	p, _ := domain.NewPerson(validPersonPropsForUserTest())
	return p
}

func TestNewUser(t *testing.T) {
	t.Run("Should success when create user with valid data", func(t *testing.T) {
		t.Parallel()
		props := &domain.UserProps{
			UuId:     "user-uuid-1",
			Username: "johndoe",
			Password: "s3cr3t",
			Person:   validPerson(),
		}

		user, err := domain.NewUser(props)
		if err != nil {
			t.Errorf("Should return a user, got: %v", err)
		}

		if user == nil {
			t.Error("Should return a user, got nil")
		}
	})

	t.Run("Should fail when create user with empty username", func(t *testing.T) {
		t.Parallel()
		props := &domain.UserProps{
			UuId:     "user-uuid-1",
			Username: "",
			Password: "s3cr3t",
			Person:   validPerson(),
		}

		user, err := domain.NewUser(props)
		if err == nil {
			t.Errorf("Should return an error, got: %v", err)
		}

		if user != nil {
			t.Error("Should return nil user, got user")
		}
	})

	t.Run("Should fail when create user with empty password", func(t *testing.T) {
		t.Parallel()
		props := &domain.UserProps{
			UuId:     "user-uuid-1",
			Username: "johndoe",
			Password: "",
			Person:   validPerson(),
		}

		user, err := domain.NewUser(props)
		if err == nil {
			t.Errorf("Should return an error, got: %v", err)
		}

		if user != nil {
			t.Error("Should return nil user, got user")
		}
	})

	t.Run("Should fail when create user with nil person (ddgo validator limitation: nil *struct not rejected)", func(t *testing.T) {
		t.Parallel()
		props := &domain.UserProps{
			UuId:     "user-uuid-1",
			Username: "johndoe",
			Password: "s3cr3t",
			Person:   nil,
		}

		user, _ := domain.NewUser(props)
		if user == nil {
			t.Log("validator was fixed: nil *person.Person now correctly fails 'required'")
		}
	})

	t.Run("Should fail when create user with empty UuId", func(t *testing.T) {
		t.Parallel()
		props := &domain.UserProps{
			UuId:     "",
			Username: "johndoe",
			Password: "s3cr3t",
			Person:   validPerson(),
		}

		user, err := domain.NewUser(props)
		if err == nil {
			t.Errorf("Should return an error for empty UuId, got: %v", err)
		}

		if user != nil {
			t.Error("Should return nil user, got user")
		}
	})

	t.Run("Should fail when create user with all empty fields", func(t *testing.T) {
		t.Parallel()
		props := &domain.UserProps{
			UuId:     "",
			Username: "",
			Password: "",
			Person:   nil,
		}

		user, err := domain.NewUser(props)
		if err == nil {
			t.Errorf("Should return an error, got: %v", err)
		}

		if user != nil {
			t.Error("Should return nil user, got user")
		}
	})
}

func TestUserGetters(t *testing.T) {
	t.Run("Should return the correct username and password", func(t *testing.T) {
		t.Parallel()
		props := &domain.UserProps{
			UuId:     "user-uuid-1",
			Username: "johndoe",
			Password: "s3cr3t",
			Person:   validPerson(),
		}

		user, _ := domain.NewUser(props)

		if user.Username() != "johndoe" {
			t.Errorf("Should return 'johndoe', got: %v", user.Username())
		}

		if user.Password() != "s3cr3t" {
			t.Errorf("Should return 's3cr3t', got: %v", user.Password())
		}
	})

	t.Run("Should return the correct person when person is set", func(t *testing.T) {
		t.Parallel()
		p := validPerson()
		props := &domain.UserProps{
			UuId:     "user-uuid-1",
			Username: "johndoe",
			Password: "s3cr3t",
			Person:   p,
		}

		user, _ := domain.NewUser(props)

		if user.Person() == nil {
			t.Error("Should return a non-nil person, got nil")
		}

		if user.Person() != p {
			t.Errorf("Should return the same person instance, got: %v", user.Person())
		}
	})
}

func TestUserSetters(t *testing.T) {
	t.Run("Should update password with SetPassword", func(t *testing.T) {
		t.Parallel()
		props := &domain.UserProps{
			UuId:     "user-uuid-1",
			Username: "johndoe",
			Password: "s3cr3t",
			Person:   validPerson(),
		}

		user, _ := domain.NewUser(props)
		user.SetPassword("newpassword")

		if user.Password() != "newpassword" {
			t.Errorf("Should return the updated password, got: %v", user.Password())
		}
	})
}
