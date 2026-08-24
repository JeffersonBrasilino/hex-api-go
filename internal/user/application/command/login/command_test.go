// Login command DTO tests.
//
// Intent: verify the login Command's routing identifier.
// Objective: cover the Name() method contract used by the command bus.
package login_test

import (
	"testing"

	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/login"
)

func TestCommand_Name(t *testing.T) {
	t.Run("should return login as the command name", func(t *testing.T) {
		t.Parallel()
		command := &login.Command{
			Username: "jdoe",
			Password: "secret",
		}

		name := command.Name()

		if name != "login" {
			t.Fatalf("Name should return %q, got: %q", "login", name)
		}
	})
}
