package rest

import (
	"net/http"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"

	"github.com/gin-gonic/gin"
)

func (h *Handler) WebScreenshot(ctx *gin.Context) {
	var (
		request models.WebScreenshotRequest
		errx    serror.SError
	)

	request.Key = ctx.Query("key")
	request.URL = ctx.Query("url")
	res, errx := h.engineUsecase.WebScreenshot(request)
	if errx != nil {
		handleError(ctx, errx.Code(), errx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url": res.URL,
	})
}

func (h *Handler) PhishingReport(ctx *gin.Context) {
	var (
		errx serror.SError
	)

	err := h.engineUsecase.PhishingWebReportFromExcel(ctx)
	if err != nil {
		errx = serror.NewFromError(err)
		handleError(ctx, errx.Code(), errx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *Handler) DownloadImageFromExcel(ctx *gin.Context) {
	var (
		errx serror.SError
	)

	err := h.engineUsecase.DownloadImageFromExcel(ctx)
	if err != nil {
		errx = serror.NewFromError(err)
		handleError(ctx, errx.Code(), errx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
