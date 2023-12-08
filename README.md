# Fiber Dependency Injection

FiberDI uses Dependency Injection in conjunction with the practices used at ClickCode.dev

## Get started

    go get -u github.com/clickcode-dev/fiberdi

### Create a server

    func main() {
    	appModule := &fiberdi.Module{}

    	app := fiberdi.New(appModule, fiber.Config{})

    	port := "3000"

    	log.Infof("Server is running on port %s", port)

    	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
    		log.Fatal(err)
    	}
    }

### Create a module

    type Controller struct {}

    func (controller *Controller) Routes(app *fiber.App) *fiber.App {
        app.Get("/', ...)
        return app
    }

    func main() {
    	appModule := &fiberdi.Module{
    		Controllers: []fiberdi.IController{
    			&Controller{},
    		},
    	}
    	app := fiberdi.New(appModule, fiber.Config{})

    	...continue
    }

### Create an injectable

    type Service struct {}

    func (service *Service) func HelloWorld(ctx *fiber.Ctx) {
    	return ctx.JSON("Hello World")
    }

    type Controller struct {
    	Service *Service
    }

    // You need the function Routes to create your paths
    func (controller *Controller) Routes(app *fiber.App) *fiber.App {
        app.Get("/', controller.Service.HelloWorld)
        return app
    }

    func main() {
    	appModule := &fiberdi.Module{
    		Controllers: []fiberdi.IController{
    			&Controller{},
    		},
    		Injectables: []interface{}{
    			&Service{},
    		},
    	}
    	app := fiberdi.New(appModule, fiber.Config{})

    	...continue
    }

### Ignore a field

    type Controller struct {
    	Service *Service
    Name string `di:"ignore"`
    }

##### \*Field Name will be ignore

### Import an injectable of another module

    catModule := &fiberdi.Module{
    	Controllers: []fiberdi.IController{
    		&CatController{},
    	},
    	Injectables: []interface{}{
    		&CatService{},
    	},
    	Exports: []interface{}{
    		&CatService{},
    	},
    }

    dogModule := &fiberdi.Module{
    	Controllers: []fiberdi.IController{
    		&DogController{},
    	},
    	Injectables: []interface{}{
    		&CatService{},
    	},
    	Imports: []interface{}{
    		catModule,
    	},
    }

### Hooks

#### Post Construct

Will be executed after inject dependencies
type Service struct {
Example Example
}
func (service \*Service) PostConstruct() {
log.Print(service.Example == nil) // false
}

#### Pre Construct

Will be executed before inject dependencies
type Service struct {
Example *Example
}
func (service *Service) PreConstruct() {
log.Print(service.Example == nil) // true
}

### Utilities

    import "github.com/clickcode-dev/fiberdi/log"

    type Service struct {
    	Log *log.Logger
    }

    func (service *Service) PreConstruct() {
    	log.Print("Wow! It's works")
    }

### Run our example

    go run .examples/main.go
