// User registration command handler tests.
//
// Intent: verify the createUser use case orchestration in isolation from real infrastructure.
// Objective: cover document uniqueness checks, aggregate build failures, password hashing, and
// repository persistence, using hand-rolled test doubles for the domain contracts.
package createuser_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/createuser"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
)

// stubUserRepository is a hand-rolled test double implementing contract.UserRepository.
type stubUserRepository struct {
	existsByDocumentResult bool
	existsByDocumentErr    error
	createErr              error
	createCalledWith       *domain.User
}

func (s *stubUserRepository) ExistsByDocument(ctx context.Context, document string) (bool, error) {
	return s.existsByDocumentResult, s.existsByDocumentErr
}

func (s *stubUserRepository) Create(ctx context.Context, aggregate *domain.User) error {
	s.createCalledWith = aggregate
	return s.createErr
}

// stubPasswordHasher is a hand-rolled test double implementing contract.PasswordHasher.
type stubPasswordHasher struct {
	hashResult     string
	hashErr        error
	hashCalledWith string
}

func (s *stubPasswordHasher) Hash(raw string) (string, error) {
	s.hashCalledWith = raw
	return s.hashResult, s.hashErr
}

func (s *stubPasswordHasher) Verify(raw, hash string) (bool, error) {
	return false, nil
}

func validCommand() *createuser.Command {
	return &createuser.Command{
		Username:   "jdoe",
		Password:   "StrongP@ss1",
		PersonName: "John Doe",
		Document:   "52998224725",
		BirthDate:  "1990-01-01",
		Email:      "jdoe@example.com",
	}
}

func TestNewComandHandler(t *testing.T) {
	t.Run("should return a non-nil Handler", func(t *testing.T) {
		t.Parallel()
		repository := &stubUserRepository{}
		hasher := &stubPasswordHasher{}

		handler := createuser.NewComandHandler(repository, hasher)

		if handler == nil {
			t.Fatal("NewComandHandler should return a non-nil handler")
		}
	})
}

func TestHandler_Handle(t *testing.T) {
	t.Run("should return an error when ExistsByDocument fails", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("datasource unavailable")
		repository := &stubUserRepository{existsByDocumentErr: expectedErr}
		hasher := &stubPasswordHasher{}
		handler := createuser.NewComandHandler(repository, hasher)

		result, err := handler.Handle(context.Background(), validCommand())

		if err != expectedErr {
			t.Fatalf("Handle should return the repository error, got: %v", err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should return AlreadyExistsError when document is already registered", func(t *testing.T) {
		t.Parallel()
		repository := &stubUserRepository{existsByDocumentResult: true}
		hasher := &stubPasswordHasher{}
		handler := createuser.NewComandHandler(repository, hasher)

		result, err := handler.Handle(context.Background(), validCommand())

		if err == nil {
			t.Fatal("Handle should return an error when the document already exists")
		}
		if err.Error() != "User already exists" {
			t.Fatalf("Handle should return an AlreadyExistsError, got: %v", err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should return an error when the aggregate fails to build", func(t *testing.T) {
		t.Parallel()
		repository := &stubUserRepository{existsByDocumentResult: false}
		hasher := &stubPasswordHasher{}
		handler := createuser.NewComandHandler(repository, hasher)
		invalidCommand := validCommand()
		invalidCommand.Username = ""

		result, err := handler.Handle(context.Background(), invalidCommand)

		if err == nil {
			t.Fatal("Handle should return an error when the aggregate build fails")
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should return an error when the password hasher fails", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("hashing failed")
		repository := &stubUserRepository{existsByDocumentResult: false}
		hasher := &stubPasswordHasher{hashErr: expectedErr}
		handler := createuser.NewComandHandler(repository, hasher)

		result, err := handler.Handle(context.Background(), validCommand())

		if err != expectedErr {
			t.Fatalf("Handle should return the hasher error, got: %v", err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should return an error when repository Create fails", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("insert failed")
		repository := &stubUserRepository{existsByDocumentResult: false, createErr: expectedErr}
		hasher := &stubPasswordHasher{hashResult: "hashed-password"}
		handler := createuser.NewComandHandler(repository, hasher)

		result, err := handler.Handle(context.Background(), validCommand())

		if err != expectedErr {
			t.Fatalf("Handle should return the repository Create error, got: %v", err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should hash the password and persist the user on success", func(t *testing.T) {
		t.Parallel()
		repository := &stubUserRepository{existsByDocumentResult: false}
		hasher := &stubPasswordHasher{hashResult: "hashed-password"}
		handler := createuser.NewComandHandler(repository, hasher)
		command := validCommand()

		result, err := handler.Handle(context.Background(), command)

		if err != nil {
			t.Fatalf("Handle should not return an error, got: %v", err)
		}
		if result != "okok" {
			t.Fatalf("Handle should return %q, got: %v", "okok", result)
		}
		if hasher.hashCalledWith != command.Password {
			t.Fatalf("Hash should be called with the raw password %q, got: %q", command.Password, hasher.hashCalledWith)
		}
		if repository.createCalledWith == nil {
			t.Fatal("Create should be called with the built aggregate")
		}
	})

	t.Run("should build the aggregate without a contact when Email is empty", func(t *testing.T) {
		t.Parallel()
		repository := &stubUserRepository{existsByDocumentResult: false}
		hasher := &stubPasswordHasher{hashResult: "hashed-password"}
		handler := createuser.NewComandHandler(repository, hasher)
		command := validCommand()
		command.Email = ""

		result, err := handler.Handle(context.Background(), command)

		if err != nil {
			t.Fatalf("Handle should not return an error, got: %v", err)
		}
		if result != "okok" {
			t.Fatalf("Handle should return %q, got: %v", "okok", result)
		}
	})
}
