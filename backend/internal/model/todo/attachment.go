package todo

import (
	"github.com/google/uuid"
	"github.com/kartikey315/go-tasker/internal/model"
)

type TodoAttachment struct {
	model.Base
	TodoID      uuid.UUID `json:"todoId" db:"todo_id"`
	Name        string    `json:"name" db:"name"`
	UploadedBy  string    `json:"uploadedBy" db:"uploaded_by"`
	DownloadKey string    `json:"downloadKey" db:"download_key"`
	Filesize    *int64    `json:"fileSize" db:"file_size"`
	MimeType    *string   `json:"mimeType" db:"mime_type"`
}
