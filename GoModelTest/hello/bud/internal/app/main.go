package main

import (
	context "context"
	errors "errors"
	budclient "github.com/livebud/bud/package/budclient"
	commander "github.com/livebud/bud/package/commander"
	gomod "github.com/livebud/bud/package/gomod"
	log "github.com/livebud/bud/package/log"
	console "github.com/livebud/bud/package/log/console"
	filter "github.com/livebud/bud/package/log/filter"
	overlay "github.com/livebud/bud/package/overlay"
	router "github.com/livebud/bud/package/router"
	os "os"
	controller "temp/GoModelTest/hello/bud/internal/app/controller"
	public "temp/GoModelTest/hello/bud/internal/app/public"
	web "temp/GoModelTest/hello/bud/internal/app/web"
)

func main() {
	os.Exit(run(context.Background(), os.Args[1:]...))
}

// Run the cli
func run(ctx context.Context, args ...string) int {
	if err := parse(ctx, args...); err != nil {
		if errors.Is(err, context.Canceled) {
			return 0
		}
		console.Error(err.Error())
		return 1
	}
	return 0
}

// Parse the arguments
func parse(ctx context.Context, args ...string) error {
	cli := commander.New("bud")
	app := new(App)
	cli.Flag("listen", "address to listen to").String(&app.Listen).Default(":3000")
	cli.Flag("log", "filter logs with a pattern").Short('L').String(&app.Log).Default("info")
	cli.Run(app.Run)
	return cli.Parse(ctx, args)
}

// App command
type App struct {
	Listen string
	Log    string
}

// logger creates a structured log that supports filtering
func (a *App) logger() (log.Interface, error) {
	handler, err := filter.Load(console.New(os.Stderr), a.Log)
	if err != nil {
		return nil, err
	}
	return log.New(handler), nil
}

// Run your app
func (a *App) Run(ctx context.Context) error {
	log, err := a.logger()
	if err != nil {
		return err
	}
	budClient, err := budclient.Try(os.Getenv("BUD_LISTEN"))
	if err != nil {
		return err
	}
	// Load the module dependency
	module, err := gomod.Find(".")
	if err != nil {
		return err
	}
	// Load the web server
	webServer, err := loadWeb(
		module, log,
	)
	if err != nil {
		return err
	}
	// Inform bud that we're ready
	budClient.Publish("app:ready", nil)
	// Start serving requests
	log.Debug("app: listening on", "listen", a.Listen)
	return webServer.Serve(ctx, a.Listen)
}

func loadWeb(gomodModule *gomod.Module, logInterface log.Interface) (*web.Server, error) {
	routerRouter := router.New()
	controllerAboutIndexAction := &controller.AboutIndexAction{}
	controllerAboutController := &controller.AboutController{Index: controllerAboutIndexAction}
	controllerPostController := &controller.PostController{}
	controllerUsersIndexAction := &controller.UsersIndexAction{}
	controllerUsersNewAction := &controller.UsersNewAction{}
	controllerUsersCreateAction := &controller.UsersCreateAction{}
	controllerUsersShowAction := &controller.UsersShowAction{}
	controllerUsersUpdateAction := &controller.UsersUpdateAction{}
	controllerUsersDeleteAction := &controller.UsersDeleteAction{}
	controllerUsersEditAction := &controller.UsersEditAction{}
	controllerUsersController := &controller.UsersController{Index: controllerUsersIndexAction, New: controllerUsersNewAction, Create: controllerUsersCreateAction, Show: controllerUsersShowAction, Update: controllerUsersUpdateAction, Delete: controllerUsersDeleteAction, Edit: controllerUsersEditAction}
	controllerController := &controller.Controller{About: controllerAboutController, Post: controllerPostController, Users: controllerUsersController}
	overlayServer, err := overlay.Serve(logInterface, gomodModule)
	if err != nil {
		return nil, err
	}
	publicMiddleware := public.New(overlayServer)
	webServer := web.New(routerRouter, controllerController, publicMiddleware)
	return webServer, err
}
