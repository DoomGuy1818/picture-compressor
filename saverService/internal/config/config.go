package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env   string `yaml:"env" env-required:"true"`
	Minio MinIO  `yaml:"minio" env-required:"true"`
	Kafka Kafka  `yaml:"kafka" env-required:"true"`
}

type MinIO struct {
	Endpoint        string `yaml:"endpoint" env-default:"localhost:9090"`
	AccessKeyID     string `yaml:"access_key_id" env-required:"true"`
	SecretAccessKey string `yaml:"secret_access_key" env-required:"true"`
	BucketName      string `yaml:"bucket_name" env-required:"true"`
}

type Kafka struct {
	URL     string `yaml:"url" env-required:"true"`
	Topic   string `yaml:"topic" env-required:"true"`
	GroupID string `yaml:"group_id" env-required:"true"`
}

func MustLoad() *Config {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("CONFIG_PATH env variable not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("file %s does not exist", cfgPath)
	}

	cfg := new(Config)

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	return cfg
}
