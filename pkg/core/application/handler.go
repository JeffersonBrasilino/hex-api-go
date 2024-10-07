package application

import "github.com/hex-api-go/pkg/core/domain"

type Handler[T domain.Action[any]] interface {
	Handler(action T)
}

type CommandHandler[T domain.Action[any]] interface {
	Handler[T]
}
