package service

import (
	"github.com/google/uuid"
	"github.com/kartikey315/go-tasker/internal/middleware"
	"github.com/kartikey315/go-tasker/internal/model"
	"github.com/kartikey315/go-tasker/internal/model/category"
	"github.com/kartikey315/go-tasker/internal/repository"
	"github.com/kartikey315/go-tasker/internal/server"
	"github.com/labstack/echo/v4"
)

type CategoryService struct {
	server       *server.Server
	categoryRepo *repository.CategoryRepository
}

func NewCategoryRepository(server *server.Server, categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		server:       server,
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) CreateCategory(ctx echo.Context, userID string, payload *category.CreateCategoryPayload) (*category.Category, error) {
	logger := middleware.GetLogger(ctx)

	categoryItem, err := s.categoryRepo.CreateCategory(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create category")
		return nil, err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "category_created").
		Str("todo_id", categoryItem.ID.String()).
		Str("name", categoryItem.Name).
		Str("color", categoryItem.Color).
		Str("description", *categoryItem.Description).
		Msg("Category Created Successfully")

	return categoryItem, nil

}

func (s *CategoryService) GetCategories(ctx echo.Context, userID string, query *category.GetCategoriesQuery) (*model.PaginatedResponse[category.Category], error) {
	logger := middleware.GetLogger(ctx)

	categories, err := s.categoryRepo.GetCategories(ctx.Request().Context(), userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get categories")
		return nil, err
	}

	return categories, nil
}

func (s *CategoryService) UpdateCategory(ctx echo.Context, userID string, categoryID uuid.UUID, payload *category.UpdateCategoryPayload) (*category.Category, error) {
	logger := middleware.GetLogger(ctx)

	// Validate category exists and belongs to user
	_, err := s.categoryRepo.GetCategoryByID(ctx.Request().Context(), userID, categoryID)
	if err != nil {
		logger.Error().Err(err).Msg("category validation failed")
		return nil, err
	}

	updatedCategory, err := s.categoryRepo.UpdateCategory(ctx.Request().Context(), userID, categoryID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update category")
		return nil, err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "category_updated").
		Str("category_id", updatedCategory.ID.String()).
		Str("Name", updatedCategory.Name).
		Str("Color", updatedCategory.Color).
		Msg("Category Updated Successfully")

	return updatedCategory, nil
}

func (s *CategoryService) DeleteCategory(ctx echo.Context, userID string, categoryID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	err := s.categoryRepo.Deletecategory(ctx.Request().Context(), userID, categoryID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete category")
		return err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "category_updated").
		Str("category_id", categoryID.String()).
		Msg("Category Updated Successfully")

	return nil

}
