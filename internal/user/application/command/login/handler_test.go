// Login command handler tests.
//
// Intent: verify the login use case orchestration in isolation from real infrastructure.
// Objective: cover the RF-01/RF-02 acceptance scenarios — successful login, invalid credentials,
// and lockout after repeated failures — using hand-rolled test doubles for the five domain
// contracts the handler consumes.
package login_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/login"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
)

// stubLoginRepository is a hand-rolled test double implementing contract.LoginRepository.
type stubLoginRepository struct {
	user            *domain.User
	err             error
	calledWith      string
	calledWithCount int
}

func (s *stubLoginRepository) FindByUsernameOrDocument(ctx context.Context, identifier string) (*domain.User, error) {
	s.calledWith = identifier
	s.calledWithCount++
	return s.user, s.err
}

// stubPasswordHasher is a hand-rolled test double implementing contract.PasswordHasher.
type stubPasswordHasher struct {
	verifyResult    bool
	verifyErr       error
	verifyCalledRaw string
}

func (s *stubPasswordHasher) Hash(raw string) (string, error) {
	return "", nil
}

func (s *stubPasswordHasher) Verify(raw, hash string) (bool, error) {
	s.verifyCalledRaw = raw
	return s.verifyResult, s.verifyErr
}

// stubTokenGenerator is a hand-rolled test double implementing contract.TokenGenerator.
type stubTokenGenerator struct {
	accessToken         string
	accessTokenErr      error
	refreshToken        string
	refreshTokenErr     error
	accessCalledUserId  string
	accessCalledGroups  []string
	refreshCalledSessId string
}

func (s *stubTokenGenerator) GenerateAccessToken(userId string, groups []string) (string, error) {
	s.accessCalledUserId = userId
	s.accessCalledGroups = groups
	return s.accessToken, s.accessTokenErr
}

func (s *stubTokenGenerator) GenerateRefreshToken(sessionId string) (string, error) {
	s.refreshCalledSessId = sessionId
	return s.refreshToken, s.refreshTokenErr
}

func (s *stubTokenGenerator) ParseRefreshToken(token string) (string, error) {
	return "", nil
}

// stubLoginAttemptStore is a hand-rolled test double implementing contract.LoginAttemptStore.
type stubLoginAttemptStore struct {
	isBlockedResult      bool
	isBlockedErr         error
	registerFailureErr   error
	registerFailureCalls int
	resetErr             error
	resetCalls           int
}

func (s *stubLoginAttemptStore) RegisterFailure(ctx context.Context, username string) error {
	s.registerFailureCalls++
	return s.registerFailureErr
}

func (s *stubLoginAttemptStore) IsBlocked(ctx context.Context, username string) (bool, error) {
	return s.isBlockedResult, s.isBlockedErr
}

func (s *stubLoginAttemptStore) Reset(ctx context.Context, username string) error {
	s.resetCalls++
	return s.resetErr
}

// stubLoginSessionStore is a hand-rolled test double implementing contract.LoginSessionStore.
type stubLoginSessionStore struct {
	saveErr          error
	saveCalledSessId string
	saveCalledRefreshToken string
}

func (s *stubLoginSessionStore) Save(ctx context.Context, sessionId, refreshToken string) error {
	s.saveCalledSessId = sessionId
	s.saveCalledRefreshToken = refreshToken
	return s.saveErr
}

func validCommand() *login.Command {
	return &login.Command{
		Username: "jdoe",
		Password: "StrongP@ss1",
	}
}

func validUser(t *testing.T) *domain.User {
	t.Helper()
	person, err := domain.NewPerson(&domain.PersonProps{
		UuId:      "person-uuid",
		Name:      "John Doe",
		BirthDate: "1990-01-01",
	})
	if err != nil {
		t.Fatalf("failed to build test person: %v", err)
	}
	group, err := domain.NewUserGroup(&domain.UserGroupProps{
		UuId: "group-uuid",
		Name: "admin",
	})
	if err != nil {
		t.Fatalf("failed to build test user group: %v", err)
	}
	user, err := domain.NewUser(&domain.UserProps{
		UuId:       "user-uuid",
		Username:   "jdoe",
		Password:   domain.NewPasswordFromHash("hashed-password"),
		Person:     person,
		UserGroups: []*domain.UserGroup{group},
	})
	if err != nil {
		t.Fatalf("failed to build test user: %v", err)
	}
	return user
}

