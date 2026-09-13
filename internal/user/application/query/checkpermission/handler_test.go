// Handler tests for permission checking.
//
// Intent: verify the checkpermission query handler's orchestration of token parsing, role
// querying, and access decisions in isolation from real infrastructure.
// Objective: cover the RF-01/RF-02/RF-03/RF-04 acceptance scenarios — invalid token,
// public route, denied access, allowed access, and infrastructure errors — using
// hand-rolled test doubles for AccessTokenValidator and PermissionRepository.
package checkpermission_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/query/checkpermission"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
)

// stubAccessTokenValidator is a hand-rolled test double implementing
// contract.AccessTokenValidator.
type stubAccessTokenValidator struct {
	userId        string
	groups        []string
	err           error
	calledWithArg string
	calledCount   int
}

func (s *stubAccessTokenValidator) ParseAccessToken(token string) (string, []string, error) {
	s.calledWithArg = token
	s.calledCount++
	return s.userId, s.groups, s.err
}

// stubPermissionRepository is a hand-rolled test double implementing
// contract.PermissionRepository.
type stubPermissionRepository struct {
	roles           []string
	err             error
	calledWithArg   string
	calledWithPath  string
	calledWithCount int
}

func (s *stubPermissionRepository) RolesWithAccess(
	ctx context.Context, method, path string,
) ([]string, error) {
	s.calledWithArg = method
	s.calledWithPath = path
	s.calledWithCount++
	return s.roles, s.err
}

