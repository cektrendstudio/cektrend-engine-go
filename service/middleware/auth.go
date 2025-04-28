package middleware

import (
	"context"
	"net/http"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/service"

	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	AuthUsecase service.AuthUsecase
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ResponseError{
				Message: "Token Not Found",
				Error:   "Unauthorized",
			})

			return
		}

		bearer := strings.HasPrefix(token, "Bearer")
		if !bearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ResponseError{
				Message: "Bearer Not Found",
				Error:   "Unauthorized",
			})

			return
		}

		tokenStr := strings.Split(token, "Bearer ")[1]
		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ResponseError{
				Message: "Token STR Not Found",
				Error:   "Unauthorized",
			})

			return
		}

		userID, err := m.AuthUsecase.ValidateAccessToken(context.Background(), tokenStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ResponseError{
				Message: err.Error(),
				Error:   "Unauthorized",
			})

			log.Errorln("ERROR:", err)
			return
		}

		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
