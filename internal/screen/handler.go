package screen

import (
	"net/http"

	"github.com/Okwujohngibson25/movie-go/internal/apperrors"
	"github.com/Okwujohngibson25/movie-go/internal/helper"
	"github.com/Okwujohngibson25/movie-go/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ScreenHandler struct {
	service  *ScreenService
	validate *validator.Validate
}

func NewScreenHandler(service *ScreenService, validate *validator.Validate) *ScreenHandler {
	return &ScreenHandler{
		service:  service,
		validate: validate,
	}
}

func (s *ScreenHandler) CreateScreen(w http.ResponseWriter, r *http.Request) {
	var req CreateScreenRequest

	if err := helper.DecodeJSON(w, r, &req); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	err := s.service.CreateScreen(r.Context(), CreateScreenInput{
		Name:     req.Name,
		Capacity: req.Capacity,
	})
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	apiResponse := response.APIResponse[any]{
		Code:    "SCREEN_CREATED_SUCCESSFULLY",
		Message: "screen created successfully",
	}
	response.WriteJSON(w, http.StatusCreated, apiResponse)
}

func (s *ScreenHandler) GetAllScreens(w http.ResponseWriter, r *http.Request) {
	screenData, err := s.service.GetAllScreens(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	apiResponse := response.APIResponse[any]{
		Code:    "SCREENS_FETCHED_SUCCESSFULLY",
		Message: "Screen fetched successfully",
		Data:    toScreenResponse(screenData),
	}

	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (s *ScreenHandler) GetScreensByID(w http.ResponseWriter, r *http.Request) {
	screenIDStr := chi.URLParam(r, "screensID")

	screenID, err := uuid.Parse(screenIDStr)
	if err != nil {
		apperrors.HandleError(w, apperrors.ErrInvalidID)
		return
	}

	screen, err := s.service.GetScreensByID(r.Context(), screenID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	apiResponse := response.APIResponse[any]{
		Code:    "SCREEN_FETCHED_SUCCESSFULLY",
		Message: "screen fetched successfully",
		Data:    toSingleScreenResponse(screen),
	}

	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (s *ScreenHandler) UpdateScreensByID(w http.ResponseWriter, r *http.Request) {
	var req UpdateScreenRequest

	screenIDStr := chi.URLParam(r, "screensID")

	screenID, err := uuid.Parse(screenIDStr)
	if err != nil {
		apperrors.HandleError(w, apperrors.ErrInvalidID)
		return
	}

	if err := helper.DecodeJSON(w, r, &req); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	err = s.service.UpdateScreensByID(r.Context(), screenID, UpdateScreenInput{
		Name: req.Name,
	})
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	apiResponse := response.APIResponse[any]{
		Code:    "SCREEN_UPDATED_SUCCESSFULLY",
		Message: "screen updated successfully",
	}

	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (s *ScreenHandler) GetSeatsByScreenID(w http.ResponseWriter, r *http.Request) {
	screenIDStr := chi.URLParam(r, "screensID")

	screenID, err := uuid.Parse(screenIDStr)
	if err != nil {
		apperrors.HandleError(w, apperrors.ErrInvalidID)
		return
	}

	seatData, err := s.service.GetSeatsByScreenID(r.Context(), screenID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	apiResponse := response.APIResponse[any]{
		Code:    "SEATS_FETCHED_SUCCESSFULLY",
		Message: "Seats retrieved for screen " + screenIDStr,
		Data:    toSeatsResponse(seatData),
	}

	response.WriteJSON(w, http.StatusOK, apiResponse)
}
