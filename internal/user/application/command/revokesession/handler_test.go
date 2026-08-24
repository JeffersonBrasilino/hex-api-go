// Revoke session command handler tests.
//
// Intent: verify the administrative session revocation use case orchestration in isolation from
// real infrastructure.
// Objective: cover the RF-04 acceptance scenarios — successful revocation and a non-existent
// session — using a hand-rolled test double for the revoke session store contract.
package revokesession_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jeffersonbrasilino/ddgo"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/revokesession"
)

// stubRevokeSessionStore is a hand-rolled test double implementing contract.RevokeSessionStore.
type stubRevokeSessionStore struct {
	deleteErr        error
	deleteCalledWith string
}

func (s *stubRevokeSessionStore) Delete(ctx context.Context, sessionId string) error {
	s.deleteCalledWith = sessionId
	return s.deleteErr
}

func validCommand() *revokesession.Command {
	return &revokesession.Command{SessionId: "a-session-id"}
}

func TestNewCommandHandler(t *testing.T) {
	t.Run("should return a non-nil Handler", func(t *testing.T) {
		t.Parallel()
		handler := revokesession.NewCommandHandler(&stubRevokeSessionStore{})

		if handler == nil {
			t.Fatal("NewCommandHandler should return a non-nil handler")
		}
	})
}

func TestHandler_Handle(t *testing.T) {
	t.Run("should return a NotFoundError when the session does not exist", func(t *testing.T) {
		t.Parallel()
		revokeSessionStore := &stubRevokeSessionStore{deleteErr: errors.New("session not found")}
		handler := revokesession.NewCommandHandler(revokeSessionStore)

		result, err := handler.Handle(context.Background(), validCommand())

		var notFoundErr *ddgo.NotFoundError
		if !errors.As(err, &notFoundErr) {
			t.Fatalf("Handle should return a *ddgo.NotFoundError, got: %T (%v)", err, err)
		}
		if result != nil {
			t.Fatalf("Handle should return a nil result on error, got: %v", result)
		}
	})

	t.Run("should revoke the session on success", func(t *testing.T) {
		t.Parallel()
		revokeSessionStore := &stubRevokeSessionStore{}
		handler := revokesession.NewCommandHandler(revokeSessionStore)
		command := validCommand()

		result, err := handler.Handle(context.Background(), command)

		if err != nil {
			t.Fatalf("Handle should not return an error, got: %v", err)
		}
		if result != "session revoked" {
			t.Fatalf("Handle should return %q, got: %v", "session revoked", result)
		}
		if revokeSessionStore.deleteCalledWith != command.SessionId {
			t.Fatalf("Delete should be called with the submitted session id, got: %q", revokeSessionStore.deleteCalledWith)
		}
	})
}