func newHandler(
	repository *stubLoginRepository,
	hasher *stubPasswordHasher,
	tokenGenerator *stubTokenGenerator,
	attemptStore *stubLoginAttemptStore,
	sessionStore *stubLoginSessionStore,
) *login.Handler {
	return login.NewCommandHandler(repository, hasher, tokenGenerator, attemptStore, sessionStore)
}

func TestNewCommandHandler(t *testing.T) {
	t.Run("should return a non-nil Handler", func(t *testing.T) {
		t.Parallel()
		handler := newHandler(
			&stubLoginRepository{},
			&stubPasswordHasher{},
			&stubTokenGenerator{},
			&stubLoginAttemptStore{},
			&stubLoginSessionStore{},
		)

		if handler == nil {
			t.Fatal("NewCommandHandler should return a non-nil handler")
		}
	})
}

func TestHandler_Handle(t *testing.T) {
	t.Run("should return an error when checking lockout state fails", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("attempt store unavailable")
		attemptStore := &stubLoginAttemptStore{isBlockedErr: expectedErr}
		handler := newHandler(&stubLoginRepository{}, &stubPasswordHasher{}, &stubTokenGenerator{}, attemptStore, &stubLoginSessionStore{})

		result, err := handler.Handle(context.Background(), validCommand())

		if err != expectedErr {
			t.Fatalf("Handle should return the attempt store error, got: %v", err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should return UserBlockedError when the user is locked out", func(t *testing.T) {
		t.Parallel()
		attemptStore := &stubLoginAttemptStore{isBlockedResult: true}
		handler := newHandler(&stubLoginRepository{}, &stubPasswordHasher{}, &stubTokenGenerator{}, attemptStore, &stubLoginSessionStore{})

		result, err := handler.Handle(context.Background(), validCommand())

		if _, ok := err.(*domain.UserBlockedError); !ok {
			t.Fatalf("Handle should return a *domain.UserBlockedError, got: %T (%v)", err, err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should register a failure and return a validation error when the user is not found", func(t *testing.T) {
		t.Parallel()
		repository := &stubLoginRepository{err: errors.New("not found")}
		attemptStore := &stubLoginAttemptStore{}
		handler := newHandler(repository, &stubPasswordHasher{}, &stubTokenGenerator{}, attemptStore, &stubLoginSessionStore{})

		result, err := handler.Handle(context.Background(), validCommand())

		if err == nil || err.Error() != "Credenciais inválidas" {
			t.Fatalf("Handle should return a 'Credenciais inválidas' validation error, got: %v", err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
		if attemptStore.registerFailureCalls != 1 {
			t.Fatalf("Handle should register exactly one failure, got: %d", attemptStore.registerFailureCalls)
		}
	})

	t.Run("should return an error when registering the failure fails after a not-found user", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("attempt store write failed")
		repository := &stubLoginRepository{err: errors.New("not found")}
		attemptStore := &stubLoginAttemptStore{registerFailureErr: expectedErr}
		handler := newHandler(repository, &stubPasswordHasher{}, &stubTokenGenerator{}, attemptStore, &stubLoginSessionStore{})

		result, err := handler.Handle(context.Background(), validCommand())

		if err != expectedErr {
			t.Fatalf("Handle should return the attempt store error, got: %v", err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should return an error when the password verification fails unexpectedly", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("bcrypt failure")
		repository := &stubLoginRepository{user: validUser(t)}
		hasher := &stubPasswordHasher{verifyErr: expectedErr}
		handler := newHandler(repository, hasher, &stubTokenGenerator{}, &stubLoginAttemptStore{}, &stubLoginSessionStore{})

		result, err := handler.Handle(context.Background(), validCommand())

		if err != expectedErr {
			t.Fatalf("Handle should return the hasher error, got: %v", err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should register a failure and return a validation error when the password does not match", func(t *testing.T) {
		t.Parallel()
		repository := &stubLoginRepository{user: validUser(t)}
		hasher := &stubPasswordHasher{verifyResult: false}
		attemptStore := &stubLoginAttemptStore{}
		handler := newHandler(repository, hasher, &stubTokenGenerator{}, attemptStore, &stubLoginSessionStore{})

		result, err := handler.Handle(context.Background(), validCommand())

		if err == nil || err.Error() != "Credenciais inválidas" {
			t.Fatalf("Handle should return a 'Credenciais inválidas' validation error, got: %v", err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
		if attemptStore.registerFailureCalls != 1 {
			t.Fatalf("Handle should register exactly one failure, got: %d", attemptStore.registerFailureCalls)
		}
	})

	t.Run("should return an error when registering the failure fails after an invalid password", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("attempt store write failed")
		repository := &stubLoginRepository{user: validUser(t)}
		hasher := &stubPasswordHasher{verifyResult: false}
		attemptStore := &stubLoginAttemptStore{registerFailureErr: expectedErr}
		handler := newHandler(repository, hasher, &stubTokenGenerator{}, attemptStore, &stubLoginSessionStore{})

		result, err := handler.Handle(context.Background(), validCommand())

		if err != expectedErr {
			t.Fatalf("Handle should return the attempt store error, got: %v", err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should return an error when resetting the attempt counter fails", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("reset failed")
		repository := &stubLoginRepository{user: validUser(t)}
		hasher := &stubPasswordHasher{verifyResult: true}
		attemptStore := &stubLoginAttemptStore{resetErr: expectedErr}
		handler := newHandler(repository, hasher, &stubTokenGenerator{}, attemptStore, &stubLoginSessionStore{})

		result, err := handler.Handle(context.Background(), validCommand())

		if err != expectedErr {
			t.Fatalf("Handle should return the reset error, got: %v", err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should return an error when issuing the access token fails", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("access token signing failed")
		repository := &stubLoginRepository{user: validUser(t)}
		hasher := &stubPasswordHasher{verifyResult: true}
		tokenGenerator := &stubTokenGenerator{accessTokenErr: expectedErr}
		handler := newHandler(repository, hasher, tokenGenerator, &stubLoginAttemptStore{}, &stubLoginSessionStore{})

		result, err := handler.Handle(context.Background(), validCommand())

		if err != expectedErr {
			t.Fatalf("Handle should return the access token error, got: %v", err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should return an error when issuing the refresh token fails", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("refresh token signing failed")
		repository := &stubLoginRepository{user: validUser(t)}
		hasher := &stubPasswordHasher{verifyResult: true}
		tokenGenerator := &stubTokenGenerator{accessToken: "access-token", refreshTokenErr: expectedErr}
		handler := newHandler(repository, hasher, tokenGenerator, &stubLoginAttemptStore{}, &stubLoginSessionStore{})

		result, err := handler.Handle(context.Background(), validCommand())

		if err != expectedErr {
			t.Fatalf("Handle should return the refresh token error, got: %v", err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should return an error when saving the session fails", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("session save failed")
		repository := &stubLoginRepository{user: validUser(t)}
		hasher := &stubPasswordHasher{verifyResult: true}
		tokenGenerator := &stubTokenGenerator{accessToken: "access-token", refreshToken: "refresh-token"}
		sessionStore := &stubLoginSessionStore{saveErr: expectedErr}
		handler := newHandler(repository, hasher, tokenGenerator, &stubLoginAttemptStore{}, sessionStore)

		result, err := handler.Handle(context.Background(), validCommand())

		if err != expectedErr {
			t.Fatalf("Handle should return the session store error, got: %v", err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should issue a session with access/refresh tokens on success", func(t *testing.T) {
		t.Parallel()
		repository := &stubLoginRepository{user: validUser(t)}
		hasher := &stubPasswordHasher{verifyResult: true}
		tokenGenerator := &stubTokenGenerator{accessToken: "access-token", refreshToken: "refresh-token"}
		attemptStore := &stubLoginAttemptStore{}
		sessionStore := &stubLoginSessionStore{}
		handler := newHandler(repository, hasher, tokenGenerator, attemptStore, sessionStore)
		command := validCommand()

		result, err := handler.Handle(context.Background(), command)

		if err != nil {
			t.Fatalf("Handle should not return an error, got: %v", err)
		}
		payload, ok := result.(map[string]string)
		if !ok {
			t.Fatalf("Handle should return a map[string]string payload, got: %T", result)
		}
		if payload["accessToken"] != "access-token" {
			t.Fatalf("Handle should return the generated access token, got: %q", payload["accessToken"])
		}
		if payload["refreshToken"] != "refresh-token" {
			t.Fatalf("Handle should return the generated refresh token, got: %q", payload["refreshToken"])
		}
		if payload["sessionId"] == "" {
			t.Fatal("Handle should return a non-empty sessionId")
		}
		if attemptStore.resetCalls != 1 {
			t.Fatalf("Handle should reset the attempt counter exactly once, got: %d", attemptStore.resetCalls)
		}
		if sessionStore.saveCalledRefreshToken != "refresh-token" {
			t.Fatalf("Handle should save the session for the authenticated user, got: %q", sessionStore.saveCalledRefreshToken)
		}
		if hasher.verifyCalledRaw != command.Password {
			t.Fatalf("Verify should be called with the raw password %q, got: %q", command.Password, hasher.verifyCalledRaw)
		}
	})
}
