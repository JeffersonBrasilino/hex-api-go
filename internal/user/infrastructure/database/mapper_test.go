// Domain-database mapper tests.
//
// Intent: verify the anti-corruption layer between the persisted Users record and the
// domain.User aggregate.
// Objective: cover toDomain, the newly implemented reconstitution path used by the login flow to
// rehydrate a domain.User from storage.
//
// Note: toDomain is package-private, so this file lives in package database (white-box) rather
// than database_test — the standard black-box convention does not apply here since the function
// under test is unexported.
package database

import "testing"

func TestToDomain(t *testing.T) {
	t.Run("should rehydrate a domain.User with the persisted uuid, username and password hash", func(t *testing.T) {
		t.Parallel()
		record := &Users{
			Uuid:     "user-uuid",
			Username: "jdoe",
			Password: "hashed-password",
			Person: Person{
				Uuid:      "person-uuid",
				Name:      "John Doe",
				BirthDate: "1990-01-01",
			},
		}

		user := toDomain(record)

		if user == nil {
			t.Fatal("toDomain should return a non-nil domain.User")
		}
		if user.Uuid() != record.Uuid {
			t.Fatalf("toDomain should preserve the uuid, got: %q, want: %q", user.Uuid(), record.Uuid)
		}
		if user.Username() != record.Username {
			t.Fatalf("toDomain should preserve the username, got: %q, want: %q", user.Username(), record.Username)
		}
		if user.Password() == nil {
			t.Fatal("toDomain should build a non-nil Password")
		}
		if user.Password().Value() != record.Password {
			t.Fatalf("toDomain should preserve the password hash, got: %q, want: %q", user.Password().Value(), record.Password)
		}
	})

	t.Run("should rehydrate the associated Person with its uuid, name and birth date", func(t *testing.T) {
		t.Parallel()
		record := &Users{
			Uuid:     "user-uuid",
			Username: "jdoe",
			Password: "hashed-password",
			Person: Person{
				Uuid:      "person-uuid",
				Name:      "John Doe",
				BirthDate: "1990-01-01",
			},
		}

		user := toDomain(record)

		if user.Person() == nil {
			t.Fatal("toDomain should build a non-nil Person")
		}
		if user.Person().Uuid() != record.Person.Uuid {
			t.Fatalf("toDomain should preserve the person uuid, got: %q, want: %q", user.Person().Uuid(), record.Person.Uuid)
		}
		if user.Person().Name() != record.Person.Name {
			t.Fatalf("toDomain should preserve the person name, got: %q, want: %q", user.Person().Name(), record.Person.Name)
		}
		if user.Person().BirthDate() != record.Person.BirthDate {
			t.Fatalf("toDomain should preserve the person birth date, got: %q, want: %q", user.Person().BirthDate(), record.Person.BirthDate)
		}
	})
}
