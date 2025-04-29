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
	userUsecase   service.UserUsecase
	authUsecase   service.AuthUsecase
	engineUsecase service.EngineUsecase
}

func CreateHandler(
	userUsecase service.UserUsecase,
	authUsecase service.AuthUsecase,
	engineUsecase service.EngineUsecase,
) *gin.Engine {
	obj := Handler{
		userUsecase:   userUsecase,
		authUsecase:   authUsecase,
		engineUsecase: engineUsecase,
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
			"message": "welcome, to cektrend engine",
		})
	})
	publicRouter.POST("/register", obj.Register)
	publicRouter.POST("/login", obj.Login)
	publicRouter.POST("/refresh-token", obj.RefreshToken)

	strictAuth := publicRouter.Group("/")
	openApiAuth := publicRouter.Group("/open-api")

	authMiddleware := &middleware.AuthMiddleware{
		AuthUsecase: authUsecase,
	}

	strictAuth.Use(authMiddleware.StrictAuthenticate([]string{"127.0.0.1"}))
	apiKeys := []string{
		"b1a2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
	}
	openApiAuth.Use(authMiddleware.OpenAPIAuthenticate(apiKeys))

	{
		// strictAuth.POST("/web-ss", obj.WebScreenshot)
	}

	{
		openApiAuth.POST("/web-ss", obj.WebScreenshot)
	}

	return r
}
