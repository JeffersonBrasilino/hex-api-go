package bus

type CommandBus interface {
	Send(queueOrTopicName string, payload []byte) (any, error)
}
