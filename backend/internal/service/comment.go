package service

import (
	"github.com/google/uuid"
	"github.com/kartikey315/go-tasker/internal/middleware"
	"github.com/kartikey315/go-tasker/internal/model/comment"
	"github.com/kartikey315/go-tasker/internal/repository"
	"github.com/kartikey315/go-tasker/internal/server"
	"github.com/labstack/echo/v4"
)

type CommentService struct {
	server      *server.Server
	commentRepo *repository.CommentRepository
	todoRepo    *repository.TodoRepository
}

func NewCommentService(server *server.Server, todoRepo *repository.TodoRepository, commentRepo *repository.CommentRepository) *CommentService {
	return &CommentService{
		server:      server,
		commentRepo: commentRepo,
		todoRepo:    todoRepo,
	}
}

func (s *CommentService) AddComment(ctx echo.Context, userID string, todoID uuid.UUID, payload *comment.AddCommentPayload) (*comment.Comment, error) {
	logger := middleware.GetLogger(ctx)

	// Validate todo exists and belongs to user
	_, err := s.todoRepo.CheckTodoExists(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("todo validation failed")
		return nil, err
	}

	commentItem, err := s.commentRepo.AddComment(ctx.Request().Context(), userID, todoID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to add comment")
		return nil, err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "comment_added").
		Str("todo_id", todoID.String()).
		Str("comment_id", commentItem.ID.String()).
		Str("content", commentItem.Content).
		Msg("Comment Added Successfully")

	return commentItem, nil

}

func (s *CommentService) GetCommentsByTodoID(ctx echo.Context, userID string, todoID uuid.UUID) ([]comment.Comment, error) {
	logger := middleware.GetLogger(ctx)

	// Validate todo exists and belongs to user
	_, err := s.todoRepo.CheckTodoExists(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("todo validation failed")
		return nil, err
	}

	comments, err := s.commentRepo.GetCommentsByTodoID(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comments by Id")
		return nil, err
	}

	return comments, nil
}

func (s *CommentService) UpdateComment(ctx echo.Context, userID string, commentID uuid.UUID, content string) (*comment.Comment, error) {
	logger := middleware.GetLogger(ctx)

	// Validate comment exists and belongs to user
	_, err := s.commentRepo.GetCommentByID(ctx.Request().Context(), userID, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("comment validation failed")
		return nil, err
	}

	updatedComment, err := s.commentRepo.UpdateComment(ctx.Request().Context(), userID, commentID, content)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update comment")
		return nil, err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "comment_updated").
		Str("comment_id", updatedComment.ID.String()).
		Str("content", updatedComment.Content).
		Msg("Comment Updated Successfully")

	return updatedComment, nil
}

func (s *CommentService) DeleteComment(ctx echo.Context, userID string, commentID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	// Validate comment exists and belongs to user
	_, err := s.commentRepo.GetCommentByID(ctx.Request().Context(), userID, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("comment validation failed")
		return err
	}

	err = s.commentRepo.DeleteComment(ctx.Request().Context(), userID, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete comment")
		return err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "comment_deleted").
		Str("comment_id", commentID.String()).
		Msg("Comment Deleted Successfully")

	return nil
}
