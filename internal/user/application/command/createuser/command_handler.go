package createuser

import (
	"fmt"
	"sync"
	"time"

	"github.com/hex-api-go/internal/user/domain"
	"github.com/hex-api-go/internal/user/domain/contract"
)

type CommandHandler struct {
	repository contract.UserRepository
}
type Response struct {
	TestReturn string
}

func NewComandHandler(repository contract.UserRepository) *CommandHandler {
	return &CommandHandler{repository}
}

func (c *CommandHandler) Handle(data *Command) (any, error) {
	var wg sync.WaitGroup
	fmt.Println("start processes create user")
	time.Sleep(time.Second * 2)
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("domain.NewUser called")
		domain.NewUser("new user", "new Password")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("c.repository.Create called")
		c.repository.Create(domain.NewUser("", ""))
	}()

	wg.Wait()
	return &Response{TestReturn: "RETURNED OK"}, nil
}
