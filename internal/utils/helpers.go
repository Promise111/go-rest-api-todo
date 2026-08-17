package utils

import "github.com/go-playground/validator/v10"

func JsonFieldNameLogin(fe validator.FieldError) string {
	switch fe.Field() {
	case "Email":
		return "email"
	case "Username":
		return "username"
	case "Password":
		return "password"
	default:
		return fe.Field()
	}
}

func ValidationMessageLogin(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "required_with":
		return "email or username is required"
	case "email":
		return "enter valid email address"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	default:
		return fe.Field() + " is invalid"
	}
}
