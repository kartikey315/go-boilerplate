package service

import (
	"fmt"

	"github.com/kartikey315/go-tasker/internal/lib/aws"
	"github.com/kartikey315/go-tasker/internal/lib/job"
	"github.com/kartikey315/go-tasker/internal/repository"
	"github.com/kartikey315/go-tasker/internal/server"
)

type Services struct {
	Auth     *AuthService
	Job      *job.JobService
	Todo     *TodoService
	Category *CategoryService
	Comment  *CommentService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	s.Job.SetAuthService(authService)

	awsClient, err := aws.NewAWS(s)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %w", err)
	}

	return &Services{
		Job:      s.Job,
		Auth:     authService,
		Todo:     NewTodoService(s, repos.Todo, repos.Category, awsClient),
		Category: NewCategoryRepository(s, repos.Category),
		Comment:  NewCommentService(s, repos.Todo, repos.Comment),
	}, nil
}
