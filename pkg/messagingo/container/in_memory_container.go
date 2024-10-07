package container

type inMemoryContainer struct {
	objects map[string]any
}

var instance *inMemoryContainer

// TODO: check is necessary singleton instance here
func CreateInMemoryContainer() *inMemoryContainer {
	if instance == nil {
		instance = &inMemoryContainer{
			objects: make(map[string]any),
		}
	}

	return instance
}

func (instance *inMemoryContainer) Get(id string) any {
	return instance.objects[id]
}

func (instance *inMemoryContainer) Has(id string) bool {
	return instance.objects[id] != nil
}

func (instance *inMemoryContainer) Set(id string, value any) {
	instance.objects[id] = value
}
