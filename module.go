package fiberdi

import (
	"fmt"
	"log"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di"
)

type IModule interface{}

type Module struct {
	Controllers []IController
	Injectables []interface{}
	Modules     []Module
	Exports     []interface{}

	injectables *[]string
	exportsName *[]string
	exports     *[]interface{}
}

func (m Module) start() {
	builder := m.generateDependecies()

	builder.Build()
}

func (m Module) buildDependency(dep interface{}, builder *di.Builder) {
	valueOf := reflect.ValueOf(dep)

	if valueOf.Kind() != reflect.Ptr {
		panic(fmt.Errorf("are you sure? maybe %T is not a pointer, please check", dep))
	}

	pointer := reflect.Indirect(valueOf)

	builder.Add(di.Def{
		Name:  pointer.Type().Name(),
		Scope: di.Request,
		Build: func(ctn di.Container) (interface{}, error) {

			return dep, nil
		},
	})
}

func (m Module) generateDependecies() *di.Builder {
	builder, err := di.NewBuilder(di.Request)

	if err != nil {
		panic(err)
	}

	for _, controller := range m.Controllers {
		m.buildDependency(controller, builder)
	}

	return builder
}

func (m Module) foundInjectables(appModule *Module, injectables *[]string) *[]string {
	for _, module := range m.Modules {
		for _, inject := range m.Injectables {
			typee := reflect.Indirect(reflect.ValueOf(inject)).Type()

			injectables = array(append(*injectables, typee.Name()))
		}

		injectables = module.foundInjectables(appModule, injectables)
	}

	appModule.injectables = injectables

	return appModule.injectables
}

func (m Module) foundExports(appModule *Module, exportsName *[]string, exports *[]interface{}) (*[]string, *[]interface{}) {
	for _, module := range m.Modules {
		for _, export := range m.Exports {
			typee := reflect.Indirect(reflect.ValueOf(export)).Type()

			exports = array(append(*exports, export))
			exportsName = array(append(*exportsName, typee.Name()))
		}

		exportsName, exports = module.foundExports(appModule, exportsName, exports)
	}

	appModule.exports = exports
	appModule.exportsName = exportsName

	return appModule.exportsName, appModule.exports
}

func (m Module) addDependencies(app *fiber.App) *fiber.App {
	builder, err := di.NewBuilder(di.Request)

	if err != nil {
		panic(err)
	}

	for _, module := range m.Modules {
		module.exportsName = m.exportsName
		module.exports = m.exports

		module.foundInjectables(&module, new([]string))
		module.foundExports(&module, new([]string), new([]interface{}))

		app = module.addDependencies(app)
	}

	for _, injectable := range m.Injectables {
		typee := reflect.Indirect(reflect.ValueOf(injectable)).Type()

		builder.Add(di.Def{
			Name:  typee.Name(),
			Scope: di.Request,
			Build: func(ctn di.Container) (interface{}, error) {
				ref := reflect.Indirect(reflect.ValueOf(injectable))

				for i := 0; i < ref.NumField(); i++ {
					name := ref.Type().Field(i).Name

					inject, err := ctn.SafeGet(name)

					if err != nil {
						panic(fmt.Errorf("are you trying to access %s inside %T before inject in module?\n\n%v", name, injectable, err))
					}

					reflect.ValueOf(injectable).Elem().Field(i).Set(reflect.ValueOf(inject))
				}

				return injectable, nil
			},
		})
	}

	for _, export := range m.Exports {
		typee := reflect.Indirect(reflect.ValueOf(export)).Type()

		builder.Add(di.Def{
			Name:  typee.Name(),
			Scope: di.Request,
			Build: func(ctn di.Container) (interface{}, error) {
				return export, nil
			},
		})
	}

	for _, controller := range m.Controllers {
		v := reflect.ValueOf(controller)

		if v.Kind() != reflect.Ptr {
			panic(fmt.Errorf("are you sure? Maybe %T is not a pointer, please check", controller))
		}

		typee := reflect.Indirect(v).Type()

		builder.Add(di.Def{
			Name:  typee.Name(),
			Scope: di.Request,
			Build: func(ctn di.Container) (interface{}, error) {
				ref := reflect.Indirect(reflect.ValueOf(controller))

				for i := 0; i < ref.NumField(); i++ {
					name := ref.Type().Field(i).Name

					inject, err := ctn.SafeGet(name)

					if err != nil {
						panic(fmt.Errorf("are you trying to access %s inside %T before inject in module?\n\n%v", name, controller, err))
					}

					log.Println(name, reflect.Indirect(reflect.ValueOf(inject)).Type().Name())

					reflect.ValueOf(controller).Elem().Field(i).Set(reflect.ValueOf(inject))
				}

				app = controller.Routes(app)

				return controller, nil
			},
		})
	}

	build := builder.Build()

	for _, controller := range m.Controllers {
		typee := reflect.Indirect(reflect.ValueOf(controller)).Type()

		build.Get(typee.Name())
	}

	return app
}
