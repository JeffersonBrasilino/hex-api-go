package createuser

import "github.com/hex-api-go/pkg/core/infrastructure/message_system/message"

type Command struct {
	Username string `json:"username"`
	Password string `json:"password"`
	//Person   any 	`json:"person"`
}

func CreateCommand(Username, Password string) *Command {
	return &Command{
		Username,
		Password,
	}
}
func (c *Command) Type() message.MessageType {
	return message.Command
}

func (c *Command)Name() string {
	return "createUser"
}