// Permission repository.
//
// Intent: implement the domain's contract.PermissionRepository, managing two
// backing stores — Redis as the primary, low-latency lookup and Postgres as
// the source of truth queried on cache miss — for RBAC permission checking.
// Objective: efficiently check which roles have access to a given HTTP route
// (method + path), storing results in Redis with a sentinel value (__EMPTY__)
// for empty permission sets, and propagating infrastructure errors to enforce
// fail-closed access control (deny all on double failure).
package database

import (
	"context"
	"fmt"

	"github.com/jeffersonbrasilino/ddgo"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// permissionRouteKeyPrefix namespaces permission role sets in Redis.
const permissionRouteKeyPrefix = "auth:route:"

// emptyPermissionSentinel is a sentinel value stored in Redis when a route has
// no permission mappings (public route).
const emptyPermissionSentinel = "__EMPTY__"

// activeStatus is the only Status value considered active/enabled, matching the
// `default:1` convention used across gorm_model.go's soft-enable columns
// (ApiRoutes.Status, UserGroupsPermissions.Status, UsersGroups.Status).
const activeStatus = 1

// PermissionRepository implements contract.PermissionRepository using
// Redis cache-aside with Postgres as the fallback store.
type PermissionRepository struct {
	redis *redis.Client
	db    *gorm.DB
}

// NewPermissionRepository creates a PermissionRepository.
//
// Intent: construct a ready-to-use permission repository combining Redis cache
// and Postgres persistence for RBAC access control.
// Parameters:
//   - rdb: the Redis client used for caching role sets by route.
//   - db: the GORM database connection used as fallback for permission lookups.
//
// Returns: a *PermissionRepository satisfying contract.PermissionRepository.
func NewPermissionRepository(
	rdb *redis.Client,
	db *gorm.DB,
) *PermissionRepository {
	return &PermissionRepository{
		redis: rdb,
		db:    db,
	}
}

// permissionRouteKey builds the Redis key for a given HTTP method and path.
func permissionRouteKey(method, path string) string {
	return permissionRouteKeyPrefix + method + ":" + path
}

// RolesWithAccess returns a list of roles (groups) that have access to the
// given HTTP method and path, using a cache-aside pattern backed by Redis.
//
// Intent: efficiently check permission mappings for a route, caching the result
// in Redis to avoid repeated Postgres queries, while enforcing fail-closed
// semantics (deny all on any infrastructure failure).
//
// Parameters:
//   - ctx: request-scoped context propagated to Redis and GORM clients.
//   - method: the HTTP method (e.g. "GET", "POST", "DELETE").
//   - path: the resource path being accessed.
//
// Returns: a list of role names that have access to the route. An empty list
// signals a public route (no permission mappings exist). On error, returns nil
// and a ddgo.DependencyError or ddgo.InternalError (fail-closed — deny all).
//
// Behavior:
//   - Cache hit (SMEMBERS returns at least one member):
//   - If the sentinel __EMPTY__ is present, return an empty list (public).
//   - Otherwise, return the member set as-is.
//   - Cache miss (SMEMBERS returns no members — note that Redis' SMEMBERS
//     replies with an empty set, not a nil/not-found error, for a key that
//     does not exist, so the miss is detected by an empty result rather than
//     by err == redis.Nil):
//   - Query Postgres for roles with access to this route via
//     rolesWithAccessFromPostgres.
//   - On Postgres success: store the result in Redis (or __EMPTY__ if empty),
//     return it.
//   - On Postgres error: return a DependencyError (fail closed).
//   - Redis error on read: return a DependencyError (fail closed).
//   - Double failure (Redis error AND Postgres error): return a DependencyError.
func (r *PermissionRepository) RolesWithAccess(
	ctx context.Context,
	method, path string,
) ([]string, error) {
	key := permissionRouteKey(method, path)

	// Try to fetch from Redis cache. SMEMBERS never errors with redis.Nil for a
	// missing key — it replies with an empty set — so a cache miss is detected
	// by an empty, error-free result rather than by comparing err to redis.Nil.
	members, err := r.redis.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		// Redis error other than key not found — fail closed.
		return nil, ddgo.NewDependencyError(
			fmt.Sprintf("Error querying permission cache: %s", err.Error()),
		)
	}

	if len(members) > 0 {
		// Cache hit. Check for sentinel value.
		if len(members) == 1 && members[0] == emptyPermissionSentinel {
			return []string{}, nil
		}
		return members, nil
	}

	// Cache miss. Query Postgres for the roles with access to this route.
	roles, postgresErr := r.rolesWithAccessFromPostgres(ctx, method, path)
	if postgresErr != nil {
		// Postgres error — fail closed.
		return nil, ddgo.NewDependencyError(
			fmt.Sprintf("Error querying permissions from database: %s",
				postgresErr.Error()),
		)
	}

	// Cache the result. Store the sentinel if the result is empty.
	if len(roles) == 0 {
		if err := r.redis.SAdd(ctx, key, emptyPermissionSentinel).Err(); err != nil {
			// Redis write error — return permission result but log the cache miss.
			// Note: We return the Postgres result (which may be empty) rather than
			// fail closed, since we successfully retrieved from the primary source.
			return roles, nil
		}
		return []string{}, nil
	}

	// Store the roles in Redis as a set.
	if err := r.redis.SAdd(ctx, key, roles).Err(); err != nil {
		// Redis write error — return permission result anyway, since we have a
		// valid answer from Postgres.
		return roles, nil
	}

	return roles, nil
}

