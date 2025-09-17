package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Validator wraps the validator instance
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

// ValidateStruct validates a struct and returns formatted error messages
func (v *Validator) ValidateStruct(s interface{}) map[string]string {
	errors := make(map[string]string)

	err := v.validator.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := getJSONFieldName(s, err.Field())
			errors[fieldName] = getErrorMessage(err)
		}
	}

	return errors
}

// getJSONFieldName returns the JSON field name for a struct field
func getJSONFieldName(s interface{}, fieldName string) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, found := t.FieldByName(fieldName)
	if !found {
		return strings.ToLower(fieldName)
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return strings.ToLower(fieldName)
	}

	// Handle json:"field_name,omitempty"
	parts := strings.Split(jsonTag, ",")
	if parts[0] == "-" {
		return strings.ToLower(fieldName)
	}

	return parts[0]
}

// getErrorMessage returns a user-friendly error message for validation errors
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("Must be at least %s characters long", err.Param())
		}
		return fmt.Sprintf("Must be at least %s", err.Param())
	case "max":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("Must be at most %s characters long", err.Param())
		}
		return fmt.Sprintf("Must be at most %s", err.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", err.Param())
	case "unique":
		return "This value already exists"
	default:
		return fmt.Sprintf("Invalid value for %s", err.Field())
	}
}

// BindAndValidate binds request data and validates it
func BindAndValidate(c *gin.Context, obj interface{}) map[string]string {
	if err := c.ShouldBindJSON(obj); err != nil {
		return map[string]string{"binding": err.Error()}
	}

	validator := NewValidator()
	return validator.ValidateStruct(obj)
}
