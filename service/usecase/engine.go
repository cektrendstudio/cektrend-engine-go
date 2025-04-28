package usecase

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/service"
	"github.com/chromedp/chromedp"

	"github.com/eko/gocache/v2/cache"
)

type EngineUsecase struct {
	userRepo service.UserRepository
	authRepo service.AuthRepository
	cache    *cache.ChainCache
}

func NewEngineUsecase(
	userRepo service.UserRepository,
	authRepo service.AuthRepository,
	cache *cache.ChainCache,
) service.EngineUsecase {
	return &EngineUsecase{
		userRepo: userRepo,
		authRepo: authRepo,
		cache:    cache,
	}
}

func (u *EngineUsecase) WebScreenshot(ctx context.Context, request models.RegisterUserRequest) (errx serror.SError) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	url := "https://openai.com"

	var buf []byte

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.FullScreenshot(&buf, 90),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("screenshot.png", buf, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Screenshot saved as screenshot.png")

	return
}
