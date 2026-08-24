// Token generator contract.
//
// Intent: define JWT issuance and parsing for access and refresh tokens.
// Objective: allow login to issue access/refresh tokens carrying a session
// id, and logout to decode a refresh token back into its session id.
package contract

// TokenGenerator defines JWT access/refresh token generation and refresh
// token parsing for the login and logout user actions.
type TokenGenerator interface {
	// GenerateAccessToken issues an access JWT for the given user and groups.
	GenerateAccessToken(userId string, groups []string) (string, error)
	// GenerateRefreshToken issues a refresh JWT carrying the given session id.
	GenerateRefreshToken(sessionId string) (string, error)
	// ParseRefreshToken decodes a refresh JWT and extracts its session id.
	ParseRefreshToken(token string) (sessionId string, err error)
}
