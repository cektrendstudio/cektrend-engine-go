package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/service"
	"github.com/cektrendstudio/cektrend-engine-go/service/helper"

	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
)

type UserUsecase struct {
	userRepo service.UserRepository
	authRepo service.AuthRepository
	cache    *cache.ChainCache
}

func NewUserUsecase(
	userRepo service.UserRepository,
	authRepo service.AuthRepository,
	cache *cache.ChainCache,
) service.UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
		authRepo: authRepo,
		cache:    cache,
	}
}

func (u *UserUsecase) Register(ctx context.Context, request models.RegisterUserRequest) (errx serror.SError) {
	userArgs := models.RegisterUserRequest{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	userCheck, err := u.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][Register] Failed to get user by email, [email: %s]", request.Email)
		return
	}

	if userCheck.UserID != 0 {
		errx = serror.Newi(http.StatusBadRequest, "Email already registered")
		return
	}

	_, err = u.userRepo.Register(ctx, userArgs)
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][Register] Failed to register user, [email: %s]", request.Email)
		return
	}

	return
}

func (u *UserUsecase) Login(ctx context.Context, request models.LoginUser) (res models.LoginResponse, errx serror.SError) {
	userDB, errx := u.userRepo.GetUserByEmail(ctx, request.Email)
	if errx != nil {
		errx.AddCommentf("[usecase][Login] Failed to get user by email, [email: %s]", request.Email)
		return
	}

	if userDB.UserID == 0 {
		errx = serror.Newi(http.StatusNotFound, "User not found")
		return
	}

	accountMatch := helper.ComparePassword([]byte(userDB.Password), []byte(request.Password))
	if !accountMatch {
		errx = serror.Newi(http.StatusBadRequest, "Password does not match")
		return
	}

	accessToken, accessExpiresAt, err := helper.GenerateAccessToken(userDB.UserID)
	if err != nil {
		errx = serror.NewFromError(err)
		return
	}

	refreshToken, refreshExpiresAt, err := helper.GenerateRefreshToken(userDB.UserID)
	if err != nil {
		errx = serror.NewFromError(err)
		return
	}

	token := models.CreateTokenRequest{
		UserID:                userDB.UserID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}
	tokenID, errx := u.authRepo.CreateToken(ctx, token)
	if errx != nil {
		errx.AddCommentf("[usecase][Login] Failed to create token, [email: %s]", request.Email)
		return
	}

	errx = u.authRepo.RevokedInactiveToken(ctx, models.RevokedInactiveTokenRequest{
		UserID:  userDB.UserID,
		TokenID: tokenID,
	})
	if errx != nil {
		errx.AddCommentf("[usecase][Login] Failed while RevokedInactiveToken, [email: %s]", request.Email)
		return
	}

	// Pakai redis supaya tidak perlu query ke db lagi
	key := fmt.Sprintf("access_token:%s", accessToken)
	tokenData := map[string]interface{}{
		"user_id":    userDB.UserID,
		"expires_at": accessExpiresAt,
	}

	tokenJSON, _ := json.Marshal(tokenData)
	err = u.cache.Set(ctx, key, tokenJSON, &store.Options{Expiration: time.Duration(1) * time.Hour})
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][Login] while set cache")
		return
	}

	res = models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt,
		User:         userDB,
	}

	return
}