// TestQueryHandler_Handle tests the handler's permission check logic.
func TestQueryHandler_Handle(t *testing.T) {
	t.Run("invalid or expired token returns InvalidSessionError", func(t *testing.T) {
		t.Parallel()

		tokenValidator := &stubAccessTokenValidator{
			err: errors.New("token expired"),
		}
		permissionRepo := &stubPermissionRepository{
			roles: []string{"admin"},
		}

		handler := checkpermission.NewQueryHandler(tokenValidator, permissionRepo)
		query := &checkpermission.Query{
			AccessToken: "invalid-token",
			Method:      "GET",
			Path:        "/users",
		}

		result, err := handler.Handle(context.Background(), query)

		if result != nil {
			t.Errorf("expected nil result, got %v", result)
		}

		var invalidSessionErr *domain.InvalidSessionError
		if !errors.As(err, &invalidSessionErr) {
			t.Errorf("expected InvalidSessionError, got %v (%T)", err, err)
		}

		if tokenValidator.calledWithArg != "invalid-token" {
			t.Errorf("expected token validator called with 'invalid-token', got '%s'",
				tokenValidator.calledWithArg)
		}
	})

	t.Run("route with no permission mappings allows access (public)", func(t *testing.T) {
		t.Parallel()

		tokenValidator := &stubAccessTokenValidator{
			userId: "user-123",
			groups: []string{"staff"},
		}
		permissionRepo := &stubPermissionRepository{
			roles: []string{}, // Empty roles = public route
		}

		handler := checkpermission.NewQueryHandler(tokenValidator, permissionRepo)
		query := &checkpermission.Query{
			AccessToken: "valid-token",
			Method:      "GET",
			Path:        "/public",
		}

		result, err := handler.Handle(context.Background(), query)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if result != true {
			t.Errorf("expected true, got %v", result)
		}

		if permissionRepo.calledWithArg != "GET" {
			t.Errorf("expected permission repo called with method 'GET', got '%s'",
				permissionRepo.calledWithArg)
		}

		if permissionRepo.calledWithPath != "/public" {
			t.Errorf("expected permission repo called with path '/public', got '%s'",
				permissionRepo.calledWithPath)
		}
	})

	t.Run("route with required roles, user missing role returns AccessDeniedError", func(t *testing.T) {
		t.Parallel()

		tokenValidator := &stubAccessTokenValidator{
			userId: "user-123",
			groups: []string{"staff"}, // User has 'staff' role
		}
		permissionRepo := &stubPermissionRepository{
			roles: []string{"admin", "manager"}, // Route requires 'admin' or 'manager'
		}

		handler := checkpermission.NewQueryHandler(tokenValidator, permissionRepo)
		query := &checkpermission.Query{
			AccessToken: "valid-token",
			Method:      "DELETE",
			Path:        "/users/123",
		}

		result, err := handler.Handle(context.Background(), query)

		if result != nil {
			t.Errorf("expected nil result, got %v", result)
		}

		var accessDeniedErr *domain.AccessDeniedError
		if !errors.As(err, &accessDeniedErr) {
			t.Errorf("expected AccessDeniedError, got %v (%T)", err, err)
		}
	})

	t.Run("route with required roles, user has matching role allows access", func(t *testing.T) {
		t.Parallel()

		tokenValidator := &stubAccessTokenValidator{
			userId: "user-123",
			groups: []string{"staff", "manager"}, // User has both 'staff' and 'manager'
		}
		permissionRepo := &stubPermissionRepository{
			roles: []string{"admin", "manager"}, // Route requires 'admin' or 'manager'
		}

		handler := checkpermission.NewQueryHandler(tokenValidator, permissionRepo)
		query := &checkpermission.Query{
			AccessToken: "valid-token",
			Method:      "PUT",
			Path:        "/users/123",
		}

		result, err := handler.Handle(context.Background(), query)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if result != true {
			t.Errorf("expected true, got %v", result)
		}
	})

	t.Run("permission repository error propagates (fail-closed)", func(t *testing.T) {
		t.Parallel()

		repoErr := errors.New("database connection error")
		tokenValidator := &stubAccessTokenValidator{
			userId: "user-123",
			groups: []string{"user"},
		}
		permissionRepo := &stubPermissionRepository{
			err: repoErr,
		}

		handler := checkpermission.NewQueryHandler(tokenValidator, permissionRepo)
		query := &checkpermission.Query{
			AccessToken: "valid-token",
			Method:      "GET",
			Path:        "/protected",
		}

		result, err := handler.Handle(context.Background(), query)

		if result != nil {
			t.Errorf("expected nil result on error, got %v", result)
		}

		if !errors.Is(err, repoErr) {
			t.Errorf("expected repository error to propagate, got %v", err)
		}
	})

	t.Run("multiple user groups, one matches route role", func(t *testing.T) {
		t.Parallel()

		tokenValidator := &stubAccessTokenValidator{
			userId: "user-456",
			groups: []string{"viewer", "editor", "commenter"},
		}
		permissionRepo := &stubPermissionRepository{
			roles: []string{"admin", "editor"},
		}

		handler := checkpermission.NewQueryHandler(tokenValidator, permissionRepo)
		query := &checkpermission.Query{
			AccessToken: "valid-token",
			Method:      "POST",
			Path:        "/content",
		}

		result, err := handler.Handle(context.Background(), query)

		if err != nil {
			t.Errorf("expected no error when user has matching role, got %v", err)
		}

		if result != true {
			t.Errorf("expected true, got %v", result)
		}
	})

	t.Run("public route allows access without a token, never calling the token validator", func(t *testing.T) {
		t.Parallel()

		tokenValidator := &stubAccessTokenValidator{}
		permissionRepo := &stubPermissionRepository{
			roles: []string{}, // Empty roles = public route
		}

		handler := checkpermission.NewQueryHandler(tokenValidator, permissionRepo)
		query := &checkpermission.Query{
			AccessToken: "", // No token at all — e.g. an unauthenticated login request.
			Method:      "POST",
			Path:        "/users/login",
		}

		result, err := handler.Handle(context.Background(), query)

		if err != nil {
			t.Errorf("expected no error for a public route with no token, got %v", err)
		}

		if result != true {
			t.Errorf("expected true, got %v", result)
		}

		if tokenValidator.calledCount != 0 {
			t.Errorf("expected the token validator to never be called for a public route, called %d times",
				tokenValidator.calledCount)
		}
	})

	t.Run("protected route with no token returns InvalidSessionError", func(t *testing.T) {
		t.Parallel()

		tokenValidator := &stubAccessTokenValidator{
			err: errors.New("empty token"),
		}
		permissionRepo := &stubPermissionRepository{
			roles: []string{"admin"},
		}

		handler := checkpermission.NewQueryHandler(tokenValidator, permissionRepo)
		query := &checkpermission.Query{
			AccessToken: "",
			Method:      "DELETE",
			Path:        "/users/123",
		}

		result, err := handler.Handle(context.Background(), query)

		if result != nil {
			t.Errorf("expected nil result, got %v", result)
		}

		var invalidSessionErr *domain.InvalidSessionError
		if !errors.As(err, &invalidSessionErr) {
			t.Errorf("expected InvalidSessionError, got %v (%T)", err, err)
		}
	})

	t.Run("user with empty groups denied access to protected route", func(t *testing.T) {
		t.Parallel()

		tokenValidator := &stubAccessTokenValidator{
			userId: "user-789",
			groups: []string{}, // No groups
		}
		permissionRepo := &stubPermissionRepository{
			roles: []string{"admin"}, // Route requires admin
		}

		handler := checkpermission.NewQueryHandler(tokenValidator, permissionRepo)
		query := &checkpermission.Query{
			AccessToken: "valid-token",
			Method:      "GET",
			Path:        "/admin/dashboard",
		}

		result, err := handler.Handle(context.Background(), query)

		if result != nil {
			t.Errorf("expected nil result, got %v", result)
		}

		var accessDeniedErr *domain.AccessDeniedError
		if !errors.As(err, &accessDeniedErr) {
			t.Errorf("expected AccessDeniedError, got %v (%T)", err, err)
		}
	})
}
