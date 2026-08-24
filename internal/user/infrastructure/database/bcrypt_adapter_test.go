package database_test

import (
	"strings"
	"testing"

	"github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/database"
)

func TestNewBcryptAdapter(t *testing.T) {
	t.Run("should return a non-nil BcryptAdapter", func(t *testing.T) {
		t.Parallel()
		adapter := database.NewBcryptAdapter()
		if adapter == nil {
			t.Fatal("NewBcryptAdapter should return a non-nil adapter")
		}
	})
}

func TestBcryptAdapter_Hash(t *testing.T) {
	t.Run("should return a hash different from the raw password", func(t *testing.T) {
		t.Parallel()
		adapter := database.NewBcryptAdapter()
		raw := "my-secret-password"

		hash, err := adapter.Hash(raw)

		if err != nil {
			t.Fatalf("Hash should not return an error, got: %v", err)
		}
		if hash == "" {
			t.Fatal("Hash should return a non-empty hash")
		}
		if hash == raw {
			t.Fatal("Hash should not return the raw password unchanged")
		}
	})

	t.Run("should return different hashes for the same password on each call", func(t *testing.T) {
		t.Parallel()
		adapter := database.NewBcryptAdapter()
		raw := "my-secret-password"

		hash1, err1 := adapter.Hash(raw)
		hash2, err2 := adapter.Hash(raw)

		if err1 != nil || err2 != nil {
			t.Fatalf("Hash should not return an error, got: %v / %v", err1, err2)
		}
		if hash1 == hash2 {
			t.Fatal("Hash should produce different salted hashes across calls")
		}
	})

	t.Run("should return an error when raw password exceeds bcrypt's length limit", func(t *testing.T) {
		t.Parallel()
		adapter := database.NewBcryptAdapter()
		raw := strings.Repeat("a", 73)

		hash, err := adapter.Hash(raw)

		if err == nil {
			t.Fatal("Hash should return an error for a password longer than 72 bytes")
		}
		if hash != "" {
			t.Fatalf("Hash should return an empty hash on error, got: %q", hash)
		}
	})
}

func TestBcryptAdapter_Verify(t *testing.T) {
	t.Run("should return true when raw password matches the hash", func(t *testing.T) {
		t.Parallel()
		adapter := database.NewBcryptAdapter()
		raw := "my-secret-password"
		hash, err := adapter.Hash(raw)
		if err != nil {
			t.Fatalf("Hash should not return an error, got: %v", err)
		}

		matches, err := adapter.Verify(raw, hash)

		if err != nil {
			t.Fatalf("Verify should not return an error, got: %v", err)
		}
		if !matches {
			t.Fatal("Verify should return true when the raw password matches the hash")
		}
	})

	t.Run("should return false without error when raw password does not match the hash", func(t *testing.T) {
		t.Parallel()
		adapter := database.NewBcryptAdapter()
		hash, err := adapter.Hash("my-secret-password")
		if err != nil {
			t.Fatalf("Hash should not return an error, got: %v", err)
		}

		matches, err := adapter.Verify("wrong-password", hash)

		if err != nil {
			t.Fatalf("Verify should not return an error for a simple mismatch, got: %v", err)
		}
		if matches {
			t.Fatal("Verify should return false when the raw password does not match the hash")
		}
	})

	t.Run("should return an error when the hash is malformed", func(t *testing.T) {
		t.Parallel()
		adapter := database.NewBcryptAdapter()

		matches, err := adapter.Verify("my-secret-password", "not-a-valid-bcrypt-hash")

		if err == nil {
			t.Fatal("Verify should return an error for a malformed hash")
		}
		if matches {
			t.Fatal("Verify should return false when it errors")
		}
	})
}
