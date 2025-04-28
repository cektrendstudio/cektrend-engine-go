package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/uttime"
	"github.com/cektrendstudio/cektrend-engine-go/service"
	"github.com/cektrendstudio/cektrend-engine-go/service/helper"

	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
	"github.com/go-redis/redis/v8"
)

type AuthUsecase struct {
	authRepo service.AuthRepository
	cache    *cache.ChainCache
}

func NewAuthUsecase(
	authRepo service.AuthRepository,
	cache *cache.ChainCache,
) service.AuthUsecase {
	return &AuthUsecase{
		authRepo: authRepo,
		cache:    cache,
	}
}

func (u *AuthUsecase) ValidateAccessToken(ctx context.Context, accessToken string) (userID int64, errx serror.SError) {
	cacheKey := fmt.Sprintf("access_token:%s", accessToken)
	cachedData, err := u.cache.Get(ctx, cacheKey)
	if err != nil && err != redis.Nil {
		errx = serror.NewFromError(err)
		return
	}

	if cachedData != nil {
		var tokenData struct {
			UserID    int64     `json:"user_id"`
			ExpiresAt time.Time `json:"expires_at"`
		}

		dataByte := cachedData.(string)
		err = json.Unmarshal([]byte(dataByte), &tokenData)
		if err != nil {
			errx = serror.NewFromError(err)
			errx.AddCommentf("[usecase][ValidateToken] while get cache")
			return
		}

		if tokenData.ExpiresAt.Before(time.Now()) {
			errx = serror.Newi(http.StatusUnauthorized, "Token invalid")
			return
		}
		return tokenData.UserID, nil
	}

	accessTokenClaims, err := helper.VerifyToken(accessToken, "access")
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][ValidateToken] while verify token")
		return
	}

	userID = int64(accessTokenClaims["user_id"].(float64))
	tokenData, err := u.authRepo.ValidateAccessToken(ctx, models.ValidateAccessTokenRequest{
		UserID:      userID,
		AccessToken: accessToken,
	})
	if tokenData.ExpiresAt.Before(time.Now()) {
		errx = serror.Newi(http.StatusUnauthorized, "Token invalid")
		return
	}

	tokenDataUpdate := map[string]interface{}{
		"user_id":    tokenData.UserID,
		"expires_at": tokenData.ExpiresAt,
	}
	tokenJSON, _ := json.Marshal(tokenDataUpdate)
	expiration := time.Until(uttime.MostParse(uttime.Format(helper.DefaultDateTimeFormat, tokenData.ExpiresAt))).Minutes()
	err = u.cache.Set(ctx, cacheKey, tokenJSON, &store.Options{Expiration: time.Duration(expiration) * time.Minute})
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][ValidateToken] while set cache")
		return
	}

	return userID, nil
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (res models.RefreshTokenResponse, errx serror.SError) {
	refreshTokenClaims, err := helper.VerifyToken(refreshToken, "refresh")
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][RefreshToken] while verify token")
		return
	}

	userId := int64(refreshTokenClaims["user_id"].(float64))
	tokenData, errx := u.authRepo.ValidateRefreshToken(ctx, models.ValidateRefreshTokenRequest{
		UserID:       int64(userId),
		RefreshToken: refreshToken,
	})
	if errx != nil {
		errx.AddCommentf("[usecase][RefreshToken] while ValidateRefreshToken")
		return
	}

	if tokenData.TokenID == 0 {
		errx = serror.Newi(http.StatusUnauthorized, "Token not found or has revoked")
		return
	}

	if tokenData.RefreshTokenExpiresAt.Before(time.Now()) {
		errx = serror.Newi(http.StatusUnauthorized, "Token has expired")
		return
	}

	newAccessToken, newAccessExpiresAt, err := helper.GenerateAccessToken(tokenData.UserID)
	if err != nil {
		errx = serror.NewFromError(err)
		return
	}

	newRefreshToken, newRefreshExpiresAt, err := helper.GenerateRefreshToken(tokenData.UserID)
	if err != nil {
		errx = serror.NewFromError(err)
		return
	}

	token := models.UpdateTokenRequest{
		TokenID:               tokenData.TokenID,
		AccessToken:           newAccessToken,
		AccessTokenExpiresAt:  newAccessExpiresAt,
		RefreshToken:          newRefreshToken,
		RefreshTokenExpiresAt: newRefreshExpiresAt,
	}
	errx = u.authRepo.UpdateToken(ctx, token)
	if errx != nil {
		errx.AddComments("[usecase][Login] Failed to UpdateToken")
		return
	}

	cacheKey := fmt.Sprintf("access_token:%s", newAccessToken)
	tokenDataUpdate := map[string]interface{}{
		"user_id":    tokenData.UserID,
		"expires_at": time.Now().Add(time.Hour),
	}
	tokenJSON, _ := json.Marshal(tokenDataUpdate)
	err = u.cache.Set(ctx, cacheKey, tokenJSON, &store.Options{Expiration: time.Duration(1) * time.Hour})
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[usecase][RefreshToken] while set cache")
		return
	}

	res = models.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    newAccessExpiresAt,
	}

	return
}
