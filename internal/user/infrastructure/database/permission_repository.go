// Permission repository.
//
// Intent: implement the domain's contract.PermissionRepository, managing two
// backing stores — Redis as the primary, low-latency lookup and Postgres as
// the source of truth queried on cache miss — for RBAC permission checking.
// Objective: efficiently check which roles have access to a given HTTP route
// (method + path), caching non-empty role sets in Redis, and propagating
// infrastructure errors to enforce fail-closed access control (deny all on
// double failure). An empty result (no permission mapping) is never cached —
// it is treated by the caller as a misconfigured protected route, a rare,
// transient condition that should be fixed quickly rather than survive a
// cache TTL.
package database

import (
	"context"
	"fmt"
	"slices"

	"github.com/jeffersonbrasilino/ddgo"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// permissionRouteKeyPrefix namespaces permission role sets in Redis.
const permissionRouteKeyPrefix = "userGroups:permissions:"

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
// means no permission mapping exists for this route — the caller treats that
// as a misconfigured protected route and denies access (this repository is
// only ever queried for routes the authorization middleware is registered
// on; genuinely public routes never reach it). On error, returns nil and a
// ddgo.DependencyError or ddgo.InternalError (fail-closed — deny all).
//
// Behavior:
//   - Cache hit (SMEMBERS returns at least one member): return the member set
//     as-is.
//   - Cache miss (SMEMBERS returns no members — note that Redis' SMEMBERS
//     replies with an empty set, not a nil/not-found error, for a key that
//     does not exist, so the miss is detected by an empty result rather than
//     by err == redis.Nil):
//   - Query Postgres for roles with access to this route via
//     rolesWithAccessFromPostgres.
//   - On Postgres success: cache the result in Redis only if non-empty (an
//     empty result is a misconfiguration expected to be fixed promptly, so it
//     is never cached — there is no sentinel to mark "verified empty"), then
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
		// Redis error other than key not found — fail closed. The raw driver error is
		// carried up as-is; translating it into a user-facing message is the HTTP
		// layer's job, not this adapter's.
		return nil, ddgo.NewDependencyError(
			fmt.Sprintf("Error querying permission cache: %s", err.Error()),
		)
	}

	if len(members) > 0 {
		// Cache hit.
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

	// An empty result means this protected route has no permission mapping — a
	// misconfiguration the caller denies access for. Don't cache it: there's no
	// sentinel to distinguish "verified empty" from "not yet checked" in a Redis set,
	// and caching would let a stale deny outlive the fix once the mapping is added.
	if len(roles) == 0 {
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
// database as a fallback when Redis cache is empty, joining the mapped
// UserGroupsPermissions entity to its ApiRoute (to match the route) and
// UserGroup (to resolve the role name) belongs-to associations. Both are
// to-one relations from UserGroupsPermissions, so Joins is safe here (unlike
// a has-many association, it maps one row to one nested struct, with no risk
// of row multiplication — see FindByUsernameOrDocument for the has-many case,
// which must use Preload instead).
//
// Parameters:
//   - ctx: request-scoped context propagated to GORM.
//   - method: the HTTP method.
//   - path: the resource path.
//
// Returns: a list of distinct role (UsersGroups.Uuid) values with access to
// the route. An empty, non-nil slice means the route has no permission
// mappings (genuinely public). An error is returned only when the query
// itself fails (e.g. connection/driver error) — it is never used to signal
// "no roles found", which the caller (RolesWithAccess) would otherwise treat
// as a fail-open bypass.
//
// Route/method columns: ApiRoutes.Route holds only the path (e.g.
// "/users/123") and UserGroupsPermissions.Action holds the HTTP verb (e.g.
// "GET") — they live on different tables, so method and path are matched as
// two separate WHERE clauses rather than a single concatenated value.
//
// Status filtering: only Status == activeStatus (1) rows are considered
// active on ApiRoutes, UserGroupsPermissions and UsersGroups, matching the
// `default:1` soft-enable convention in gorm_model.go.
//
// Deduplication: since each row is one UserGroupsPermissions (one per
// method+path+role), the same role can appear more than once if it has
// multiple matching permission rows; distinctness is enforced in Go rather
// than via SQL DISTINCT, since Distinct() operates on the base entity's
// columns, not on a joined association's.
func (r *PermissionRepository) rolesWithAccessFromPostgres(
	ctx context.Context,
	method, 
	path string,
) ([]string, error) {
	currentTable := clause.Table{Name: clause.CurrentTable}
	permissions, err := gorm.G[UserGroupsPermissions](r.db).
		Joins(clause.InnerJoin.Association("ApiRoute"), func(
			db gorm.JoinBuilder,
			joinTable clause.Table,
			curTable clause.Table,
		) error {
			db.Where("?.status = ?", joinTable, 1).
				Where("?.route = ?", joinTable, path).
				Select("ID")
			return nil
		}).
		Joins(clause.InnerJoin.Association("UserGroup"), func(
			db gorm.JoinBuilder,
			joinTable clause.Table,
			curTable clause.Table,
		) error {
			db.Where("?.status = ?", joinTable, 1).
				Select("Uuid")
			return nil
		}).
		Where("?.action = ?", currentTable, method).
		Where("?.status = ?", currentTable, 1).
		Find(ctx)
	if err != nil {
		return nil, err
	}

	roles := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		uuid := permission.UserGroup.Uuid
		if slices.Contains(roles, uuid) || uuid == "" {
			continue
		}
		roles = append(roles, uuid)
	}

	return roles, nil
}
