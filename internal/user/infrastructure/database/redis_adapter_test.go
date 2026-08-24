// Redis adapter tests.
//
// Intent: exercise database.RedisAdapter's contract.LoginAttemptStore,
// contract.LoginSessionStore, and contract.RevokeSessionStore behavior.
// Objective: verify RegisterFailure, IsBlocked, and Reset (failed-login lockout), Save (session
// persistence), and Delete (session revocation), including their Redis-error and malformed-data
// paths. Each (sub)test runs against its own isolated in-memory Redis server started via
// miniredis.RunT, so tests never depend on an external Redis instance and never collide on keys
// even when run in parallel.
package database_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/database"
	"github.com/redis/go-redis/v9"
)

// unreachableRedisAddr points to an address nothing listens on, used to force Redis command
// errors deterministically without depending on any external Redis instance.
const unreachableRedisAddr = "127.0.0.1:1"

// newTestRedisClient starts an isolated in-memory Redis server via miniredis, scoped to the
// lifetime of t (miniredis.RunT registers its own t.Cleanup to shut the server down), and returns
// a client pointed at it. Each call gets its own server instance, so parallel (sub)tests never
// share state or collide on keys.
func newTestRedisClient(t *testing.T) *redis.Client {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	t.Cleanup(func() {
		_ = client.Close()
	})
	return client
}

// newUnreachableRedisClient returns a client pointed at an address nothing listens on, with a
// short timeout, so commands fail fast with a real Redis error.
func newUnreachableRedisClient(t *testing.T) *redis.Client {
	t.Helper()
	client := redis.NewClient(&redis.Options{
		Addr:        unreachableRedisAddr,
		DialTimeout: 200 * time.Millisecond,
		ReadTimeout: 200 * time.Millisecond,
	})
	t.Cleanup(func() {
		_ = client.Close()
	})
	return client
}

// uniqueName builds a per-subtest unique identifier so subtests never collide on keys even when
// they happen to share a Redis instance.
func uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func TestNewRedisAdapter(t *testing.T) {
	t.Run("should return a non-nil RedisAdapter", func(t *testing.T) {
		t.Parallel()
		client := newTestRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		if adapter == nil {
			t.Fatal("NewRedisAdapter should return a non-nil adapter")
		}
	})
}

func TestRedisAdapter_RegisterFailure(t *testing.T) {
	t.Run("should set the counter to 1 and start the lockout window on the first failure", func(t *testing.T) {
		t.Parallel()
		client := newTestRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()
		username := uniqueName("user")
		key := "login_attempt:" + username

		err := adapter.RegisterFailure(ctx, username)

		if err != nil {
			t.Fatalf("RegisterFailure should not return an error, got: %v", err)
		}
		value, getErr := client.Get(ctx, key).Result()
		if getErr != nil {
			t.Fatalf("expected key %q to exist, got error: %v", key, getErr)
		}
		if value != "1" {
			t.Fatalf("expected counter to be 1, got: %s", value)
		}
		ttl, ttlErr := client.TTL(ctx, key).Result()
		if ttlErr != nil {
			t.Fatalf("TTL should not return an error, got: %v", ttlErr)
		}
		if ttl <= 0 || ttl > 5*time.Minute {
			t.Fatalf("expected TTL to be set within the 5-minute lockout window, got: %v", ttl)
		}
	})

	t.Run("should increment the counter on subsequent failures without resetting the TTL", func(t *testing.T) {
		t.Parallel()
		client := newTestRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()
		username := uniqueName("user")
		key := "login_attempt:" + username

		if err := adapter.RegisterFailure(ctx, username); err != nil {
			t.Fatalf("first RegisterFailure should not return an error, got: %v", err)
		}
		firstTtl, _ := client.TTL(ctx, key).Result()

		if err := adapter.RegisterFailure(ctx, username); err != nil {
			t.Fatalf("second RegisterFailure should not return an error, got: %v", err)
		}

		value, getErr := client.Get(ctx, key).Result()
		if getErr != nil {
			t.Fatalf("expected key %q to exist, got error: %v", key, getErr)
		}
		if value != "2" {
			t.Fatalf("expected counter to be 2, got: %s", value)
		}
		secondTtl, ttlErr := client.TTL(ctx, key).Result()
		if ttlErr != nil {
			t.Fatalf("TTL should not return an error, got: %v", ttlErr)
		}
		if secondTtl <= 0 {
			t.Fatalf("expected TTL to remain set after a second failure, got: %v", secondTtl)
		}
		if secondTtl > firstTtl {
			t.Fatalf("expected TTL not to be extended on subsequent failures, first: %v second: %v", firstTtl, secondTtl)
		}
	})

	t.Run("should return an internal error when Redis is unreachable", func(t *testing.T) {
		t.Parallel()
		client := newUnreachableRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()

		err := adapter.RegisterFailure(ctx, uniqueName("user"))

		if err == nil {
			t.Fatal("RegisterFailure should return an error when Redis is unreachable")
		}
	})
}

