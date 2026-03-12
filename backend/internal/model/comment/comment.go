package comment

import (
	"github.com/google/uuid"
	"github.com/kartikey315/go-tasker/internal/model"
)

type Comment struct {
	model.Base
	TodoId  uuid.UUID `json:"todoId" db:"todo_id"`
	UserId  string    `json:"userId" db:"user_id"`
	Content string    `json:"content" db:"content"`
}
