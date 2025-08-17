package validator

import "github.com/go-playground/validator/v10"

// Setup provides a validator instance for DI.
func Setup() (*validator.Validate, error) {
	v := validator.New()
	return v, nil
}
