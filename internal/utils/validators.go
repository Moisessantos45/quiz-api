package utils

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

var PhoneRegex = regexp.MustCompile(`^\+?[0-9\s\-\(\)]{7,20}$`)

func ValidateParamsId(c *gin.Context, params string) (uint64, error) {
	newParams := params
	if len(strings.TrimSpace(params)) == 0 {
		newParams = "id"
	}

	idStr := c.Param(newParams)
	if len(strings.TrimSpace(idStr)) == 0 {
		return 0, errors.New("ID de producto no proporcionado")
	}

	id, newErr := strconv.ParseUint(idStr, 10, 64)
	if newErr != nil {
		return 0, errors.New("ID de producto inválido: " + newErr.Error())
	}

	return id, nil
}

func ExtractedParamsJwt(c *gin.Context) (string, string, error) {

	emailIface := c.MustGet("email")
	emailStr, ok := emailIface.(string)
	if !ok {
		log.Printf("UserId no es string: %T %+v", emailIface, emailStr)
		return "", "", fmt.Errorf("email no es string")
	}

	tokenIface := c.MustGet("token")
	tokenStr, ok := tokenIface.(string)
	if !ok {
		return "", "", fmt.Errorf("Token no es string")
	}

	return tokenStr, emailStr, nil
}

func ValidateQueryPagination(c *gin.Context) (int, int, error) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, 0, errors.New("Parámetro 'page' inválido: debe ser un número entero positivo")
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		return 0, 0, errors.New("Parámetro 'page_size' inválido: debe ser un número entero positivo")
	}

	return page, pageSize, nil
}
