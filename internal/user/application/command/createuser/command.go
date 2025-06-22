package createuser

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

func (c *Command) Name() string {
	return "createUser"
}

type CreatedCommand struct {
	Username string `json:"username"`
	Password string `json:"password"`
	//Person   any 	`json:"person"`
}

func NewCreatedCommand(Username, Password string) *CreatedCommand {
	return &CreatedCommand{
		Username,
		Password,
	}
}

func (c *CreatedCommand) Name() string {
	return "CreatedCommand"
}
