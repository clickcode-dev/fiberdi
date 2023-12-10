package fiberdi

import (
	"fmt"
	"os"
	"reflect"

	"github.com/Goldziher/go-utils/sliceutils"

	"github.com/clickcode-dev/fiberdi/flog"
	"github.com/gofiber/fiber/v2"

	"github.com/charmbracelet/log"

	"github.com/sarulabs/di"
)

type IModule interface {
	setExportsNames(names []string)
	setInjectableCurrentInAppModuleFn(fn func(string) []string)
	setInjectablesNames(names []string)
	addContainer(container di.Container)
	setParent(module *Module)
	setInjectableCurrent(name string) []string
	mappedInjectables(modules []IModule, names []string) []string
	start(app *fiber.App) *fiber.App
	verifyIfIsAnAttemptToInjectController(structt interface{}, inject interface{})
	configureThirdDependecies(pointer reflect.Value, inject interface{}) interface{}
	injectThirdDependencies()
	isAttemptToImportSomethingNotExported(name string, injecteds []string) bool
	getDependency(container di.Container, pointer reflect.Value, name string) (interface{}, error)
	injectDependencies(ctn di.Container, field reflect.StructField, logger *log.Logger, pointer reflect.Value, structt interface{}, valueOf reflect.Value) reflect.Value
	buildDependency(app *fiber.App, structt interface{}, injecteds *[]string) *fiber.App
	getStructTag(f reflect.StructField, tagName string) string
	generateDependecies(app *fiber.App) (*fiber.App, di.Container)
	mappedExports(modules []IModule, names []string) []string
}

type Module struct {
	Controllers []IController
	Injectables []interface{}
	Modules     []IModule

	Imports []IModule
	Exports []interface{}

	injectablesNames []string
	exportsNames     []string

	currentInjectable []string

	setInjectableCurrentInAppModule func(string) []string

	builder          *di.Builder
	parentContainers []di.Container
	parentModule     IModule
	container        di.Container
}

func (m *Module) addContainer(container di.Container) {
	m.parentContainers = append(m.parentContainers, container)
}

func (m *Module) setParent(module *Module) {
	m.parentModule = module
}

func (m *Module) setInjectableCurrentInAppModuleFn(fn func(string) []string) {
	m.setInjectableCurrentInAppModule = fn
}

func (m *Module) setInjectableCurrent(name string) []string {
	m.currentInjectable = append(m.currentInjectable, name)

	return m.currentInjectable
}

func (m *Module) setInjectablesNames(names []string) {
	m.injectablesNames = names
}

func (m *Module) setExportsNames(names []string) {
	m.exportsNames = names
}

func (m *Module) mappedInjectables(modules []IModule, names []string) []string {
	for _, module := range modules {
		for _, injectable := range module.(*Module).Injectables {
			valueOf := reflect.ValueOf(injectable)

			if valueOf.Kind() != reflect.Ptr {
				panic(fmt.Errorf("are you sure? maybe %T is not a pointer, please check", injectable))
			}

			pointer := reflect.Indirect(valueOf)
			name := pointer.Type().Name()

			if sliceutils.FindIndexOf(names, name) == -1 {
				names = append(names, name)
			}
		}

		names = module.mappedInjectables(module.(*Module).Modules, names)
	}

	return names
}

func (m *Module) mappedExports(modules []IModule, names []string) []string {
	for _, module := range modules {
		for _, exports := range module.(*Module).Exports {
			valueOf := reflect.ValueOf(exports)

			if valueOf.Kind() != reflect.Ptr {
				panic(fmt.Errorf("are you sure? maybe %T is not a pointer, please check", exports))
			}

			pointer := reflect.Indirect(valueOf)
			name := pointer.Type().Name()

			if sliceutils.FindIndexOf(names, name) == -1 {
				names = append(names, name)
			} else {
				panic(fmt.Errorf("are you trying to export %s in two modules?", name))
			}
		}

		names = module.mappedExports(module.(*Module).Modules, names)
	}

	return names
}

func (m *Module) start(app *fiber.App) *fiber.App {
	builder, err := di.NewBuilder(di.Request)

	if err != nil {
		panic(err)
	}

	m.builder = builder

	app, container := m.generateDependecies(app)

	m.container = container

	for _, module := range m.Modules {
		module.setInjectableCurrentInAppModuleFn(m.setInjectableCurrent)
		module.addContainer(container)

		for _, importt := range module.(*Module).Imports {
			module.addContainer(importt.(*Module).container)
		}

		module.setParent(m)
		app = module.start(app)
	}

	return app
}

func (m *Module) verifyIfIsAnAttemptToInjectController(structt interface{}, inject interface{}) {
	_, structtIsController := structt.(IController)
	_, injectIsController := inject.(IController)

	if injectIsController && structtIsController {
		panic(fmt.Errorf("you cannot inject %T inside another controller", inject))
	}
}

