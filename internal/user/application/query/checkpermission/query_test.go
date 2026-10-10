// Package checkpermission_test provides unit tests for the checkpermission query.
// Tests cover the Query struct's Name() method and field accessibility,
// ensuring correct CQRS query behavior for permission checking operations.
package checkpermission_test

import (
	"testing"

	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/query/checkpermission"
)

// TestQuery_Name tests the Name method of the Query struct.
// It verifies that the query returns the correct CQRS routing name
// and that the name is consistent regardless of field values.
func TestQuery_Name(t *testing.T) {
	t.Run("should return the correct query name", func(t *testing.T) {
		t.Parallel()
		query := &checkpermission.Query{
			AccessToken: "test-token",
			Method:      "GET",
			Path:        "/users/123",
		}
		expectedName := "checkPermission"
		actualName := query.Name()
		if actualName != expectedName {
			t.Errorf("Query.Name() = %q, want %q",
				actualName, expectedName)
		}
	})

	t.Run("should return the same name regardless of query fields",
		func(t *testing.T) {
			t.Parallel()
			queries := []*checkpermission.Query{
				{
					AccessToken: "token1",
					Method:      "POST",
					Path:        "/users",
				},
				{
					AccessToken: "token2",
					Method:      "DELETE",
					Path:        "/users/456",
				},
				{
					AccessToken: "",
					Method:      "",
					Path:        "",
				},
			}
			expectedName := "checkPermission"
			for _, q := range queries {
				if q.Name() != expectedName {
					t.Errorf("Query.Name() = %q, want %q",
						q.Name(), expectedName)
				}
			}
		})
}

// TestQuery_Fields tests that the Query struct fields are properly accessible.
// It validates that all fields (AccessToken, Method, Path) can be set and
// retrieved correctly.
func TestQuery_Fields(t *testing.T) {
	t.Run("should have AccessToken field accessible", func(t *testing.T) {
		t.Parallel()
		expectedToken := "test-access-token-123"
		query := &checkpermission.Query{
			AccessToken: expectedToken,
			Method:      "GET",
			Path:        "/api/resource",
		}
		if query.AccessToken != expectedToken {
			t.Errorf("Query.AccessToken = %q, want %q",
				query.AccessToken, expectedToken)
		}
	})

	t.Run("should have Method field accessible", func(t *testing.T) {
		t.Parallel()
		expectedMethod := "POST"
		query := &checkpermission.Query{
			AccessToken: "token",
			Method:      expectedMethod,
			Path:        "/api/resource",
		}
		if query.Method != expectedMethod {
			t.Errorf("Query.Method = %q, want %q",
				query.Method, expectedMethod)
		}
	})

	t.Run("should have Path field accessible", func(t *testing.T) {
		t.Parallel()
		expectedPath := "/users/123/settings"
		query := &checkpermission.Query{
			AccessToken: "token",
			Method:      "GET",
			Path:        expectedPath,
		}
		if query.Path != expectedPath {
			t.Errorf("Query.Path = %q, want %q",
				query.Path, expectedPath)
		}
	})

	t.Run("should allow modification of fields", func(t *testing.T) {
		t.Parallel()
		query := &checkpermission.Query{
			AccessToken: "initial-token",
			Method:      "GET",
			Path:        "/initial",
		}
		// Modify fields
		query.AccessToken = "modified-token"
		query.Method = "PUT"
		query.Path = "/modified"

		if query.AccessToken != "modified-token" {
			t.Errorf(
				"Query.AccessToken after modification = %q, want %q",
				query.AccessToken, "modified-token")
		}
		if query.Method != "PUT" {
			t.Errorf(
				"Query.Method after modification = %q, want %q",
				query.Method, "PUT")
		}
		if query.Path != "/modified" {
			t.Errorf(
				"Query.Path after modification = %q, want %q",
				query.Path, "/modified")
		}
	})

	t.Run("should initialize with zero values when fields are not set",
		func(t *testing.T) {
			t.Parallel()
			query := &checkpermission.Query{}
			if query.AccessToken != "" {
				t.Errorf("Query.AccessToken = %q, want empty string",
					query.AccessToken)
			}
			if query.Method != "" {
				t.Errorf("Query.Method = %q, want empty string",
					query.Method)
			}
			if query.Path != "" {
				t.Errorf("Query.Path = %q, want empty string",
					query.Path)
			}
		})
}

// TestQuery_StructCreation tests various ways of creating Query instances.
// It validates that queries can be created with different field values and
// that the Name() method works correctly across various scenarios.
func TestQuery_StructCreation(t *testing.T) {
	t.Run("should create Query with pointer receiver", func(t *testing.T) {
		t.Parallel()
		query := &checkpermission.Query{
			AccessToken: "test-token",
			Method:      "GET",
			Path:        "/users",
		}
		if query == nil {
			t.Error("Query creation failed, got nil")
		}
		if query.Name() != "checkPermission" {
			t.Errorf("Query.Name() = %q, want %q",
				query.Name(), "checkPermission")
		}
	})

	t.Run("should handle various HTTP methods", func(t *testing.T) {
		t.Parallel()
		methods := []string{
			"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS",
		}
		for _, method := range methods {
			query := &checkpermission.Query{
				AccessToken: "token",
				Method:      method,
				Path:        "/resource",
			}
			if query.Method != method {
				t.Errorf("Query.Method = %q, want %q",
					query.Method, method)
			}
		}
	})

	t.Run("should handle various resource paths", func(t *testing.T) {
		t.Parallel()
		paths := []string{
			"/users",
			"/users/123",
			"/users/123/settings",
			"/api/v1/users",
			"/api/v1/users/123/permissions",
			"/",
		}
		for _, path := range paths {
			query := &checkpermission.Query{
				AccessToken: "token",
				Method:      "GET",
				Path:        path,
			}
			if query.Path != path {
				t.Errorf("Query.Path = %q, want %q",
					query.Path, path)
			}
		}
	})
}
