package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(req interface{}) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	value := reflect.ValueOf(req)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return errors.New("request cannot be nil")
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return errors.New("request must be a struct")
	}

	return v.validateStruct(value)
}

func (v *Validator) validateStruct(value reflect.Value) error {
	structType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := structType.Field(i)
		
		if !field.CanInterface() {
			continue
		}

		tag := fieldType.Tag.Get("binding")
		jsonTag := fieldType.Tag.Get("json")
		
		fieldName := fieldType.Name
		if jsonTag != "" {
			if parts := strings.Split(jsonTag, ","); len(parts) > 0 && parts[0] != "" {
				fieldName = parts[0]
			}
		}

		if strings.Contains(tag, "required") {
			if err := v.validateRequired(field, fieldName); err != nil {
				return err
			}
		}

		if err := v.validateFieldType(field, fieldType, fieldName); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) validateRequired(field reflect.Value, fieldName string) error {
	switch field.Kind() {
	case reflect.String:
		if strings.TrimSpace(field.String()) == "" {
			return fmt.Errorf("field '%s' is required", fieldName)
		}
	case reflect.Slice, reflect.Array:
		if field.Len() == 0 {
			return fmt.Errorf("field '%s' is required", fieldName)
		}
	case reflect.Ptr, reflect.Interface:
		if field.IsNil() {
			return fmt.Errorf("field '%s' is required", fieldName)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Bool:
	default:
		if field.IsZero() {
			return fmt.Errorf("field '%s' is required", fieldName)
		}
	}
	return nil
}

func (v *Validator) validateFieldType(field reflect.Value, fieldType reflect.StructField, fieldName string) error {
	switch fieldName {
	case "url", "webhook", "webhookurl", "webhookURL":
		if field.Kind() == reflect.String && field.String() != "" {
			return v.validateURL(field.String(), fieldName)
		}
	case "phone":
		if field.Kind() == reflect.String && field.String() != "" {
			return v.validatePhoneNumber(field.String(), fieldName)
		}
	case "email":
		if field.Kind() == reflect.String && field.String() != "" {
			return v.validateEmail(field.String(), fieldName)
		}
	}
	return nil
}

func (v *Validator) validateURL(url, fieldName string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return nil // Empty is valid if not required
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("field '%s' must be a valid URL (http:// or https://)", fieldName)
	}

	if len(url) < 10 { // Minimum: https://a.b
		return fmt.Errorf("field '%s' must be a valid URL", fieldName)
	}

	return nil
}

func (v *Validator) validatePhoneNumber(phone, fieldName string) error {
	if phone == "" {
		return fmt.Errorf("field '%s' cannot be empty", fieldName)
	}

	if len(phone) < 7 || len(phone) > 20 {
		return fmt.Errorf("field '%s' must be between 7 and 20 characters", fieldName)
	}

	return nil
}

func (v *Validator) validateEmail(email, fieldName string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("field '%s' must be a valid email address", fieldName)
	}

	return nil
}

func (v *Validator) ValidateStringSlice(slice []string, fieldName string) error {
	if len(slice) == 0 {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}

	for i, item := range slice {
		if strings.TrimSpace(item) == "" {
			return fmt.Errorf("%s[%d] cannot be empty", fieldName, i)
		}
	}

	return nil
}

var DefaultValidator = NewValidator()
