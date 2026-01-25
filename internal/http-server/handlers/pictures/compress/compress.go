package compress

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"picCompressor/internal/lib/compressor"
	"picCompressor/internal/lib/sl"
	"picCompressor/internal/object"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type Request struct {
	FilePath string `json:"file_path"`
	Filename string `json:"filename"`
}

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error" validate:"error,omitempty"`
}

type PictureSaver interface {
	AddPicture(id uuid.UUID, path string) error
}

func New(
	log *slog.Logger,
	saver PictureSaver,
	compressor compressor.Compressor,
	vault object.SaverInVault,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.Decode(r, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, Response{Status: http.StatusBadRequest, Error: "invalid request"})

			return
		}

		outFilePath, err := compressor.Compress(req.FilePath, req.Filename)
		if err != nil {
			render.JSON(w, r, Response{Status: http.StatusInternalServerError, Error: "failed to compress file"})
			return
		}

		id, err := uuid.NewV7()
		if err != nil {
			render.JSON(w, r, Response{Status: http.StatusInternalServerError, Error: "failed to generate id"})
			return
		}

		err = saver.AddPicture(id, outFilePath)
		if err != nil {
			render.JSON(w, r, Response{Status: http.StatusInternalServerError, Error: "failed to save file"})
			return
		}

		err = vault.PutObject(r.Context(), outFilePath)
		if err != nil {
			render.JSON(w, r, Response{Status: http.StatusInternalServerError, Error: "failed to save file in vault"})
			return
		}
	}
}
