package messagingo

import (
	"fmt"

	"github.com/hex-api-go/pkg/messagingo/container"
)

type ConfiguredMessageSystem struct {
	container container.Container
}

func (m *ConfiguredMessageSystem) GetCommandBus() {}

func (m *ConfiguredMessageSystem) GetEventBus() {
	fmt.Println("event bus CALLED", m.container)
}
