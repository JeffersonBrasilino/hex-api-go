package domain

/*
Abstração de Entidade de Domínio.

Uma Entidade de Domínio é o local onde as regras de negócio residem.
as regras implementadas na entidade devem ser totalmente agnósticas
de detalhes de implementeções(banco de dados, integrações de terceiros e etc).

ex:

	type item struct{
		*Entity
	}

	var itemInstance = &item{
		Entity:     domain.NewEntity(""),
	}
*/
type Entity struct {
	uuid string
}

// cria a instancia de entidade de dominio
func NewEntity(uuid string) *Entity {
	return &Entity{
		uuid: uuid,
	}
}

func (entity *Entity) Uuid() string {
	return entity.uuid
}
