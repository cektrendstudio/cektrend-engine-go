package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/logger"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	// Ensure the serror package is correctly imported and available
)

func handleError(ctx *gin.Context, statusCode int, errx serror.SError) (result gin.H) {
	if statusCode == 0 || statusCode == http.StatusInternalServerError {
		logger.Err(errx)
		ctx.JSON(errx.Code(), models.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
			Error:      errx.Error(),
		})
		return
	}

	ctx.JSON(errx.Code(), models.ResponseError{
		StatusCode: errx.Code(),
		Message:    errx.Error(),
	})
	return
}

func handleValidationError(ctx *gin.Context, validationErrors interface{}) (result gin.H) {
	ctx.JSON(http.StatusUnprocessableEntity, models.ResponseError{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    "Validation error",
		Error:      validationErrors,
	})

	return
}

func BuildAndGetValidationMessage(err error) map[string]string {
	validationErrors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		fieldName := strings.ToLower(err.Field()) // Convert field name to lowercase
		var errorMessage string

		switch err.Tag() {
		case "required":
			errorMessage = fmt.Sprintf("%s is required", fieldName)
		case "min":
			errorMessage = fmt.Sprintf("%s must be at least %s characters long", fieldName, err.Param())
		case "max":
			errorMessage = fmt.Sprintf("%s must not exceed %s characters", fieldName, err.Param())
		case "email":
			errorMessage = fmt.Sprintf("%s must be a valid email address", fieldName)
		default:
			errorMessage = fmt.Sprintf("%s failed validation on rule '%s'", fieldName, err.Tag())
		}

		validationErrors[fieldName] = errorMessage
	}

	return validationErrors
}
