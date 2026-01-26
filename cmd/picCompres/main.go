package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"picCompressor/internal/config"
	"picCompressor/internal/http-server/handlers/pictures/compress"
	"picCompressor/internal/lib/compressor"
	"picCompressor/internal/lib/compressor/baselib"
	"picCompressor/internal/lib/sl"
	vault "picCompressor/internal/object"
	"picCompressor/internal/object/s3"
	eventSender "picCompressor/internal/services/event-sender"
	kafkaBroker "picCompressor/internal/services/queue/kafka-broker"
	"picCompressor/internal/storage/pg"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()

	log := initLogger(cfg.Env)

	log.Info("start picture compressor", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled!")

	minio, err := s3.New(cfg.Minio.Endpoint, cfg.Minio.AccessKeyID, cfg.Minio.SecretAccessKey, cfg.Minio.BucketName)
	if err != nil {
		log.Error("failed to init S3 client", sl.Err(err))
	}

	if err = minio.CreateBucketWithCheck(ctx, cfg.Minio.BucketName); err != nil {
		log.Error("failed to init S3 bucket", sl.Err(err))
		os.Exit(1)
	} else {
		log.Info("bucket has been created!", slog.String("bucket", cfg.Minio.BucketName))
	}

	str, err := pg.New(cfg.StorageURL)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	kafka := kafkaBroker.New(cfg.Kafka.URL, cfg.Kafka.Topic)

	compr := baselib.New(10, -1)

	router := initRouter(log, str, compr, minio)

	log.Info("starting server", slog.String("address", cfg.HTTPServer.URL))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.URL,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	sender := eventSender.New(str, log, kafka)
	sender.RunProcessingEvents(context.Background(), 5*time.Second)

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")
}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func initRouter(
	log *slog.Logger,
	saver compress.PictureSaver,
	compressor compressor.Compressor,
	vault vault.SaverInVault,
) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route(
		"/compress", func(r chi.Router) {
			r.Post("/", compress.New(log, saver, compressor, vault))
		},
	)

	return router
}
