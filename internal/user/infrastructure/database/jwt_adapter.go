// JWT token generator adapter.
//
// Intent: implement the domain's contract.TokenGenerator using HS256-signed JWTs from the
// github.com/golang-jwt/jwt/v5 library.
// Objective: issue access tokens carrying the user id and groups, issue refresh tokens carrying
// a session id, and parse a refresh token back into its session id, without exposing JWT-specific
// details to the domain or application layers.
package database

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// accessTokenTtl is the lifetime of an access token before it must be renewed.
const accessTokenTtl = 2 * time.Hour

// refreshTokenTtl is the lifetime of a refresh token before it must be renewed.
const refreshTokenTtl = 15 * 24 * time.Hour

// accessTokenClaims are the JWT claims carried by an access token.
type accessTokenClaims struct {
	Groups []string `json:"groups"`
	jwt.RegisteredClaims
}

// refreshTokenClaims are the JWT claims carried by a refresh token.
type refreshTokenClaims struct {
	jwt.RegisteredClaims
}

// JwtAdapter implements contract.TokenGenerator using HS256-signed JWTs.
type JwtAdapter struct{}

// NewJwtAdapter creates a JwtAdapter.
//
// Intent: construct a ready-to-use JWT-based token generator.
// Returns: a *JwtAdapter satisfying contract.TokenGenerator.
func NewJwtAdapter() *JwtAdapter {
	return &JwtAdapter{}
}

// secret returns the signing key read from the JWT_SECRET environment variable.
func (a *JwtAdapter) secret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

// GenerateAccessToken issues an access JWT for the given user and groups.
//
// Intent: produce a short-lived, HS256-signed token carrying the user id (as `sub`) and group
// names, used to authorize subsequent requests.
// Parameters:
//   - userId: the id of the authenticated user, set as the token's `sub` claim.
//   - groups: the names of the groups the user belongs to.
//
// Returns: the signed JWT string, or an error if signing fails.
func (a *JwtAdapter) GenerateAccessToken(userId string, groups []string) (string, error) {
	claims := &accessTokenClaims{
		Groups: groups,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTtl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secret())
}

// GenerateRefreshToken issues a refresh JWT carrying the given session id.
//
// Intent: produce a long-lived, HS256-signed token carrying the session id (as `jti`), used to
// end a session on logout.
// Parameters:
//   - sessionId: the id of the session to associate with the refresh token.
//
// Returns: the signed JWT string, or an error if signing fails.
func (a *JwtAdapter) GenerateRefreshToken(sessionId string) (string, error) {
	claims := &refreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        sessionId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTtl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secret())
}

// ParseRefreshToken decodes a refresh JWT and extracts its session id.
//
// Intent: validate the token's signature and expiration, then recover the session id it carries,
// so logout can identify which session to end.
// Parameters:
//   - tokenString: the refresh JWT to decode.
//
// Returns: the session id carried by the token, or an error if the token is malformed, has an
// invalid signature, or is expired.
func (a *JwtAdapter) ParseRefreshToken(tokenString string) (sessionId string, err error) {
	claims := &refreshTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return a.secret(), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	return claims.ID, nil
}
