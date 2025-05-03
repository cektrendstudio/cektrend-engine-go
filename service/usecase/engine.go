package usecase

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/service"
	"github.com/chromedp/chromedp"
	"github.com/nfnt/resize"
	"github.com/xuri/excelize/v2"

	"github.com/disintegration/imaging"
	"github.com/eko/gocache/v2/cache"
)

type EngineUsecase struct {
	userRepo   service.UserRepository
	authRepo   service.AuthRepository
	cache      *cache.ChainCache
	s3Repo     service.S3Repository
	engineRepo service.EngineRepository
}

func NewEngineUsecase(
	userRepo service.UserRepository,
	authRepo service.AuthRepository,
	cache *cache.ChainCache,
	s3Repo service.S3Repository,
	engineRepo service.EngineRepository,
) service.EngineUsecase {
	return &EngineUsecase{
		userRepo:   userRepo,
		authRepo:   authRepo,
		cache:      cache,
		s3Repo:     s3Repo,
		engineRepo: engineRepo,
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

func (u *EngineUsecase) PhishingWebReportFromExcel(ctx context.Context) error {
	var version int8 = 2
	f, err := excelize.OpenFile("./tmp/phishing-web-report.xlsx")
	if err != nil {
		return err
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return err
	}

	var urls []string
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) > 0 {
			urls = append(urls, row[0])
		}
	}

	if len(urls) == 0 {
		return fmt.Errorf("no URLs found in the Excel file")
	}

	res, errx := u.engineRepo.GetPhishingReportByURLs(ctx, urls, version)
	if errx != nil {
		return errx
	}

	// remove url from the slice if it already exists in the database
	for _, report := range res {
		for i, url := range urls {
			if report.SiteURL == url {
				urls = append(urls[:i], urls[i+1:]...)
				break
			}
		}
	}

	const workerCount = 1
	var wg sync.WaitGroup
	jobs := make(chan string)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				fmt.Println("Processing URL:", url)
				u.handlePhishingWebScreenshot(ctx, url, version)
				fmt.Println("Finished processing URL:", url)
			}
		}()
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	wg.Wait()
	return nil
}

func (u *EngineUsecase) DownloadImageFromExcel(ctx context.Context) error {
	f, err := excelize.OpenFile("./tmp/phishing-web-report.xlsx")
	if err != nil {
		return err
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return err
	}

	var siteData []models.CreatePhishingWebReportRequest
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) > 0 {
			siteData = append(siteData, models.CreatePhishingWebReportRequest{
				SiteURL:  row[0],
				ImageURL: row[1],
			})
		}
	}

	if len(siteData) == 0 {
		return fmt.Errorf("no URLs found in the Excel file")
	}

	const workerCount = 10
	var wg sync.WaitGroup
	jobs := make(chan models.CreatePhishingWebReportRequest)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for data := range jobs {
				fmt.Println("Processing Download Image URL:", data.ImageURL)
				parsedURL, err := url.Parse(data.SiteURL)
				if err != nil {
					fmt.Println("Error parsing URL:", err)
					continue
				}
				filename := parsedURL.Hostname() + ".png"
				outputPath := filepath.Join("./tmp/download_images/", filename)
				u.downloadImage(data.ImageURL, outputPath)
				fmt.Println("Finished processing download URL:", data.ImageURL)
			}
		}()
	}

	for _, site := range siteData {
		jobs <- site
	}
	close(jobs)

	wg.Wait()
	return nil
}

func (u *EngineUsecase) downloadImage(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("request gagal: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return fmt.Errorf("gagal decode gambar: %w", err)
	}

	resized := imaging.Resize(img, 500, 0, imaging.Lanczos)

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("gagal membuat file: %w", err)
	}
	defer out.Close()

	err = png.Encode(out, resized)
	if err != nil {
		return err
	}

	return err
}

func (u *EngineUsecase) handlePhishingWebScreenshot(ctx context.Context, url string, version int8) {
	newCtx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctxTimeout, cancelTimeout := context.WithTimeout(newCtx, 30*time.Second)
	defer cancelTimeout()

	var file []byte
	if err := chromedp.Run(ctxTimeout,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(url),
		chromedp.Evaluate(`document.body.insertAdjacentHTML('afterbegin', '<div style="position:fixed;top:0;left:0;width:100%;padding:8px;background:#eee;font-family:sans-serif;z-index:9999;border-bottom:1px solid #ccc;">`+url+`</div>')`, nil),
		chromedp.Sleep(2*time.Second),
		chromedp.CaptureScreenshot(&file),
	); err != nil {
		fmt.Printf("Error capturing screenshot: %s, err: %v", url, err)
		return
	}

	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	resizedImg := resize.Resize(1024, 0, img, resize.Lanczos3)

	var buf bytes.Buffer
	encoder := png.Encoder{
		CompressionLevel: png.DefaultCompression,
	}
	if err := encoder.Encode(&buf, resizedImg); err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}

	compressed := buf.Bytes()

	newFileName := "cektrend-engine-storage/web-ss-" + time.Now().Format("20060102150405") + ".png"
	imageUrl, err := u.s3Repo.UploadFile(ctx, compressed, newFileName, "image/png")
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}

	_, err = u.engineRepo.InsertPhisingWebReport(ctx, models.CreatePhishingWebReportRequest{
		SiteURL:  url,
		ImageURL: imageUrl,
		Version:  version,
	})
	if err != nil {
		fmt.Println("Error inserting phishing web report:", err)
		return
	}
}
