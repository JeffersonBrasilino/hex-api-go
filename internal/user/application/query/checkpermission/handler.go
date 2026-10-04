// Package checkpermission defines the query handler for checking user permissions.
//
// Intent: orchestrate permission checking by parsing the access token, querying roles with
// access, and deciding whether to allow or deny.
//
// Objective: enforce RBAC access control by checking if the authenticated user's groups have
// permission to access the requested resource. Genuinely public routes never reach this
// handler at all — the authorization middleware that dispatches this query is only ever
// registered on protected routes/groups — so a route with no permission mapping here is a
// misconfiguration, not a public route, and is denied.
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
// Returns: true (as an any value) if access is allowed (user has permission to the
// resource), or an error:
//   - domain.InvalidSessionError when the access token is missing, malformed, or expired
//   - domain.AccessDeniedError when the user lacks the required permissions, or when the
//     route has no permission mapping at all (misconfiguration — fail closed)
//   - any error surfaced by the underlying contracts (permission repository failure)
//
// Business rule (RN-01 revised):
//   - Every route reaching this handler is protected — the authorization middleware is
//     only ever registered on routes/groups that require access control — so the access
//     token is parsed first, unconditionally.
//   - Empty roles list (no permission mapping for this route) = deny via
//     AccessDeniedError. This is a misconfigured protected route, not a public one:
//     genuinely public routes have no middleware attached and never dispatch this query.
//   - Non-empty roles list with no intersection with user's groups = deny via
//     AccessDeniedError.
//   - Non-empty roles list with intersection with user's groups = allow access.
func (h *Handler) Handle(ctx context.Context, q *Query) (any, error) {
	// Every caller of this handler is hitting a protected route, so a token is always
	// required.
	_, userGroups, errParse := h.tokenValidator.ParseAccessToken(q.AccessToken)
	if errParse != nil {
		return nil, domain.NewInvalidSessionError(
			"Invalid or expired access token",
		)
	}

	rolesWithAccess, errRoles := h.permissionRepository.RolesWithAccess(
		ctx, q.Method, q.Path,
	)
	if errRoles != nil {
		// Infrastructure error on permission repository — propagate as-is (fail closed at
		// HTTP layer).
		return nil, errRoles
	}

	// Business rule: empty roles list = no permission mapping configured for this
	// protected route — fail closed rather than treat it as public.
	if len(rolesWithAccess) == 0 {
		return nil, domain.NewAccessDeniedError(
			"No permission mapping configured for this resource",
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
