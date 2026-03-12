package v1

import (
	"github.com/kartikey315/go-tasker/internal/handler"
	"github.com/kartikey315/go-tasker/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerCommentRoutes(r *echo.Group, h *handler.CommentHandler, auth *middleware.AuthMiddleware) {
	// Comment operations
	comments := r.Group("/comments")
	comments.Use(auth.RequireAuth)

	// Individual comment operations
	dynamicComments := comments.Group("/:id")
	dynamicComments.PATCH("", h.UpdateComment)
	dynamicComments.DELETE("", h.DeleteComment)
}
