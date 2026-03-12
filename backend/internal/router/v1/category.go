package v1

import (
	"github.com/kartikey315/go-tasker/internal/handler"
	"github.com/kartikey315/go-tasker/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerCategoryRoutes(r *echo.Group, h *handler.CategoryHandler, auth *middleware.AuthMiddleware) {
	// Comment operations
	categories := r.Group("/categories")
	categories.Use(auth.RequireAuth)

	categories.POST("", h.CreateCategory)
	categories.GET("", h.GetCategories)

	// Individual comment operations
	dynamicCategories := categories.Group("/:id")
	dynamicCategories.PATCH("", h.UpdateCategory)
	dynamicCategories.DELETE("", h.DeleteCategory)
}
