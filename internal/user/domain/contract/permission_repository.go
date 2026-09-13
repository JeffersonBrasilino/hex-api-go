// Permission repository contract.
//
// Intent: define RBAC permission checking by route, abstracted from storage
// details (cache or persistent database).
// Objective: allow permission checks to query roles with access to a given
// route without depending on infrastructure specifics.
package contract

import "context"

// PermissionRepository defines read access to roles with permission for a
// given HTTP route (method + path).
type PermissionRepository interface {
	// RolesWithAccess returns a list of roles (groups) that have access to
	// the given HTTP method and path. Returns an empty slice when the route
	// has no permission mappings — an empty result signals a public route.
	// A non-empty slice indicates the route is protected; only those roles
	// may access it. On error, the permission check fails closed (deny all).
	RolesWithAccess(ctx context.Context, method, path string) ([]string, error)
}
