package bus

type GenericCommandBus struct{}

func NewGenericCommandBus() *GenericCommandBus {
	return &GenericCommandBus{}
}

func (c *GenericCommandBus) Send(queueOrTopicName string, payload []byte) (any, error) {

	return nil, nil
}
