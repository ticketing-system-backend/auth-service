package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Respond(c *gin.Context, code int, success bool, message string, payload interface{}) {
	key := "data"
	if !success {
		key = "errors"
	}
	c.JSON(code, gin.H{
		"success": success,
		"message": message,
		key:       payload,
	})
}

func ValidateAndFormat(data interface{}) []FieldError {
	if err := validate.Struct(data); err != nil {
		var res []FieldError
		for _, e := range err.(validator.ValidationErrors) {
			field := toSnakeCase(e.Field())
			msg := field + " is invalid"
			switch e.Tag() {
			case "required":
				msg = field + " is required"
			case "email":
				msg = "invalid email format"
			case "min":
				msg = field + " must be at least " + e.Param() + " characters"
			case "max":
				msg = field + " must be at most " + e.Param() + " characters"
			}
			res = append(res, FieldError{Field: field, Message: msg})
		}
		return res
	}
	return nil
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
