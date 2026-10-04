// Permission cache repository tests.
//
// Intent: exercise the PermissionRepository's cache-aside pattern,
// testing the Redis caching logic with miniredis and the real Postgres JOIN
// fallback (mocked via go-sqlmock, following the convention established in
// gorm_user_repository_test.go).
// Objective: verify cache hit, cache miss triggering the Postgres JOIN
// fallback (mapped route returns roles and caches them; unmapped route
// returns empty and is never cached), Postgres query errors (fail-closed
// DependencyError), Redis errors, and fail-closed semantics (deny all on
// infrastructure failure).
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

// rolesJoinQuery matches the Postgres JOIN emitted by rolesWithAccessFromPostgres via the
// gorm.G[UserGroupsPermissions] generics API: UserGroupsPermissions -> ApiRoute (route match)
// and UserGroupsPermissions -> UserGroup (role uuid), filtered by activeStatus (1) on all three
// tables. Argument order follows the Joins declaration order, then the base Where calls:
// ApiRoute.status, ApiRoute.route, UserGroup.status, action, status.
const rolesJoinQuery = `SELECT .+ ` +
	`FROM "hex-api-go"\."user_groups_permissions" ` +
	`INNER JOIN "hex-api-go"\."api_routes" "ApiRoute" ON .+ ` +
	`INNER JOIN "hex-api-go"\."users_groups" "UserGroup" ON .+ ` +
	`WHERE .+`

// rolesQueryColumns names the columns rolesWithAccessFromPostgres selects: the base entity's
// minimal "id" (unused, kept only because the generics Select requires at least one column),
// ApiRoute's "id" (unused, joined only to filter by route/status) and UserGroup's "uuid" (the
// role identifier returned to the caller).
var rolesQueryColumns = []string{"id", "ApiRoute__id", "UserGroup__uuid"}

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

// TestPermissionRepository_RolesWithAccess tests the cache-aside pattern.
func TestPermissionRepository_RolesWithAccess(t *testing.T) {
	t.Run("cache hit returns stored roles", func(t *testing.T) {
		t.Parallel()

		rdb := newTestRedisClient(t)
		db := &gorm.DB{} // Not used in this test

		ctx := context.Background()
		key := "userGroups:permissions:GET:/users"

		// Pre-populate Redis with roles.
		rdb.SAdd(ctx, key, "admin", "operator")

		repo := database.NewPermissionRepository(rdb, db)
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
		key := "userGroups:permissions:GET:/public"

		// Pre-populate Redis with a mapped role.
		rdb.SAdd(ctx, key, "admin")

		repo := database.NewPermissionRepository(rdb, db)
		roles, err := repo.RolesWithAccess(ctx, "GET", "/public")

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(roles) != 1 || roles[0] != "admin" {
			t.Fatalf("expected roles {admin}, got: %v", roles)
		}
	})

	t.Run("cache miss, unmapped route: JOIN returns empty and is never cached", func(t *testing.T) {
		t.Parallel()

		rdb := newTestRedisClient(t)
		db, mock := newMockedPostgresDB(t)
		ctx := context.Background()

		mock.ExpectQuery(rolesJoinQuery).
			WithArgs(1, "/items", 1, "DELETE", 1).
			WillReturnRows(sqlmock.NewRows(rolesQueryColumns))

		repo := database.NewPermissionRepository(rdb, db)

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

		// The empty result must never be cached — there's no sentinel anymore, and
		// caching it would let a misconfiguration outlive its fix until the key expires.
		exists, err := rdb.Exists(ctx, "userGroups:permissions:DELETE:/items").Result()
		if err != nil {
			t.Fatalf("unexpected Redis error checking cache: %v", err)
		}
		if exists != 0 {
			t.Fatalf("expected no cache entry for an unmapped route, found one")
		}
	})

	t.Run("cache miss, mapped route: JOIN returns roles and caches them", func(t *testing.T) {
		t.Parallel()

		rdb := newTestRedisClient(t)
		db, mock := newMockedPostgresDB(t)
		ctx := context.Background()

		mock.ExpectQuery(rolesJoinQuery).
			WithArgs(1, "/orders", 1, "GET", 1).
			WillReturnRows(sqlmock.NewRows(rolesQueryColumns).
				AddRow(1, 1, "admin").
				AddRow(2, 1, "operator"))

		repo := database.NewPermissionRepository(rdb, db)

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
		members, err := rdb.SMembers(ctx, "userGroups:permissions:GET:/orders").Result()
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
			WithArgs(1, "/broken", 1, "POST", 1).
			WillReturnError(errors.New("connection reset by peer"))

		repo := database.NewPermissionRepository(rdb, db)

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

		// A query error must never write a cache entry — an empty/failed lookup
		// masquerading as "verified" would persist until the key expires.
		exists, err := rdb.Exists(ctx, "userGroups:permissions:POST:/broken").Result()
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

			repo := database.NewPermissionRepository(unreachableClient, db)
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

	t.Run("cache hit returns the roles", func(t *testing.T) {
		t.Parallel()

		rdb := newTestRedisClient(t)
		db := &gorm.DB{}

		ctx := context.Background()
		key := "userGroups:permissions:POST:/content"

		// Pre-populate with multiple roles.
		rdb.SAdd(ctx, key, "editor", "reviewer")

		repo := database.NewPermissionRepository(rdb, db)
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
		key := "userGroups:permissions:DELETE:/users/123"

		// Pre-populate with a single role.
		rdb.SAdd(ctx, key, "admin")

		repo := database.NewPermissionRepository(rdb, db)
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
			WithArgs(1, "/test", 1, "GET", 1).
			WillReturnRows(sqlmock.NewRows(rolesQueryColumns))

		repo := database.NewPermissionRepository(rdb, db)

		// Calling with a valid context should succeed (no panic or error) and reach
		// the Postgres JOIN fallback on cache miss.
		roles, err := repo.RolesWithAccess(ctx, "GET", "/test")

		if err != nil {
			t.Fatalf("expected no error with valid context, got: %v", err)
		}
		if roles == nil {
			t.Fatal("expected non-nil (empty) roles slice for an unmapped route")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet sqlmock expectations: %v", err)
		}
	})
}
