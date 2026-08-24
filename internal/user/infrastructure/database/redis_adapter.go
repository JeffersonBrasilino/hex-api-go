// Redis adapter.
//
// Intent: implement the domain's contract.LoginAttemptStore, contract.LoginSessionStore, and
// contract.RevokeSessionStore using a single Redis client, grouped in one file per the
// {technology}_adapter.go convention for infrastructure concerns sharing the same technology.
// Objective: control failed login attempts and lockout state, persist login sessions, and revoke
// sessions (shared by logout and administrative session revocation) — all backed by Redis.
package database

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jeffersonbrasilino/ddgo"
	"github.com/redis/go-redis/v9"
)

// loginAttemptTtl is the lockout window applied after the first failed login attempt.
const loginAttemptTtl = 5 * time.Minute

// maxLoginAttempts is the number of failed attempts allowed before a username is locked out.
const maxLoginAttempts = 3

// loginAttemptKeyPrefix namespaces failed-login-attempt counters in Redis.
const loginAttemptKeyPrefix = "login_attempt:"

// sessionKeyPrefix namespaces login session records in Redis.
const sessionKeyPrefix = "session:"

// redisRefreshTokenTtl is the lifetime of a refresh token before it must be renewed.
const redisRefreshTokenTtl = 15 * 24 * time.Hour

// RedisAdapter implements contract.LoginAttemptStore, contract.LoginSessionStore, and
// contract.RevokeSessionStore using Redis.
type RedisAdapter struct {
	client *redis.Client
}

// NewRedisAdapter creates a RedisAdapter.
//
// Intent: construct a ready-to-use Redis-based adapter for login attempts and session storage.
// Parameters:
//   - client: the Redis client used to execute commands.
//
// Returns: a *RedisAdapter satisfying contract.LoginAttemptStore, contract.LoginSessionStore, and
// contract.RevokeSessionStore.
func NewRedisAdapter(client *redis.Client) *RedisAdapter {
	return &RedisAdapter{client: client}
}

// loginAttemptKey builds the Redis key holding the failed-attempt counter for username.
func loginAttemptKey(username string) string {
	return loginAttemptKeyPrefix + username
}

// sessionKey builds the Redis key holding the session record for sessionId.
func sessionKey(sessionId string) string {
	return sessionKeyPrefix + sessionId
}

// RegisterFailure records a failed login attempt for the given username.
//
// Intent: increment the username's failed-attempt counter, starting the 5-minute lockout window
// on the first failure.
// Parameters:
//   - ctx: request-scoped context propagated to the Redis client.
//   - username: the username that failed to authenticate.
//
// Returns: nil on success, or a ddgo.InternalError if the Redis operation fails.
func (a *RedisAdapter) RegisterFailure(ctx context.Context, username string) error {
	key := loginAttemptKey(username)
	count, err := a.client.Incr(ctx, key).Result()
	if err != nil {
		return ddgo.NewInternalError(fmt.Sprintf("Error registering login failure: %s", err.Error()))
	}

	if count == 1 {
		if err := a.client.Expire(ctx, key, loginAttemptTtl).Err(); err != nil {
			return ddgo.NewInternalError(fmt.Sprintf("Error setting login attempt lockout window: %s", err.Error()))
		}
	}

	return nil
}

// IsBlocked reports whether the given username is currently locked out.
//
// Intent: check whether the failed-attempt counter for username has reached maxLoginAttempts.
// Parameters:
//   - ctx: request-scoped context propagated to the Redis client.
//   - username: the username to check.
//
// Returns: true if the username is locked out, false otherwise; a ddgo.InternalError if the
// Redis operation fails unexpectedly.
func (a *RedisAdapter) IsBlocked(ctx context.Context, username string) (bool, error) {
	value, err := a.client.Get(ctx, loginAttemptKey(username)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, ddgo.NewInternalError(fmt.Sprintf("Error checking login lockout state: %s", err.Error()))
	}

	count, err := strconv.Atoi(value)
	if err != nil {
		return false, ddgo.NewInternalError(fmt.Sprintf("Error parsing login attempt counter: %s", err.Error()))
	}

	return count >= maxLoginAttempts, nil
}

// Reset clears the failed attempt count for the given username.
//
// Intent: remove the failed-attempt counter for username, typically after a successful login.
// Parameters:
//   - ctx: request-scoped context propagated to the Redis client.
//   - username: the username whose counter should be cleared.
//
// Returns: nil on success, or a ddgo.InternalError if the Redis operation fails.
func (a *RedisAdapter) Reset(ctx context.Context, username string) error {
	if err := a.client.Del(ctx, loginAttemptKey(username)).Err(); err != nil {
		return ddgo.NewInternalError(fmt.Sprintf("Error resetting login attempt counter: %s", err.Error()))
	}
	return nil
}

// Save persists a session identified by sessionId for userId, expiring after ttl.
//
// Intent: store the session record in Redis so it can later be validated or revoked.
// Parameters:
//   - ctx: request-scoped context propagated to the Redis client.
//   - sessionId: the id of the session being created.
//   - userId: the id of the user the session belongs to.
//   - ttl: the duration after which the session record expires.
//
// Returns: nil on success, or a ddgo.InternalError if the Redis operation fails.
func (a *RedisAdapter) Save(ctx context.Context, sessionId, refreshToken string) error {
	if err := a.client.Set(ctx, sessionKey(sessionId), refreshToken, redisRefreshTokenTtl).Err(); err != nil {
		return ddgo.NewInternalError(fmt.Sprintf("Error saving login session: %s", err.Error()))
	}
	return nil
}

// Delete removes the session identified by sessionId. It returns a "not found" error when the
// session does not exist.
//
// Intent: end a session so it can no longer be used to access the system, shared by logout and
// administrative session revocation.
// Parameters:
//   - ctx: request-scoped context propagated to the Redis client.
//   - sessionId: the id of the session to remove.
//
// Returns: nil on success, a ddgo.NotFoundError if no session with sessionId exists, or a
// ddgo.InternalError if the Redis operation fails unexpectedly.
func (a *RedisAdapter) Delete(ctx context.Context, sessionId string) error {
	deleted, err := a.client.Del(ctx, sessionKey(sessionId)).Result()
	if err != nil {
		return ddgo.NewInternalError(fmt.Sprintf("Error deleting session: %s", err.Error()))
	}
	if deleted == 0 {
		return ddgo.NewNotFoundError("Sessão não encontrada")
	}
	return nil
}
