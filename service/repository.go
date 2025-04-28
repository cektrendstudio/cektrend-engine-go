package service

import (
	"context"
	"time"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Register(ctx context.Context, req models.RegisterUserRequest) (userId int64, errx serror.SError)
	GetUserByID(ctx context.Context, userID int64) (res models.User, errx serror.SError)
	GetUserByEmail(ctx context.Context, email string) (res models.User, errx serror.SError)
}

type MessageRepository interface {
	Publish(queue string, data interface{}) (errx serror.SError)
}

type AuthRepository interface {
	CreateToken(ctx context.Context, req models.CreateTokenRequest) (tokenID int64, errx serror.SError)
	RevokedInactiveToken(ctx context.Context, req models.RevokedInactiveTokenRequest) (errx serror.SError)
	ValidateAccessToken(ctx context.Context, req models.ValidateAccessTokenRequest) (res models.ValidateAccessTokenResponse, errx serror.SError)
	ValidateRefreshToken(ctx context.Context, req models.ValidateRefreshTokenRequest) (res models.ValidateRefreshTokenResponse, errx serror.SError)
	UpdateToken(ctx context.Context, req models.UpdateTokenRequest) (errx serror.SError)
}

type CacheRepository interface {
	Set(key string, value interface{}, duration time.Duration) error
	Get(key string) (data string)
}

type SQLRepository interface {
	BeginTxx() (tx *sqlx.Tx, err error)
	Commit(tx *sqlx.Tx) (err error)
	Rollback(tx *sqlx.Tx) (err error)
}
