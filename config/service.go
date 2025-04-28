package config

import (
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/service/handler/rest"
	"github.com/cektrendstudio/cektrend-engine-go/service/repository/postgres"
	"github.com/cektrendstudio/cektrend-engine-go/service/usecase"
)

func (cfg *Config) InitService() (errx serror.SError) {
	// messageRepo := postgres.NewMessageRepository(cfg.RabbitMQChannel)

	authRepo := postgres.NewAuthRepository(cfg.DB)
	authUsecase := usecase.NewAuthUsecase(authRepo, cfg.Cache)

	userRepo := postgres.NewUserRepository(cfg.DB)
	userUsecase := usecase.NewUserUsecase(userRepo, authRepo, cfg.Cache)
	route := rest.CreateHandler(
		userUsecase,
		authUsecase,
	)

	cfg.Server = route

	// broker.NewRabbitMQHandler(true, cfg.RabbitMQChannel, transactionUsecase)

	return nil
}
