package validator

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	validate *validator.Validate
}

func New() *CustomValidator {
	v := validator.New()

	return &CustomValidator{
		validate: v,
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validate.Struct(i)
}
