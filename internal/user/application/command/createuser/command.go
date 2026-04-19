package createuser

type Command struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	PersonName string `json:"name"`
	Document   string `json:"document"`
	BirthDate  string `json:"birthDate"`
	Email      string `json:"email"`
}

func (c *Command) Name() string {
	return "createUser"
}

