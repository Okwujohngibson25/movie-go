package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

func NewValidator() *validator.Validate {
	v := validator.New()

	err := v.RegisterValidation("rating", func(fl validator.FieldLevel) bool {
		rating, ok := fl.Field().Interface().(decimal.Decimal)
		if !ok {
			return false
		}

		return rating.GreaterThanOrEqual(decimal.NewFromInt(1)) &&
			rating.LessThanOrEqual(decimal.NewFromInt(5))
	})
	if err != nil {
		panic(err)
	}

	return v
}
