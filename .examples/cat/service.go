package cat

import "github.com/charmbracelet/log"

type CatService struct {
	Logger *log.Logger
}

func (service CatService) HelloWorld() string {
	service.Logger.Info("testando")

	return "Cat"
}
