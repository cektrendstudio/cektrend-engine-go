package config

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/service/handler/rest"
	"github.com/cektrendstudio/cektrend-engine-go/service/repository/postgres"
	"github.com/cektrendstudio/cektrend-engine-go/service/repository/storage"
	"github.com/cektrendstudio/cektrend-engine-go/service/usecase"
)

func (cfg *Config) InitService() (errx serror.SError) {
	// messageRepo := postgres.NewMessageRepository(cfg.RabbitMQChannel)

	authRepo := postgres.NewAuthRepository(cfg.DB)
	authUsecase := usecase.NewAuthUsecase(authRepo, cfg.Cache)

	userRepo := postgres.NewUserRepository(cfg.DB)
	userUsecase := usecase.NewUserUsecase(userRepo, authRepo, cfg.Cache)

	s3Repo := storage.NewS3Repository(s3.New(cfg.AWSConfig.S3Session), cfg.AWSConfig)

	engineUsecase := usecase.NewEngineUsecase(userRepo, authRepo, cfg.Cache, s3Repo)
	route := rest.CreateHandler(
		userUsecase,
		authUsecase,
		engineUsecase,
	)

	cfg.Server = route

	// broker.NewRabbitMQHandler(true, cfg.RabbitMQChannel, transactionUsecase)

	return nil
}
