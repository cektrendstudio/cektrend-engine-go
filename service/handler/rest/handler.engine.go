package rest

import (
	"io"
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

	imageResp, err := http.Get(res.URL)
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, serror.NewFromError(err))
		return
	}
	defer imageResp.Body.Close()

	ctx.Header("Content-Type", "image/png")
	ctx.Status(http.StatusOK)

	_, err = io.Copy(ctx.Writer, imageResp.Body)
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, serror.NewFromError(err))
		return
	}
}
