package service

import (
	"context"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
)

type AuthUsecase interface {
	ValidateAccessToken(ctx context.Context, accessToken string) (userID int64, errx serror.SError)
	RefreshToken(ctx context.Context, refreshToken string) (res models.RefreshTokenResponse, errx serror.SError)
}

type UserUsecase interface {
	Register(ctx context.Context, request models.RegisterUserRequest) (errx serror.SError)
	Login(ctx context.Context, request models.LoginUser) (res models.LoginResponse, errx serror.SError)
}

type EngineUsecase interface {
	WebScreenshot(request models.WebScreenshotRequest) (res models.WebScreenshotResponse, errx serror.SError)
	PhishingWebReportFromExcel(ctx context.Context) error
	DownloadImageFromExcel(ctx context.Context) error
}
