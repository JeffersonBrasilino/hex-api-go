// Revoke session command DTO.
//
// Intent: carry the data required for an administrator to revoke a session.
// Objective: transport the session identifier submitted by the client to the revoke session use
// case handler.
package revokesession

// Command holds the identifier of the session to be revoked.
type Command struct {
	SessionId string `json:"sessionId"`
}

// Name returns the command's identifier, used for routing by the command bus.
func (c *Command) Name() string {
	return "revokeSession"
}
