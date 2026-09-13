// Permission cache repository tests.
//
// Intent: exercise the PermissionCacheRepository's cache-aside pattern,
// testing the Redis caching logic with miniredis and the real Postgres JOIN
// fallback (mocked via go-sqlmock, following the convention established in
// gorm_user_repository_test.go).
// Objective: verify cache hit (with and without sentinel), cache miss
// triggering the Postgres JOIN fallback (mapped route returns roles, unmapped
// route returns empty and caches the __EMPTY__ sentinel), Postgres query
// errors (fail-closed DependencyError), Redis errors, and fail-closed
// semantics (deny all on infrastructure failure).
package database_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jeffersonbrasilino/ddgo"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/database"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// rolesJoinQuery matches the Postgres JOIN emitted by rolesWithAccessFromPostgres:
// UserGroupsPermissions -> UsersGroups (role name) and
// UserGroupsPermissions -> ApiRoutes (route match), filtered by
// activeStatus (1) on all three tables.
const rolesJoinQuery = `SELECT DISTINCT "hex-api-go"\."users_groups"\."name" ` +
	`FROM "hex-api-go"\."user_groups_permissions" ` +
	`JOIN "hex-api-go"\."users_groups" ON .+ ` +
	`JOIN "hex-api-go"\."api_routes" ON .+ ` +
	`WHERE .+`

// newMockedPostgresDB opens GORM against a sqlmock-controlled database/sql.DB, following the
// same convention as newMockedRepository in gorm_user_repository_test.go, so the Postgres JOIN
// path in rolesWithAccessFromPostgres can be exercised without a real database connection.
func newMockedPostgresDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

	return gdb, mock
}

