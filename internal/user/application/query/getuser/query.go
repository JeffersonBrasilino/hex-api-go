package getuser

type Query struct {
	DataSource    string
	mapMockedData []float64
}

func NewQuery() *Query {
	return &Query{}
}

func (c *Query) Name() string {
	return "getUser"
}
