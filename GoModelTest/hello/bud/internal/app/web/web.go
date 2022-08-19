package web

// GENERATED. DO NOT EDIT.

import (
	context "context"
	webrt "github.com/livebud/bud/framework/web/webrt"
	middleware "github.com/livebud/bud/package/middleware"
	router "github.com/livebud/bud/package/router"
	http "net/http"
	controller "temp/GoModelTest/hello/bud/internal/app/controller"
	public "temp/GoModelTest/hello/bud/internal/app/public"
)

// New web server
func New(
	router *router.Router,
	controller *controller.Controller,
	public public.Middleware,
) *Server {
	// Action routing
	router.Get(`/about`, controller.About.Index)
	router.Get(`/users`, controller.Users.Index)
	router.Get(`/users/new`, controller.Users.New)
	router.Post(`/users`, controller.Users.Create)
	router.Get(`/users/:id`, controller.Users.Show)
	router.Patch(`/users/:id`, controller.Users.Update)
	router.Delete(`/users/:id`, controller.Users.Delete)
	router.Get(`/users/:id/edit`, controller.Users.Edit)
	// Compose the middleware together
	middleware := middleware.Compose(
		router,
		public,
	)
	// 404 at the bottom of the middleware
	handler := middleware.Middleware(http.NotFoundHandler())
	return &Server{handler}
}

type Server struct {
	http.Handler
}

func (s *Server) Serve(ctx context.Context, address string) error {
	listener, err := webrt.Listen("WEB", address)
	if err != nil {
		return err
	}
	return webrt.Serve(ctx, listener, s)
}
