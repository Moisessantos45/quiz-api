package validator

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var URLRegex = regexp.MustCompile(`^https?:\/\/.+\..+`)

func ValidateURL(link string, fieldName string) error {
	if strings.TrimSpace(link) == "" {
		return nil
	}

	if !URLRegex.MatchString(link) {
		return fmt.Errorf("el campo %s tiene un formato de URL inválido", fieldName)
	}

	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("el campo %s tiene un formato de URL inválido: %s", fieldName, err.Error())
	}

	return nil
}

func ValidateJSONFormat(data string, fieldName string) error {
	if strings.TrimSpace(data) == "" {
		return fmt.Errorf("el campo %s es obligatorio y debe tener un formato JSON válido", fieldName)
	}

	var js json.RawMessage
	if err := json.Unmarshal([]byte(data), &js); err != nil {
		return fmt.Errorf("el campo %s tiene un formato JSON inválido: %s", fieldName, err.Error())
	}

	return nil
}

func ValidateJSONArray(data string, fieldName string) error {
	if strings.TrimSpace(data) == "" {
		return fmt.Errorf("el campo %s es obligatorio y debe ser un array JSON válido", fieldName)
	}

	var arr []any
	if err := json.Unmarshal([]byte(data), &arr); err != nil {
		return fmt.Errorf("el campo %s debe ser un array JSON válido: %s", fieldName, err.Error())
	}

	return nil
}

func ValidateJSONObject(data string, fieldName string) error {
	if strings.TrimSpace(data) == "" {
		return fmt.Errorf("el campo %s es obligatorio y debe ser un objeto JSON válido", fieldName)
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(data), &obj); err != nil {
		return fmt.Errorf("el campo %s debe ser un objeto JSON válido: %s", fieldName, err.Error())
	}

	return nil
}
