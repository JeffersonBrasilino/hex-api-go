// Package checkpermission defines the query for checking user permissions.
// It orchestrates the evaluation of access control rules based on an access token,
// HTTP method, and resource path, determining whether the user has permission to
// access the requested resource.
package checkpermission

// Query represents a request to check whether a user has permission to access a
// specific resource. It holds the access token, HTTP method, and resource path
// required to make the access control decision.
type Query struct {
	// AccessToken is the JWT or bearer token of the authenticated user.
	AccessToken string
	// Method is the HTTP method of the request (e.g. "GET", "POST", "DELETE").
	Method string
	// Path is the resource path being accessed (e.g. "/users/123", "/settings").
	Path string
}

// Name returns the query name for CQRS bus routing and identification.
func (q *Query) Name() string {
	return "checkPermission"
}
