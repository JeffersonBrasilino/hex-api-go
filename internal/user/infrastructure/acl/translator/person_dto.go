package translator

type PersonDto struct {
	Name      string `json:"name"`
	Age       int    `json:"age"`
	BirthDate string `json:"birth_date"`
	Email     string `json:"email"`
}
