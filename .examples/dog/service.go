package dog

import "github.com/charmbracelet/log"

type DogService struct {
	Logger *log.Logger

	Name string `di:"ignore"`
}

func (service *DogService) HelloWorld() string {
	return "Dog"
}
