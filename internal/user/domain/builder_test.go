package domain_test

import (
	"testing"

	domain "github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
)

func validPersonProps() *domain.WithPersonProps {
	return &domain.WithPersonProps{
		UuId:      "person-uuid-1",
		Name:      "John Doe",
		BirthDate: "1990-01-01",
		Document:  "123.456.789-00",
		Contacts: []*domain.WithContactProps{
			{
				UuId:        "contact-uuid-1",
				Description: "test@example.com",
				Main:        true,
				ContactType: "contact-type-uuid-1",
			},
		},
	}
}

func TestNewBuilder(t *testing.T) {
	t.Run("Should return a non-nil Builder instance", func(t *testing.T) {
		t.Parallel()
		b := domain.NewBuilder()
		if b == nil {
			t.Error("NewBuilder() should return a non-nil *Builder, got nil")
		}
	})
}

func TestBuilder_WithUuId(t *testing.T) {
	t.Run("Should set uuId and return the same builder (fluent interface)", func(t *testing.T) {
		t.Parallel()
		b := domain.NewBuilder()
		returned := b.WithUuId("some-uuid")
		if returned == nil {
			t.Error("WithUuId() should return a non-nil *Builder, got nil")
		}
		if returned != b {
			t.Error("WithUuId() should return the same *Builder instance for method chaining")
		}
	})
}

func TestBuilder_WithUsername(t *testing.T) {
	t.Run("Should set username and return the same builder (fluent interface)", func(t *testing.T) {
		t.Parallel()
		b := domain.NewBuilder()
		returned := b.WithUsername("johndoe")
		if returned == nil {
			t.Error("WithUsername() should return a non-nil *Builder, got nil")
		}
		if returned != b {
			t.Error("WithUsername() should return the same *Builder instance for method chaining")
		}
	})
}

func TestBuilder_WithPassword(t *testing.T) {
	t.Run("Should set password and return the same builder (fluent interface)", func(t *testing.T) {
		t.Parallel()
		b := domain.NewBuilder()
		returned := b.WithPassword("StrongP@ss123")
		if returned == nil {
			t.Error("WithPassword() should return a non-nil *Builder, got nil")
		}
		if returned != b {
			t.Error("WithPassword() should return the same *Builder instance for method chaining")
		}
	})
}

func TestBuilder_WithPerson(t *testing.T) {
	t.Run("Should set person successfully when props are valid", func(t *testing.T) {
		t.Parallel()
		b := domain.NewBuilder()
		returned := b.WithPerson(validPersonProps())
		if returned == nil {
			t.Error("WithPerson() should return a non-nil *Builder on valid props, got nil")
		}
		if returned != b {
			t.Error("WithPerson() should return the same *Builder instance for method chaining")
		}
	})

	t.Run("Should append error and return builder when PersonProps are invalid", func(t *testing.T) {
		t.Parallel()
		b := domain.NewBuilder()
		invalidProps := &domain.WithPersonProps{}
		returned := b.WithPerson(invalidProps)
		if returned == nil {
			t.Error("WithPerson() should return a non-nil *Builder even on invalid props, got nil")
		}
		if returned != b {
			t.Error("WithPerson() should return the same *Builder instance for method chaining")
		}
	})
}

func TestBuilder_Build(t *testing.T) {
	t.Run("Should succeed and return a User when all fields are valid", func(t *testing.T) {
		t.Parallel()
		user, err := domain.NewBuilder().
			WithUuId("user-uuid-1").
			WithUsername("johndoe").
			WithPassword("StrongP@ss123").
			WithPerson(validPersonProps()).
			Build()

		if err != nil {
			t.Errorf("Build() should succeed with valid data, got error: %v", err)
		}
		if user == nil {
			t.Error("Build() should return a non-nil User, got nil")
		}
	})

	t.Run("Should fail when WithPerson received invalid props", func(t *testing.T) {
		t.Parallel()
		user, err := domain.NewBuilder().
			WithUuId("user-uuid-1").
			WithUsername("johndoe").
			WithPassword("StrongP@ss123").
			WithPerson(&domain.WithPersonProps{}).
			Build()

		if err == nil {
			t.Error("Build() should return an error when WithPerson received invalid props")
		}
		if user != nil {
			t.Errorf("Build() should return nil user on error, got: %v", user)
		}
	})

	t.Run("Should fail when UserProps are invalid after Build()", func(t *testing.T) {
		t.Parallel()
		user, err := domain.NewBuilder().
			WithPassword("StrongP@ss123").
			WithPerson(validPersonProps()).
			Build()

		if err == nil {
			t.Error("Build() should return an error when required UserProps fields are missing")
		}
		if user != nil {
			t.Errorf("Build() should return nil user on error, got: %v", user)
		}
	})

	t.Run("Should support method chaining across all With* setters", func(t *testing.T) {
		t.Parallel()
		b := domain.NewBuilder().
			WithUuId("u1").
			WithUsername("jane").
			WithPassword("StrongP@ss123").
			WithPerson(validPersonProps())

		if b == nil {
			t.Error("Method chain should return a non-nil *Builder")
		}
	})
}
