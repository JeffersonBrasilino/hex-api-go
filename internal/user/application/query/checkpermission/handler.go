// Package checkpermission defines the query handler for checking user permissions.
//
// Intent: orchestrate permission checking by parsing the access token, querying roles with
// access, and deciding whether to allow, deny, or allow as a public route.
//
// Objective: enforce RBAC access control by checking if the authenticated user's groups have
// permission to access the requested resource, while treating routes without any permission
// mappings as public.
package checkpermission

import (
	"context"

	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain/contract"
)

// Handler orchestrates permission checking: it parses the access token, queries the roles with
// access to the requested resource, and decides whether to allow, deny, or treat the route as
// public (no permission mappings exist).
type Handler struct {
	tokenValidator       contract.AccessTokenValidator
	permissionRepository contract.PermissionRepository
}

// NewQueryHandler creates a Handler for the checkpermission query.
//
// Intent: wire the handler with the domain contracts it needs to execute the permission
// check.
//
// Parameters:
//   - tokenValidator: contract used to parse the access token and extract user ID and
//     groups.
//   - permissionRepository: contract used to query roles with access to a given HTTP
//     method and path.
//
// Returns: a ready-to-use *Handler.
func NewQueryHandler(
	tokenValidator contract.AccessTokenValidator,
	permissionRepository contract.PermissionRepository,
) *Handler {
	return &Handler{
		tokenValidator:       tokenValidator,
		permissionRepository: permissionRepository,
	}
}

// Handle executes the permission check use case.
//
// Intent: verify that the authenticated user's groups have permission to access the
// requested resource, based on the access token, HTTP method, and resource path.
//
// Parameters:
//   - ctx: request-scoped context propagated to the contracts.
//   - q: the checkpermission Query with the access token, HTTP method, and resource
//     path.
//
// Returns: true (as an any value) if access is allowed (user has permission or route is
// public), or an error:
//   - domain.InvalidSessionError when the route is protected and the access token is
//     missing, malformed, or expired
//   - domain.AccessDeniedError when the user lacks the required permissions
//   - any error surfaced by the underlying contracts (permission repository failure)
//
// Business rule (RN-01 revised):
//   - Empty roles list (no permission mappings) = public route, allow access — evaluated
//     BEFORE the access token is parsed, so a request with no token at all (e.g. an
//     unauthenticated call to a genuinely public route such as login) is never rejected
//     for lacking a token it was never meant to carry.
//   - Non-empty roles list with no intersection with user's groups = deny via
//     AccessDeniedError.
//   - Non-empty roles list with intersection with user's groups = allow access.
//
// This ordering is deliberate: RF-04's "no static public route list" requirement means the
// only way to know a route is public is to ask the permission repository first — checking
// the token before that would force every caller, including unauthenticated ones hitting a
// public route, to present a token it doesn't need.
func (h *Handler) Handle(ctx context.Context, q *Query) (any, error) {
	// Query the roles with access to the requested resource (method + path) first — this is
	// the only way to know whether the route is public, and must not depend on a token being
	// present.
	rolesWithAccess, errRoles := h.permissionRepository.RolesWithAccess(
		ctx, q.Method, q.Path,
	)
	if errRoles != nil {
		// Infrastructure error on permission repository — propagate as-is (fail closed at
		// HTTP layer).
		return nil, errRoles
	}

	// Business rule: empty roles list = public route, allow access without requiring a token.
	if len(rolesWithAccess) == 0 {
		return true, nil
	}

	// The route is protected — a token is now required. Parse it to extract user groups.
	_, userGroups, errParse := h.tokenValidator.ParseAccessToken(q.AccessToken)
	if errParse != nil {
		return nil, domain.NewInvalidSessionError(
			"Invalid or expired access token",
		)
	}

	// Check if any of the user's groups intersects with the roles that have access.
	userGroupSet := make(map[string]bool)
	for _, group := range userGroups {
		userGroupSet[group] = true
	}

	for _, role := range rolesWithAccess {
		if userGroupSet[role] {
			// User has at least one group with access to this resource.
			return true, nil
		}
	}

	// User lacks the required permissions — deny access.
	return nil, domain.NewAccessDeniedError(
		"User does not have permission to access this resource",
	)
}
