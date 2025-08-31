package models

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("breedvalidator", BreedValidator)
	}
}

func MessageForValidationTag(tag string) string {
	switch tag {
	case "breedvalidator":
		return "This breed is not listed."
	default:
		return "Default validation error"
	}

}
