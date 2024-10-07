package container

type Container interface {
	Has(id string) bool
	Get(id string) any
	Set(id string, value any)
}
