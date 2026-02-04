package compressMany

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"picCompressor/internal/http-server/handlers/pictures/compress"
	"picCompressor/internal/lib/api/response"
	"picCompressor/internal/lib/sl"
	"picCompressor/internal/services/PictureWorker"
	"sync"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Items []CompressItem `json:"items"`
}

type CompressItem struct {
	Filepath string `json:"file_path" validate:"required,path"`
	Filename string `json:"filename" validate:"required"`
}

func New(log *slog.Logger, saver compress.PictureSaver, pool *PictureWorker.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pictures.compressMany.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request body"))

			return
		}

		if err != nil {
			log.Error("can not decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("can not decode request body"))

			return
		}

		log.Info("request body decoded!", slog.Any("request", req))

		reply := make(chan PictureWorker.Result, len(req.Items))
		results := make([]PictureWorker.Result, 0, len(req.Items))

		var wg sync.WaitGroup
		wg.Add(len(req.Items))

		for _, item := range req.Items {
			pool.Enqueue(
				PictureWorker.Job{
					Path:  item.Filepath,
					Alias: item.Filename,
					Reply: reply,
					Wg:    &wg,
				},
			)
		}

		go func() {
			wg.Wait()
			close(reply)
		}()

		for res := range reply {
			results = append(results, res)
		}

		ok := make([]PictureWorker.Result, 0, len(results))

		for _, res := range results {
			if err != nil {
				log.Error("compress failed", sl.Err(err))
				continue
			}
			ok = append(ok, res)
		}

		if len(ok) == 0 {
			render.JSON(w, r, response.Error("all images failed"))
			return
		}

		err = saver.AddManyPictures(results)
		if err != nil {
			log.Error("can not add pictures", sl.Err(err))

			render.JSON(w, r, response.Error("can not add pictures"))

			return
		}

		return
	}
}
