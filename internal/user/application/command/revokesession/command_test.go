// Revoke session command DTO tests.
//
// Intent: verify the revokeSession Command's routing identifier.
// Objective: cover the Name() method contract used by the command bus.
package revokesession_test

import (
	"testing"

	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/revokesession"
)

func TestCommand_Name(t *testing.T) {
	t.Run("should return revokeSession as the command name", func(t *testing.T) {
		t.Parallel()
		command := &revokesession.Command{
			SessionId: "session-id-value",
		}

		name := command.Name()

		if name != "revokeSession" {
			t.Fatalf("Name should return %q, got: %q", "revokeSession", name)
		}
	})
}
