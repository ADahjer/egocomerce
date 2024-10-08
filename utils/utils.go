package utils

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"

	"github.com/ADahjer/egocomerce/types"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	maxFileSize = 3 * 1024 * 1024 //3MB
)

func ApiErrorHandler(err error, c echo.Context) {

	switch e := err.(type) {
	case types.ApiError:

		c.JSON(e.Status, types.Map{
			"Error":       e.Msg,
			"Status Code": e.Status,
		})

	case validator.ValidationErrors:
		var errsMsg []string

		for _, err := range e {
			switch err.Tag() {
			case "min":
				errsMsg = append(errsMsg, fmt.Sprintf("%s should be at least %s characters long", err.Field(), err.Param()))
			case "required":
				errsMsg = append(errsMsg, fmt.Sprintf("%s field is required", err.Field()))
			case "email":
				errsMsg = append(errsMsg, fmt.Sprintf("%s field should be a valid email", err.Field()))
			default:
				errsMsg = append(errsMsg, err.Error())
			}
		}

		c.JSON(http.StatusBadRequest, types.Map{"Errors": errsMsg})

	default:
		c.JSON(http.StatusInternalServerError, types.Map{"Error": err.Error()})
	}

}

func ValidatePassword(candidate string) bool {
	if len(candidate) < 8 {
		return false
	}

	lowercase := regexp.MustCompile(`[a-z]`).MatchString(candidate)
	uppercase := regexp.MustCompile(`[A-Z]`).MatchString(candidate)
	number := regexp.MustCompile(`\d`).MatchString(candidate)
	specialChar := regexp.MustCompile(`[\W_]`).MatchString(candidate)

	return lowercase && uppercase && number && specialChar
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePassword(candidate string, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(candidate))

	if err != nil {
		return false, err
	}

	return true, nil
}

func ValidateImageType(fileHeader *multipart.FileHeader) (string, bool) {
	mimeType := fileHeader.Header.Get("Content-Type")
	if fileHeader.Size > maxFileSize {
		return "maximun file size supported is 3MB", false
	}

	if mimeType != "image/png" && mimeType != "image/jpeg" && mimeType != "image/jpg" {
		return "file types allowed are image/png and image/jpg", false
	}

	return "", true
}
