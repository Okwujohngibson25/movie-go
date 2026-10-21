package user

import (
	"net/http"
	"time"

	"github.com/Okwujohngibson25/movie-go/internal/apperrors"
	"github.com/Okwujohngibson25/movie-go/internal/helper"
	"github.com/Okwujohngibson25/movie-go/internal/response"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service  *UserService
	validate *validator.Validate
}

func NewUserHandler(service *UserService, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		service:  service,
		validate: validate,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := helper.DecodeJSON(w, r, &req); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	err := h.service.CreateUser(r.Context(), CreateUserInput{
		Fullname: req.Fullname,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	apiResponse := response.APIResponse[any]{
		Code:    "USER_CREATED_SUCCESSFULLY",
		Message: "user created successfully",
	}

	response.WriteJSON(w, http.StatusCreated, apiResponse)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginUserRequest

	if err := helper.DecodeJSON(w, r, &req); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	accessToken, refreshToken, err := h.service.LoginUser(r.Context(), LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()), // 7 days
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	apiResponse := response.APIResponse[any]{
		Code:    "LOGGED_IN_SUCCESSFUL",
		Message: "logged in successful",
		Data:    accessToken,
	}

	response.WriteJSON(w, http.StatusOK, apiResponse)
}
