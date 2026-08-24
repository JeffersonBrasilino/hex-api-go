// Login command DTO.
//
// Intent: carry the credentials required to authenticate a user.
// Objective: transport the username and password submitted by the client to the login use case
// handler.
package login

// Command holds the credentials submitted for a login attempt.
type Command struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Name returns the command's identifier, used for routing by the command bus.
func (c *Command) Name() string {
	return "login"
}
