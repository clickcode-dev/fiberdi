package fiberdi

type IPostConstruct interface {
	PostConstruct()
}

type IPreConstruct interface {
	PreConstruct()
}
