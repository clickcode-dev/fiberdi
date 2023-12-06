package dog

type DogService struct {
	Name string `di:"ignore"`
}

func (service DogService) HelloWorld() string {
	return "Dog"
}
