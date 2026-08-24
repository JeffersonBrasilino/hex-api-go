// GormUserRepository tests.
//
// Intent: cover FindByUsernameOrDocument, the repository method added to back the login flow
// (RF-01), without touching a real database.
// Objective: exercise the found-by-username, found-by-document, not-found and raw-error paths by
// mocking the underlying database/sql driver GORM sits on top of via go-sqlmock, and asserting the
// SQL statement GORM emits, the mapped domain.User result and the ddgo error taxonomy on failure.
//
// Convention: this project has no prior sqlmock convention, so this file establishes it —
// sqlmock.New() opens a mocked database/sql.DB, gorm.Open(postgres.New(postgres.Config{Conn: ...}))
// wires GORM on top of it (no real network connection is made), and QueryMatcherRegexp mode is used
// so expectations match the shape of the emitted SQL rather than the exact column list, which keeps
// the tests readable and resilient to unrelated column reordering.
package database_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jeffersonbrasilino/ddgo"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// fixedTimestamp is used to populate created_at/updated_at columns in the mocked rows: GORM scans
// these into time.Time, so a driver.Value-compatible time is required (a plain string fails to scan).
var fixedTimestamp = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

// userRows are the columns GORM selects for the Users/Person join used by
// FindByUsernameOrDocument.
var userRows = []string{
	"id", "created_at", "updated_at", "deleted_at",
	"uuid", "username", "password", "verification_code", "person_id", "status",
	"Person__id", "Person__created_at", "Person__updated_at", "Person__deleted_at",
	"Person__uuid", "Person__name", "Person__document", "Person__birth_date", "Person__status",
}

// newMockedRepository opens GORM against a sqlmock-controlled database/sql.DB so
// FindByUsernameOrDocument can be exercised without a real database connection.
func newMockedRepository(t *testing.T) (*database.GormUserRepository, sqlmock.Sqlmock) {
	t.Helper()

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock.New should not return an error, got: %v", err)
	}
	t.Cleanup(func() {
		_ = mockDB.Close()
	})

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open should not return an error, got: %v", err)
	}

	return database.NewGormUserRepository(gdb), mock
}

func TestGormUserRepository_FindByUsernameOrDocument(t *testing.T) {
	const joinQuery = `SELECT .+ FROM "hex-api-go"\."users" INNER JOIN "hex-api-go"\."persons" "Person" ON .+ WHERE \(username = \$1 OR document = \$2\).*`

	t.Run("should return the domain.User when found by username", func(t *testing.T) {
		t.Parallel()
		repo, mock := newMockedRepository(t)
		identifier := "jdoe"
		rows := sqlmock.NewRows(userRows).AddRow(
			1, fixedTimestamp, fixedTimestamp, nil,
			"user-uuid", "jdoe", "hashed-password", "", 1, 1,
			1, fixedTimestamp, fixedTimestamp, nil,
			"person-uuid", "John Doe", "12345678900", "1990-01-01", 1,
		)
		mock.ExpectQuery(joinQuery).WithArgs(identifier, identifier, 1).WillReturnRows(rows)

		user, err := repo.FindByUsernameOrDocument(t.Context(), identifier)

		if err != nil {
			t.Fatalf("FindByUsernameOrDocument should not return an error, got: %v", err)
		}
		if user == nil {
			t.Fatal("FindByUsernameOrDocument should return a non-nil domain.User")
		}
		if user.Uuid() != "user-uuid" {
			t.Errorf("FindByUsernameOrDocument should preserve the uuid, got: %q", user.Uuid())
		}
		if user.Username() != "jdoe" {
			t.Errorf("FindByUsernameOrDocument should preserve the username, got: %q", user.Username())
		}
		if user.Password().Value() != "hashed-password" {
			t.Errorf("FindByUsernameOrDocument should preserve the password hash, got: %q", user.Password().Value())
		}
		if user.Person() == nil || user.Person().Name() != "John Doe" {
			t.Errorf("FindByUsernameOrDocument should preserve the joined Person, got: %+v", user.Person())
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet sqlmock expectations: %v", err)
		}
	})

	t.Run("should return the domain.User when found by document", func(t *testing.T) {
		t.Parallel()
		repo, mock := newMockedRepository(t)
		identifier := "12345678900"
		rows := sqlmock.NewRows(userRows).AddRow(
			2, fixedTimestamp, fixedTimestamp, nil,
			"user-uuid-2", "asilva", "hashed-password-2", "", 2, 1,
			2, fixedTimestamp, fixedTimestamp, nil,
			"person-uuid-2", "Ana Silva", identifier, "1985-05-05", 1,
		)
		mock.ExpectQuery(joinQuery).WithArgs(identifier, identifier, 1).WillReturnRows(rows)

		user, err := repo.FindByUsernameOrDocument(t.Context(), identifier)

		if err != nil {
			t.Fatalf("FindByUsernameOrDocument should not return an error, got: %v", err)
		}
		if user == nil {
			t.Fatal("FindByUsernameOrDocument should return a non-nil domain.User")
		}
		if user.Uuid() != "user-uuid-2" {
			t.Errorf("FindByUsernameOrDocument should preserve the uuid, got: %q", user.Uuid())
		}
		if user.Person() == nil || user.Person().Name() != "Ana Silva" {
			t.Errorf("FindByUsernameOrDocument should preserve the joined Person, got: %+v", user.Person())
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet sqlmock expectations: %v", err)
		}
	})

	t.Run("should return a ddgo.NotFoundError when no user matches the identifier", func(t *testing.T) {
		t.Parallel()
		repo, mock := newMockedRepository(t)
		identifier := "unknown"
		mock.ExpectQuery(joinQuery).
			WithArgs(identifier, identifier, 1).
			WillReturnRows(sqlmock.NewRows(userRows))

		user, err := repo.FindByUsernameOrDocument(t.Context(), identifier)

		if user != nil {
			t.Errorf("FindByUsernameOrDocument should return a nil user on not-found, got: %+v", user)
		}
		var notFoundErr *ddgo.NotFoundError
		if !errors.As(err, &notFoundErr) {
			t.Fatalf("FindByUsernameOrDocument should return a *ddgo.NotFoundError, got: %T (%v)", err, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet sqlmock expectations: %v", err)
		}
	})

	t.Run("should return a ddgo.InternalError when the query fails", func(t *testing.T) {
		t.Parallel()
		repo, mock := newMockedRepository(t)
		identifier := "jdoe"
		mock.ExpectQuery(joinQuery).
			WithArgs(identifier, identifier, 1).
			WillReturnError(errors.New("connection reset by peer"))

		user, err := repo.FindByUsernameOrDocument(t.Context(), identifier)

		if user != nil {
			t.Errorf("FindByUsernameOrDocument should return a nil user on query error, got: %+v", user)
		}
		var internalErr *ddgo.InternalError
		if !errors.As(err, &internalErr) {
			t.Fatalf("FindByUsernameOrDocument should return a *ddgo.InternalError, got: %T (%v)", err, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet sqlmock expectations: %v", err)
		}
	})
}
