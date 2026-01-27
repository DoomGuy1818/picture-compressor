package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"saverService/internal/config"
	"saverService/internal/lib/sl"
	"saverService/internal/object/s3"
	eventSender "saverService/internal/services/event-sender"
	kafkaBroker "saverService/internal/services/queue/kafka-broker"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := initLogger(cfg.Env)

	minio, err := s3.New(cfg.Minio.Endpoint, cfg.Minio.AccessKeyID, cfg.Minio.SecretAccessKey, cfg.Minio.BucketName)
	if err != nil {
		log.Error("failed to init S3 client", sl.Err(err))
	}

	if err = minio.CreateBucketWithCheck(context.Background(), cfg.Minio.BucketName); err != nil {
		log.Error("failed to init S3 bucket", sl.Err(err))
		os.Exit(1)
	} else {
		log.Info("bucket has been created!", slog.String("bucket", cfg.Minio.BucketName))
	}

	kafka := kafkaBroker.New(cfg.Kafka.URL, cfg.Kafka.Topic, cfg.Kafka.GroupID)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	consumer := eventSender.New(log, kafka, minio)
	consumer.RunProcessingEvents(context.Background(), 5*time.Second)

	<-done
	log.Info("stopping server")
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