func TestRedisAdapter_IsBlocked(t *testing.T) {
	t.Run("should return false when no attempts have been registered", func(t *testing.T) {
		t.Parallel()
		client := newTestRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()
		username := uniqueName("user")

		blocked, err := adapter.IsBlocked(ctx, username)

		if err != nil {
			t.Fatalf("IsBlocked should not return an error, got: %v", err)
		}
		if blocked {
			t.Fatal("IsBlocked should return false when the username has no failed attempts")
		}
	})

	t.Run("should return false when the counter is below the maximum allowed attempts", func(t *testing.T) {
		t.Parallel()
		client := newTestRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()
		username := uniqueName("user")
		key := "login_attempt:" + username
		if err := client.Set(ctx, key, "2", time.Minute).Err(); err != nil {
			t.Fatalf("failed to seed counter: %v", err)
		}

		blocked, err := adapter.IsBlocked(ctx, username)

		if err != nil {
			t.Fatalf("IsBlocked should not return an error, got: %v", err)
		}
		if blocked {
			t.Fatal("IsBlocked should return false when the counter is below the maximum")
		}
	})

	t.Run("should return true when the counter reached the maximum allowed attempts", func(t *testing.T) {
		t.Parallel()
		client := newTestRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()
		username := uniqueName("user")
		key := "login_attempt:" + username
		if err := client.Set(ctx, key, "3", time.Minute).Err(); err != nil {
			t.Fatalf("failed to seed counter: %v", err)
		}

		blocked, err := adapter.IsBlocked(ctx, username)

		if err != nil {
			t.Fatalf("IsBlocked should not return an error, got: %v", err)
		}
		if !blocked {
			t.Fatal("IsBlocked should return true when the counter reached the maximum")
		}
	})

	t.Run("should return an internal error when the stored counter is not a valid number", func(t *testing.T) {
		t.Parallel()
		mr := miniredis.RunT(t)
		client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
		t.Cleanup(func() {
			_ = client.Close()
		})
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()
		username := uniqueName("user")
		key := "login_attempt:" + username
		if err := mr.Set(key, "not-a-number"); err != nil {
			t.Fatalf("failed to seed counter: %v", err)
		}

		blocked, err := adapter.IsBlocked(ctx, username)

		if err == nil {
			t.Fatal("IsBlocked should return an error when the stored counter cannot be parsed")
		}
		if blocked {
			t.Fatal("IsBlocked should return false alongside the error")
		}
	})

	t.Run("should return an internal error when Redis is unreachable", func(t *testing.T) {
		t.Parallel()
		client := newUnreachableRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()

		blocked, err := adapter.IsBlocked(ctx, uniqueName("user"))

		if err == nil {
			t.Fatal("IsBlocked should return an error when Redis is unreachable")
		}
		if blocked {
			t.Fatal("IsBlocked should return false alongside the error")
		}
	})
}

func TestRedisAdapter_Reset(t *testing.T) {
	t.Run("should clear an existing failed-attempt counter", func(t *testing.T) {
		t.Parallel()
		client := newTestRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()
		username := uniqueName("user")
		key := "login_attempt:" + username
		if err := client.Set(ctx, key, "2", time.Minute).Err(); err != nil {
			t.Fatalf("failed to seed counter: %v", err)
		}

		err := adapter.Reset(ctx, username)

		if err != nil {
			t.Fatalf("Reset should not return an error, got: %v", err)
		}
		exists, existsErr := client.Exists(ctx, key).Result()
		if existsErr != nil {
			t.Fatalf("Exists should not return an error, got: %v", existsErr)
		}
		if exists != 0 {
			t.Fatal("Reset should remove the failed-attempt counter")
		}
	})

	t.Run("should not return an error when there is no counter to clear", func(t *testing.T) {
		t.Parallel()
		client := newTestRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()

		err := adapter.Reset(ctx, uniqueName("user"))

		if err != nil {
			t.Fatalf("Reset should not return an error for a non-existent counter, got: %v", err)
		}
	})

	t.Run("should return an internal error when Redis is unreachable", func(t *testing.T) {
		t.Parallel()
		client := newUnreachableRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()

		err := adapter.Reset(ctx, uniqueName("user"))

		if err == nil {
			t.Fatal("Reset should return an error when Redis is unreachable")
		}
	})
}

