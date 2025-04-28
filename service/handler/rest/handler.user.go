package rest

import (
	"net/http"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (h *Handler) Register(ctx *gin.Context) {
	var (
		request models.RegisterUserRequest
		errx    serror.SError
	)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		errx = serror.NewFromErrori(http.StatusBadRequest, err)
		errx.AddComments("[handler][Register] while BodyJSONBind")
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

	errx = h.userUsecase.Register(ctx, request)
	if errx != nil {
		handleError(ctx, errx.Code(), errx)
		return
	}

	ctx.JSON(http.StatusCreated, models.ResponseSuccess{
		Message: "User has successfully to registered",
	})
}

func (h *Handler) Login(ctx *gin.Context) {
	var (
		request models.LoginUser
		errx    serror.SError
	)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		errx = serror.NewFromErrori(http.StatusBadRequest, err)
		errx.AddComments("[handler][login] while BodyJSONBind")
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

	res, errx := h.userUsecase.Login(ctx, request)
	if errx != nil {
		handleError(ctx, errx.Code(), errx)
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseSuccess{
		Message: "User has successfully to login",
		Data:    res,
	})
}

func (h *Handler) RefreshToken(ctx *gin.Context) {
	var (
		request models.RefreshTokenRequest
		errx    serror.SError
	)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		errx = serror.NewFromErrori(http.StatusBadRequest, err)
		errx.AddComments("[handler][RefreshToken] while BodyJSONBind")
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

	res, errx := h.authUsecase.RefreshToken(ctx, request.RefreshToken)
	if errx != nil {
		handleError(ctx, errx.Code(), errx)
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseSuccess{
		Message: "Token has been successfully refreshed",
		Data:    res,
	})
}
