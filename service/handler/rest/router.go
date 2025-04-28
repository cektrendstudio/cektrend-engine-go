package rest

import (
	"os"
	"time"

	"github.com/cektrendstudio/cektrend-engine-go/service"
	"github.com/cektrendstudio/cektrend-engine-go/service/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userUsecase service.UserUsecase
	authUsecase service.AuthUsecase
}

func CreateHandler(
	userUsecase service.UserUsecase,
	authUsecase service.AuthUsecase,
) *gin.Engine {
	obj := Handler{
		userUsecase: userUsecase,
		authUsecase: authUsecase,
	}

	r := gin.Default()

	gin.SetMode(gin.DebugMode)
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(middleware.LoggingMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.RateLimitMiddleware(10, 20))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET, POST, PUT, PATCH, DELETE, OPTIONS"},
		AllowHeaders:     []string{"Access-Control-Allow-Headers", "Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	publicRouter := r.Group("/v1")
	publicRouter.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "welcome, to golang project starter",
		})
	})
	publicRouter.POST("/register", obj.Register)
	publicRouter.POST("/login", obj.Login)
	publicRouter.POST("/refresh-token", obj.RefreshToken)

	authRouter := publicRouter.Group("/")

	authMiddleware := &middleware.AuthMiddleware{
		AuthUsecase: authUsecase,
	}
	authRouter.Use(authMiddleware.Authenticate())
	{
		authRouter.POST("/web-ss", obj.WebScreenshot)
	}

	return r
}