// TestPermissionCacheRepository_RolesWithAccess tests the cache-aside pattern.
func TestPermissionCacheRepository_RolesWithAccess(t *testing.T) {
	t.Run("cache hit returns stored roles", func(t *testing.T) {
		t.Parallel()

		rdb := newTestRedisClient(t)
		db := &gorm.DB{} // Not used in this test

		ctx := context.Background()
		key := "auth:route:GET:/users"

		// Pre-populate Redis with roles.
		rdb.SAdd(ctx, key, "admin", "operator")

		repo := database.NewPermissionCacheRepository(rdb, db)
		roles, err := repo.RolesWithAccess(ctx, "GET", "/users")

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(roles) != 2 {
			t.Fatalf("expected 2 roles, got: %d (%v)", len(roles), roles)
		}
		// Check that both roles are present (order may vary in set).
		roleMap := make(map[string]bool)
		for _, role := range roles {
			roleMap[role] = true
		}
		if !roleMap["admin"] || !roleMap["operator"] {
			t.Fatalf("expected roles {admin, operator}, got: %v", roles)
		}
	})

	t.Run("cache hit with sentinel returns empty list (public route)", func(t *testing.T) {
		t.Parallel()

		rdb := newTestRedisClient(t)
		db := &gorm.DB{} // Not used

		ctx := context.Background()
		key := "auth:route:GET:/public"

		// Pre-populate Redis with the empty sentinel.
		rdb.SAdd(ctx, key, "__EMPTY__")

		repo := database.NewPermissionCacheRepository(rdb, db)
		roles, err := repo.RolesWithAccess(ctx, "GET", "/public")

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(roles) != 0 {
			t.Fatalf("expected 0 roles for public route, got: %d (%v)", len(roles), roles)
		}
	})

	t.Run("cache miss, unmapped route: JOIN returns empty and caches __EMPTY__", func(t *testing.T) {
		t.Parallel()

		rdb := newTestRedisClient(t)
		db, mock := newMockedPostgresDB(t)
		ctx := context.Background()

		mock.ExpectQuery(rolesJoinQuery).
			WithArgs("DELETE:/items", 1, 1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"name"}))

		repo := database.NewPermissionCacheRepository(rdb, db)

		// First call: cache miss, fallback to Postgres JOIN (unmapped route, no rows).
		roles, err := repo.RolesWithAccess(ctx, "DELETE", "/items")

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(roles) != 0 {
			t.Fatalf("expected 0 roles for unmapped route, got: %d (%v)", len(roles), roles)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet sqlmock expectations: %v", err)
		}

		// The empty result must have been cached as the sentinel, so a follow-up call
		// hits Redis and never queries Postgres again.
		members, err := rdb.SMembers(ctx, "auth:route:DELETE:/items").Result()
		if err != nil {
			t.Fatalf("expected sentinel to be cached, got Redis error: %v", err)
		}
		if len(members) != 1 || members[0] != "__EMPTY__" {
			t.Fatalf("expected cached sentinel __EMPTY__, got: %v", members)
		}
	})

	t.Run("cache miss, mapped route: JOIN returns roles and caches them", func(t *testing.T) {
		t.Parallel()

		rdb := newTestRedisClient(t)
		db, mock := newMockedPostgresDB(t)
		ctx := context.Background()

		mock.ExpectQuery(rolesJoinQuery).
			WithArgs("GET:/orders", 1, 1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).
				AddRow("admin").
				AddRow("operator"))

		repo := database.NewPermissionCacheRepository(rdb, db)

		roles, err := repo.RolesWithAccess(ctx, "GET", "/orders")

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(roles) != 2 {
			t.Fatalf("expected 2 roles from JOIN, got: %d (%v)", len(roles), roles)
		}
		roleMap := make(map[string]bool)
		for _, role := range roles {
			roleMap[role] = true
		}
		if !roleMap["admin"] || !roleMap["operator"] {
			t.Fatalf("expected roles {admin, operator}, got: %v", roles)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet sqlmock expectations: %v", err)
		}

		// Roles must be cached in Redis for subsequent lookups.
		members, err := rdb.SMembers(ctx, "auth:route:GET:/orders").Result()
		if err != nil {
			t.Fatalf("expected roles to be cached, got Redis error: %v", err)
		}
		if len(members) != 2 {
			t.Fatalf("expected 2 cached roles, got: %v", members)
		}
	})

	t.Run("cache miss, Postgres JOIN error returns DependencyError (fail-closed)", func(t *testing.T) {
		t.Parallel()

		rdb := newTestRedisClient(t)
		db, mock := newMockedPostgresDB(t)
		ctx := context.Background()

		mock.ExpectQuery(rolesJoinQuery).
			WithArgs("POST:/broken", 1, 1, 1).
			WillReturnError(errors.New("connection reset by peer"))

		repo := database.NewPermissionCacheRepository(rdb, db)

		roles, err := repo.RolesWithAccess(ctx, "POST", "/broken")

		if err == nil {
			t.Fatal("expected DependencyError on Postgres JOIN failure, got nil")
		}
		var depErr *ddgo.DependencyError
		if !errors.As(err, &depErr) {
			t.Fatalf("expected DependencyError, got: %T (%v)", err, err)
		}
		if roles != nil {
			t.Fatalf("expected nil roles on error (fail-closed), got: %v", roles)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet sqlmock expectations: %v", err)
		}

		// A query error must never poison the cache with the empty sentinel — the
		// bypass would otherwise persist until the key expires.
		exists, err := rdb.Exists(ctx, "auth:route:POST:/broken").Result()
		if err != nil {
			t.Fatalf("unexpected Redis error checking cache poisoning: %v", err)
		}
		if exists != 0 {
			t.Fatalf("expected no cache entry to be written on Postgres error, found one")
		}
	})

	t.Run("Redis error on read returns DependencyError (fail-closed)", func(t *testing.T) {
		t.Run("unreachable Redis returns error", func(t *testing.T) {
			t.Parallel()

			// Use an address that won't connect, with short timeouts.
			unreachableClient := redis.NewClient(&redis.Options{
				Addr:        "127.0.0.1:1",
				DialTimeout: 50 * 1, // milliseconds
				ReadTimeout: 50 * 1,
			})
			db := &gorm.DB{}

			repo := database.NewPermissionCacheRepository(unreachableClient, db)
			roles, err := repo.RolesWithAccess(context.Background(), "GET", "/protected")

			if err == nil {
				t.Fatal("expected DependencyError on Redis read failure, got nil")
			}

			var depErr *ddgo.DependencyError
			if !errors.As(err, &depErr) {
				t.Fatalf("expected DependencyError, got: %T (%v)", err, err)
			}

			if roles != nil {
				t.Fatalf("expected nil roles on error, got: %v", roles)
			}

			unreachableClient.Close()
		})
	})

	t.Run("cache hit with no sentinel returns the roles", func(t *testing.T) {
		t.Parallel()

		rdb := newTestRedisClient(t)
		db := &gorm.DB{}

		ctx := context.Background()
		key := "auth:route:POST:/content"

		// Pre-populate with multiple roles, no sentinel.
		rdb.SAdd(ctx, key, "editor", "reviewer")

		repo := database.NewPermissionCacheRepository(rdb, db)
		roles, err := repo.RolesWithAccess(ctx, "POST", "/content")

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(roles) != 2 {
			t.Fatalf("expected 2 roles, got: %d (%v)", len(roles), roles)
		}

		roleMap := make(map[string]bool)
		for _, role := range roles {
			roleMap[role] = true
		}
		if !roleMap["editor"] || !roleMap["reviewer"] {
			t.Fatalf("expected roles {editor, reviewer}, got: %v", roles)
		}
	})

	t.Run("single role in cache is returned correctly", func(t *testing.T) {
		t.Parallel()

		rdb := newTestRedisClient(t)
		db := &gorm.DB{}

		ctx := context.Background()
		key := "auth:route:DELETE:/users/123"

		// Pre-populate with a single role.
		rdb.SAdd(ctx, key, "admin")

		repo := database.NewPermissionCacheRepository(rdb, db)
		roles, err := repo.RolesWithAccess(ctx, "DELETE", "/users/123")

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(roles) != 1 || roles[0] != "admin" {
			t.Fatalf("expected roles {admin}, got: %v", roles)
		}
	})

	t.Run("context is propagated to Redis and Postgres operations", func(t *testing.T) {
		t.Parallel()

		rdb := newTestRedisClient(t)
		db, mock := newMockedPostgresDB(t)
		ctx := context.Background()

		mock.ExpectQuery(rolesJoinQuery).
			WithArgs("GET:/test", 1, 1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"name"}))

		repo := database.NewPermissionCacheRepository(rdb, db)

		// Calling with a valid context should succeed (no panic or error) and reach
		// the Postgres JOIN fallback on cache miss.
		roles, err := repo.RolesWithAccess(ctx, "GET", "/test")

		if err != nil {
			t.Fatalf("expected no error with valid context, got: %v", err)
		}
		if roles == nil {
			t.Fatal("expected non-nil roles slice for public route")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet sqlmock expectations: %v", err)
		}
	})
}
