package service

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
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

type OSSRepository interface {
	UploadFileFromURL(ctx context.Context, name, fileURL string, isPublic, allowFail bool) (newURL string, errx serror.SError)
	UploadFileFromBytes(ctx context.Context, name string, file []byte) (newURL string, serr serror.SError)
	GetFileExtension(ctx context.Context, inputUrl string) (string, error)
}

type GoogleDriveRepository interface {
	GetFile(ctx context.Context, fileID string) (file *http.Response, serr serror.SError)
	GetDriveFileExtension(ctx context.Context, fileURL string) (fileFormat string, serr serror.SError)
}

type StorageRepository interface {
	UploadFile(ctx context.Context, fileheader *multipart.FileHeader, object string, fileExt string, expire time.Time) (string, error)
	UploadFileWithContentType(ctx context.Context, file multipart.File, contentType string, object string, expire time.Time) (string, error)
	UploadFileWithBytes(ctx context.Context, buffer bytes.Buffer, contentType string, object string, expire time.Time) (string, error)
	GetFilePathSignedURL(object string) (string, error)
	GetFilePathSignedURLCustom(object string, expire time.Time) (string, error)
	GetFilePathSignedURLExcel(object string) (string, error)
	DownloadFile(ctx context.Context, object string, destFileName string) error
	DeleteFile(ctx context.Context, object string) error
}

type S3Repository interface {
	UploadFile(ctx context.Context, file []byte, fileName string, contentType string) (string, error)
}
