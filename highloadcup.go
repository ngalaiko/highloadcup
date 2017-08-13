package highloadcup

import (
	"context"
	"log"

	"github.com/ngalayko/highloadcup/config"
	"github.com/ngalayko/highloadcup/database"
	"github.com/ngalayko/highloadcup/web"
	"github.com/ngalayko/highloadcup/views"
)

type Application struct {
	ctx context.Context

	config *config.Config
}

var (
	services = []func(context.Context, interface{}) context.Context{
		config.NewContext,
		database.NewContext,
		web.NewContext,
		views.NewContext,
	}
)

func NewApp() *Application {
	app := &Application{
		ctx: context.Background(),
	}

	app.initServices()
	app.initData()

	return app
}

func (app *Application) ServeHTTP() error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered", r)
		}
	}()

	return web.FromContext(app.ctx).ServeHTTP()
}

func (app *Application) initData() {
	app.config = config.FromContext(app.ctx)

	db := database.FromContext(app.ctx)

	if err := db.ParseData(app.config.DataPath); err != nil {
		log.Panic(err)
	}
}

func (app *Application) initServices() {
	for _, service := range services {
		app.ctx = service(app.ctx, nil)
	}
}
