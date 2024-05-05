package core

import (
	uuidLib "github.com/google/uuid"
)

type Entity struct {
	uuid string
}

func NewEntity(uuid string) Entity {
	entity := Entity{}
	if uuid == "" {
		entity.uuid = uuidLib.NewString()
	}else{
		entity.uuid = uuid
	}
	return entity
}

func (entity *Entity) GetUuid() string {
	return entity.uuid
}