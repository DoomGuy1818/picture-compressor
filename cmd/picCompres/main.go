package main

import (
	"context"
	"log/slog"
	"os"
	"picCompressor/internal/config"
	"picCompressor/internal/lib/sl"
	"picCompressor/internal/object/s3"
	"picCompressor/internal/storage/pg"

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

	_, err = pg.New(cfg.StorageURL)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	initRouter()

	//TODO: run server
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

func initRouter() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
}
