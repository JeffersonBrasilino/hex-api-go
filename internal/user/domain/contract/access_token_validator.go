// Access token validator contract.
//
// Intent: define JWT access token decoding for permission checking.
// Objective: allow permission checks to extract userId and groups from an
// access token JWT, independent from the login/logout actions which use
// TokenGenerator for refresh token parsing (Interface Segregation).
package contract

// AccessTokenValidator defines JWT access token parsing for the permission
// checking user action.
type AccessTokenValidator interface {
	// ParseAccessToken decodes an access JWT and extracts its user ID and groups.
	ParseAccessToken(token string) (userId string, groups []string, err error)
}