// rolesWithAccessFromPostgres queries Postgres for roles with access to the
// given route (method + path).
//
// Intent: retrieve the permission mapping for a route from the persistent
// database as a fallback when Redis cache is empty, joining
// UserGroupsPermissions -> ApiRoutes (to match the route) and
// UserGroupsPermissions -> UsersGroups (to resolve the role name).
//
// Parameters:
//   - ctx: request-scoped context propagated to GORM.
//   - method: the HTTP method.
//   - path: the resource path.
//
// Returns: a list of distinct role (UsersGroups.Name) values with access to
// the route. An empty, non-nil slice means the route has no permission
// mappings (genuinely public). An error is returned only when the query
// itself fails (e.g. connection/driver error) — it is never used to signal
// "no roles found", which the caller (RolesWithAccess) would otherwise treat
// as a fail-open bypass.
//
// Route format assumption: ApiRoutes.Route is a single string column with no
// separate HTTP method column. No other usage, seed data, or fixture in the
// repository documents its stored format, so this method adopts the same
// convention already used by permissionRouteKey for the Redis key — Route is
// stored as "{METHOD}:{PATH}" (colon-separated), e.g. "GET:/users/123". This
// is an explicit assumption, not verified against a real seeded row.
//
// Status filtering: only Status == activeStatus (1) rows are considered
// active on ApiRoutes, UserGroupsPermissions and UsersGroups, matching the
// `default:1` soft-enable convention in gorm_model.go.
func (r *PermissionRepository) rolesWithAccessFromPostgres(
	ctx context.Context,
	method, path string,
) ([]string, error) {
	route := method + ":" + path

	roles := []string{}
	err := r.db.WithContext(ctx).
		Table(`"hex-api-go"."user_groups_permissions"`).
		Joins(
			`JOIN "hex-api-go"."users_groups" ON "hex-api-go"."users_groups"."id" = `+
				`"hex-api-go"."user_groups_permissions"."user_group_id"`,
		).
		Joins(
			`JOIN "hex-api-go"."api_routes" ON "hex-api-go"."api_routes"."id" = `+
				`"hex-api-go"."user_groups_permissions"."api_route_id"`,
		).
		Where(`"hex-api-go"."api_routes"."route" = ?`, route).
		Where(`"hex-api-go"."api_routes"."status" = ?`, activeStatus).
		Where(`"hex-api-go"."user_groups_permissions"."status" = ?`, activeStatus).
		Where(`"hex-api-go"."users_groups"."status" = ?`, activeStatus).
		Distinct(`"hex-api-go"."users_groups"."name"`).
		Pluck(`"hex-api-go"."users_groups"."name"`, &roles).Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}
