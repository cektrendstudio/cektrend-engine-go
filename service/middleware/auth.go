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

func (m *AuthMiddleware) StrictAuthenticate(allowedIPs []string) gin.HandlerFunc {
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

		clientIP := ctx.ClientIP()
		allowed := false
		for _, ip := range allowedIPs {
			if clientIP == ip {
				allowed = true
				break
			}
		}

		if !allowed {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden - unauthorized IP address"})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", userID)
		ctx.Next()
	}
}

func (m *AuthMiddleware) OpenAPIAuthenticate(apiKeys []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.Request.Header.Get("X-API-Key")
		if key == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ResponseError{
				Message: "API Key Not Found",
				Error:   "Unauthorized",
			})
			return
		}

		valid := false
		for _, apiKey := range apiKeys {
			if key == apiKey {
				valid = true
				break
			}
		}

		if !valid {
			ctx.AbortWithStatusJSON(http.StatusForbidden, models.ResponseError{
				Message: "Invalid API Key",
				Error:   "Forbidden",
			})
			return
		}

		ctx.Next()
	}
}
