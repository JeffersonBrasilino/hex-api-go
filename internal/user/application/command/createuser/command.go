package createuser

type Command struct {
	Username string `json:"username"`
	Password string `json:"password"`
	//Person   any 	`json:"person"`
}

func (c *Command) Payload() any {
	return c
}

func (c *Command) Headers() any {
	return nil
}

func CreateCommand(Username, Password string) *Command {
	return &Command{
		Username,
		Password,
	}
}
