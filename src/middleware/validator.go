package middleware

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	ErrorResponse struct {
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		validator *validator.Validate
	}

	ErrorResp struct {
		Message string `json:"errors"`
	}
)

var validate = validator.New()

func (v XValidator) Validate(data interface{}) []string {
	if errs := validate.Struct(data); errs != nil {
		errorMessages := make([]string, 0)
		for _, fe := range errs.(validator.ValidationErrors) {
			fieldName := fe.Field()

			tag := fe.Tag()

			switch tag {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", fieldName))
			case "email":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is not a valid email address", fieldName))
			case "len":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be %s character length",
					fieldName, fe.Param()))
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be greater than %s character length",
					fieldName, fe.Param()))
			case "max":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be less than %s character length", fieldName, fe.Param()))
			case "gte":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be greater than or equal to %s",
					fieldName, fe.Param()))
			case "lte":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be less than or equal to %s", fieldName,
					fe.Param()))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("Validation failed on %s with tag '%s'", fieldName,
					tag))
			}
		}
		return errorMessages
	}
	return nil
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResp{Message: err.Error()})
}
