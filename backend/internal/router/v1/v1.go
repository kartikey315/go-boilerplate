package v1

import (
	"github.com/kartikey315/go-tasker/internal/handler"
	"github.com/kartikey315/go-tasker/internal/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterV1Routes(routes *echo.Group, handlers *handler.Handlers, middleware *middleware.Middlewares) {

	// Register todo routes
	registerTodoRoutes(routes, handlers.TodoHandler, handlers.CommentHandler, middleware.Auth)

	// Register comment routes
	registerCommentRoutes(routes, handlers.CommentHandler, middleware.Auth)

	// Register category routes
	registerCategoryRoutes(routes, handlers.CategoryHandler, middleware.Auth)
}
