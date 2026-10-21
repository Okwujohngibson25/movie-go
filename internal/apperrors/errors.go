package apperrors

import (
	"errors"
	"log"
	"net/http"

	"github.com/Okwujohngibson25/movie-go/internal/response"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInternalServerError     = errors.New("internal server error")
	ErrEmailAlreadyExist       = errors.New("email already exist")
	ErrInvalidLoginCredentials = errors.New("invalid login credentials")
	ErrEmailDoesNotExist       = errors.New("email does not exist")
	ErrTokenMissing            = errors.New("token is missing")
	ErrInvalidToken            = errors.New("token is invalid")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrInvalidID               = errors.New("invalid ID")
	ErrNotFound                = errors.New("resource not found")
	ErrMovieNotFound           = errors.New("movie not found")
	ErrScreenNotFound          = errors.New("screen not found")
	ErrMovieAlreadyExists      = errors.New("movie already exist")
	ErrInvalidJSON             = errors.New("invalid JSON")
	ErrBodyTooLarge            = errors.New("request body too large")
)

func HandleError(w http.ResponseWriter, err error) {
	var ValidationErrors validator.ValidationErrors

	switch {
	case errors.As(err, &ValidationErrors):
		details := make([]response.ErrorDetails, 0)

		for _, value := range ValidationErrors {
			details = append(details, response.ErrorDetails{
				Field:   value.Field(),
				Message: value.Tag(),
			})
		}
		apiResponse := response.APIResponse[any]{
			Code:    "VALIDATION_ERROR",
			Message: "validation failed",
			Details: details,
		}
		response.WriteJSON(w, http.StatusBadRequest, apiResponse)

	case errors.Is(err, ErrInvalidJSON):
		apiResponse := response.APIResponse[any]{
			Code:    "INVALID_JSON",
			Message: "request body contains invalid json",
		}
		response.WriteJSON(w, http.StatusBadRequest, apiResponse)
	case errors.Is(err, ErrBodyTooLarge):
		apiResponse := response.APIResponse[any]{
			Code:    "REQUEST_BODY_TOO_LARGE",
			Message: "request body is too large",
		}
		response.WriteJSON(w, http.StatusConflict, apiResponse)
	case errors.Is(err, ErrEmailAlreadyExist):
		apiResponse := response.APIResponse[any]{
			Code:    "EMAIL_ALREADY_EXIST",
			Message: "email already exist",
		}
		response.WriteJSON(w, http.StatusConflict, apiResponse)
	case errors.Is(err, ErrInvalidLoginCredentials):
		apiResponse := response.APIResponse[any]{
			Code:    "INVALID_LOGIN_CREDENTAILS",
			Message: "invalid email or password credentials",
		}
		response.WriteJSON(w, http.StatusUnauthorized, apiResponse)
	case errors.Is(err, ErrTokenMissing):
		apiResponse := response.APIResponse[any]{
			Code:    "AUTH_TOKEN_MISSING",
			Message: "authentication token is required",
		}
		response.WriteJSON(w, http.StatusUnauthorized, apiResponse)
	case errors.Is(err, ErrInvalidToken):
		apiResponse := response.APIResponse[any]{
			Code:    "INVALID_AUTH_TOKEN",
			Message: "auth token is invalid",
		}
		response.WriteJSON(w, http.StatusUnauthorized, apiResponse)
	case errors.Is(err, ErrUnauthorized):
		apiResponse := response.APIResponse[any]{
			Code:    "UNAUTHORIZED",
			Message: "You are not authorized to perform this action",
		}
		response.WriteJSON(w, http.StatusForbidden, apiResponse)
	case errors.Is(err, ErrInvalidID):
		apiResponse := response.APIResponse[any]{
			Code:    "INVALID_ID",
			Message: "invalid request ID",
		}
		response.WriteJSON(w, http.StatusBadRequest, apiResponse)
	case errors.Is(err, ErrNotFound):
		apiResponse := response.APIResponse[any]{
			Code:    "RESOURCE_NOT_FOUND",
			Message: "resource not found",
		}
		response.WriteJSON(w, http.StatusNotFound, apiResponse)
	case errors.Is(err, ErrMovieNotFound):
		apiResponse := response.APIResponse[any]{
			Code:    "MOVIE_NOT_FOUND",
			Message: "movie not found",
		}
		response.WriteJSON(w, http.StatusNotFound, apiResponse)
	case errors.Is(err, ErrScreenNotFound):
		apiResponse := response.APIResponse[any]{
			Code:    "SCREEN_NOT_FOUND",
			Message: "screen not found",
		}
		response.WriteJSON(w, http.StatusNotFound, apiResponse)
	default:
		log.Printf("internal error: %v", err)
		apiResponse := response.APIResponse[any]{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "internal server error",
		}
		response.WriteJSON(w, http.StatusInternalServerError, apiResponse)
	}
}