func (m *Module) configureThirdDependecies(pointer reflect.Value, inject interface{}) interface{} {
	if _, ok := inject.(*log.Logger); ok {
		logger := flog.NewLogger(ternary(os.Getenv("FIBER_MODE") != "release", log.DebugLevel, log.InfoLevel))
		logger.SetPrefix(pointer.Type().Name())

		return logger
	}

	return inject
}

func (m *Module) injectThirdDependencies() {
	m.builder.Add(di.Def{
		Name:  "Logger",
		Scope: di.Request,
		Build: func(ctn di.Container) (interface{}, error) {
			return flog.NewLogger(ternary(os.Getenv("FIBER_MODE") != "release", log.DebugLevel, log.InfoLevel)), nil
		},
	})
}

func (m *Module) getDependency(container di.Container, pointer reflect.Value, name string) (interface{}, error) {
	var err error
	var inject interface{}

	for _, parent := range m.parentContainers {
		if inject != nil {
			break
		}

		if parent == nil {
			continue
		}

		inject, _ = parent.SafeGet(name)
	}

	if inject == nil {
		inject, err = container.SafeGet(name)
	}

	inject = m.configureThirdDependecies(pointer, inject)

	return inject, err
}

func (m *Module) injectDependencies(
	ctn di.Container,
	field reflect.StructField,
	logger *log.Logger,
	pointer reflect.Value,
	structt interface{},
	valueOf reflect.Value,
) reflect.Value {
	if field.Type.Kind() != reflect.Ptr {
		panic(fmt.Errorf("are you trying to inject %s as a non-pointer dependency?", field.Type.Name()))
	}

	injectName := field.Type.Elem().Name()

	inject, err := m.getDependency(ctn, pointer, injectName)

	if err != nil {
		panic(fmt.Errorf("%v\n\nTIP:\n - Are you trying to access %s inside %T before inject in module?\n - Are you trying to inject a dependency that was supposed to be ignored? If so, remember to use di:\"ignore\"", err, injectName, structt))
	}

	if hook, ok := inject.(IPostConstruct); ok {
		hook.PostConstruct()
	}

	injectValueOf := reflect.ValueOf(inject)

	m.verifyIfIsAnAttemptToInjectController(structt, inject)

	return injectValueOf
}

func (m *Module) isAttemptToImportSomethingNotExported(name string, injecteds []string) bool {
	filterInjectable := filter(injecteds, func(inject string) bool {
		return inject == name
	})

	return len(filterInjectable) > 1
}

func (m *Module) buildDependency(app *fiber.App, structt interface{}, injecteds *[]string) *fiber.App {
	valueOf := reflect.ValueOf(structt)

	if valueOf.Kind() != reflect.Ptr {
		panic(fmt.Errorf("are you sure? maybe %T is not a pointer, please check", structt))
	}

	pointer := reflect.Indirect(valueOf)

	name := pointer.Type().Name()

	if injecteds != nil && m.isAttemptToImportSomethingNotExported(name, *injecteds) {
		return app
	}

	if hook, ok := structt.(IPreConstruct); ok {
		hook.PreConstruct()
	}

	logger := flog.NewLogger(ternary(os.Getenv("FIBER_MODE") != "release", log.DebugLevel, log.InfoLevel))

	logger.SetPrefix(name)

	m.builder.Add(di.Def{
		Name:  name,
		Scope: di.Request,
		Build: func(ctn di.Container) (interface{}, error) {
			for i := 0; i < pointer.NumField(); i++ {
				field := pointer.Type().Field(i)

				tagDi := m.getStructTag(field, "di")

				if tagDi == "ignore" {
					logger.Debugf("Ignoring field %s because tag ignore is set", field.Name)
					continue
				}

				injectValueOf := m.injectDependencies(ctn, field, logger, pointer, structt, valueOf)

				valueOf.Elem().Field(i).Set(injectValueOf)
			}

			if controller, ok := structt.(IController); ok {
				app = controller.Routes(app)
			}

			logger.Debugf("dependencies initialized")

			if hook, ok := structt.(IPostConstruct); ok {
				hook.PostConstruct()
			}

			return structt, nil
		},
	})

	return app
}

func (m *Module) getStructTag(f reflect.StructField, tagName string) string {
	return string(f.Tag.Get(tagName))
}

func (m *Module) generateDependecies(app *fiber.App) (*fiber.App, di.Container) {
	m.injectThirdDependencies()

	for _, injectable := range m.Injectables {
		valueOf := reflect.ValueOf(injectable)
		pointer := reflect.Indirect(valueOf)

		current := m.setInjectableCurrentInAppModule(pointer.Type().Name())

		app = m.buildDependency(app, injectable, &current)
	}

	for _, controller := range m.Controllers {
		app = m.buildDependency(app, controller, nil)
	}

	container := m.builder.Build()

	for _, controller := range m.Controllers {
		typee := reflect.Indirect(reflect.ValueOf(controller)).Type()

		container.Get(typee.Name())
	}

	return app, container
}
