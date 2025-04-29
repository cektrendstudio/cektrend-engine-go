package rest

import (
	"net/http"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (h *Handler) WebScreenshot(ctx *gin.Context) {
	var (
		request models.WebScreenshotRequest
		errx    serror.SError
	)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		errx = serror.NewFromErrori(http.StatusBadRequest, err)
		errx.AddComments("[handler][WebScreenshot] while BodyJSONBind")
		handleError(ctx, errx.Code(), errx)
		return
	}

	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		validationMessages := BuildAndGetValidationMessage(err)
		handleValidationError(ctx, validationMessages)

		return
	}

	res, errx := h.engineUsecase.WebScreenshot(request)
	if errx != nil {
		handleError(ctx, errx.Code(), errx)
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseSuccess{
		Message: "WEB Screenshot success",
		Data:    res,
	})
}
