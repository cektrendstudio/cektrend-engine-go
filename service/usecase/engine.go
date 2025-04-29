package usecase

import (
	"context"
	"net/http"
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
	s3Repo   service.S3Repository
}

func NewEngineUsecase(
	userRepo service.UserRepository,
	authRepo service.AuthRepository,
	cache *cache.ChainCache,
	s3Repo service.S3Repository,
) service.EngineUsecase {
	return &EngineUsecase{
		userRepo: userRepo,
		authRepo: authRepo,
		cache:    cache,
		s3Repo:   s3Repo,
	}
}

func (u *EngineUsecase) WebScreenshot(request models.WebScreenshotRequest) (res models.WebScreenshotResponse, errx serror.SError) {
	apiKeys := []string{
		"b1a2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
	}

	if request.Key == "" {
		errx = serror.Newi(http.StatusBadRequest, "API Key Not Found")
		return
	}

	valid := false
	for _, apiKey := range apiKeys {
		if request.Key == apiKey {
			valid = true
			break
		}
	}

	if !valid {
		errx = serror.Newi(http.StatusBadRequest, "Invalid API Key")
		return
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctxTimeout, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	var file []byte
	if err := chromedp.Run(ctxTimeout,
		chromedp.Navigate(request.URL),
		chromedp.Sleep(2*time.Second),
		chromedp.CaptureScreenshot(&file),
	); err != nil {
		errx = serror.NewFromErrorc(err, "failed to capture screenshot")
		return
	}

	newFileName := "cektrend-engine-storage/web-ss-" + time.Now().Format("20060102150405") + ".png"
	newURL, err := u.s3Repo.UploadFile(ctx, file, newFileName, "image/png")
	if err != nil {
		errx = serror.NewFromErrori(http.StatusInternalServerError, err)
		errx.AddComments("[usecase][WebScreenshot] while UploadFile")
		return
	}

	res = models.WebScreenshotResponse{
		URL: newURL,
	}

	return
}