func TestRedisAdapter_Save(t *testing.T) {
	t.Run("should persist the session with the given ttl", func(t *testing.T) {
		t.Parallel()
		client := newTestRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()
		sessionId := uniqueName("session")
		userId := uniqueName("user-id")
		key := "session:" + sessionId

		err := adapter.Save(ctx, sessionId, userId)

		if err != nil {
			t.Fatalf("Save should not return an error, got: %v", err)
		}
		value, getErr := client.Get(ctx, key).Result()
		if getErr != nil {
			t.Fatalf("expected key %q to exist, got error: %v", key, getErr)
		}
		if value != userId {
			t.Fatalf("expected stored value to be %q, got: %q", userId, value)
		}
		ttl, ttlErr := client.TTL(ctx, key).Result()
		if ttlErr != nil {
			t.Fatalf("TTL should not return an error, got: %v", ttlErr)
		}
		if ttl <= 0 || ttl > time.Minute {
			t.Fatalf("expected TTL to be within the requested window, got: %v", ttl)
		}
	})

	t.Run("should overwrite an existing session with the same id", func(t *testing.T) {
		t.Parallel()
		client := newTestRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()
		sessionId := uniqueName("session")
		key := "session:" + sessionId

		if err := adapter.Save(ctx, sessionId, "first-user"); err != nil {
			t.Fatalf("first Save should not return an error, got: %v", err)
		}
		if err := adapter.Save(ctx, sessionId, "second-user"); err != nil {
			t.Fatalf("second Save should not return an error, got: %v", err)
		}

		value, getErr := client.Get(ctx, key).Result()
		if getErr != nil {
			t.Fatalf("expected key %q to exist, got error: %v", key, getErr)
		}
		if value != "second-user" {
			t.Fatalf("expected stored value to be overwritten to %q, got: %q", "second-user", value)
		}
	})

	t.Run("should return an internal error when Redis is unreachable", func(t *testing.T) {
		t.Parallel()
		client := newUnreachableRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()

		err := adapter.Save(ctx, uniqueName("session"), uniqueName("user-id"))

		if err == nil {
			t.Fatal("Save should return an error when Redis is unreachable")
		}
	})
}

func TestRedisAdapter_Delete(t *testing.T) {
	t.Run("should remove an existing session", func(t *testing.T) {
		t.Parallel()
		client := newTestRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()
		sessionId := uniqueName("session")
		key := "session:" + sessionId
		if err := client.Set(ctx, key, "user-id", time.Minute).Err(); err != nil {
			t.Fatalf("failed to seed session: %v", err)
		}

		err := adapter.Delete(ctx, sessionId)

		if err != nil {
			t.Fatalf("Delete should not return an error, got: %v", err)
		}
		exists, existsErr := client.Exists(ctx, key).Result()
		if existsErr != nil {
			t.Fatalf("Exists should not return an error, got: %v", existsErr)
		}
		if exists != 0 {
			t.Fatal("Delete should remove the session record")
		}
	})

	t.Run("should return a not found error when the session does not exist", func(t *testing.T) {
		t.Parallel()
		client := newTestRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()

		err := adapter.Delete(ctx, uniqueName("session"))

		if err == nil {
			t.Fatal("Delete should return an error when the session does not exist")
		}
	})

	t.Run("should return an internal error when Redis is unreachable", func(t *testing.T) {
		t.Parallel()
		client := newUnreachableRedisClient(t)
		adapter := database.NewRedisAdapter(client)
		ctx := context.Background()

		err := adapter.Delete(ctx, uniqueName("session"))

		if err == nil {
			t.Fatal("Delete should return an error when Redis is unreachable")
		}
	})
}
