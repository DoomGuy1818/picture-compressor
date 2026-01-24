package compress

import (
	"context"
	"log/slog"
	"net/http"
	"picCompressor/internal/object/s3"
)

type Request struct {
	FilePath string `json:"file_path" validate:"required, path"`
	Alias    string `json:"alias" validate:"omitempty"`
}

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error" validate:"error,omitempty"`
	Alias  string `json:"alias" validate:"error,omitempty"`
}

type PictureSaver interface {
	AddPicture(ctx context.Context, path string, vault s3.Minio) error
}

func New(log *slog.Logger, saver PictureSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: написать хендлер для добавления сжатой картинки в s3
	}
}
