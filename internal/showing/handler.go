package showing

import (
	"net/http"

	"github.com/Okwujohngibson25/movie-go/internal/apperrors"
	"github.com/Okwujohngibson25/movie-go/internal/helper"
	"github.com/Okwujohngibson25/movie-go/internal/response"
	"github.com/google/uuid"
)

type ShowingHandler struct {
	service *ShowingService
}

func NewShowingHandler(service *ShowingService) *ShowingHandler {
	return &ShowingHandler{service: service}
}

func (s *ShowingHandler) GetAllShowings(w http.ResponseWriter, r *http.Request) {
	showingData, err := s.service.GetAllShowings(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	apiResponse := response.APIResponse[any]{
		Code:    "FETCHED_ALL_SHOWINGS_SUCCESSFULLY",
		Message: "fetched all showings successfully",
		Data:    toShowingResponse(showingData),
	}

	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (s *ShowingHandler) CreateShowing(w http.ResponseWriter, r *http.Request) {
	var req CreateShowingDTO

	if err := helper.DecodeJSON(w, r, &req); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// convert movie id from Str to UUID
	movieID, err := uuid.Parse(req.MovieID)
	if err != nil {
		apperrors.HandleError(w, apperrors.ErrInvalidID)
		return
	}

	// convert screen id id from Str to UUID
	screenID, err := uuid.Parse(req.ScreenID)
	if err != nil {
		apperrors.HandleError(w, apperrors.ErrInvalidID)
		return
	}

	err = s.service.CreateShowing(r.Context(), CreateShowingInput{
		MovieID:   movieID,
		ScreenID:  screenID,
		StartTime: req.StartTime,
	})
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	apiResponse := response.APIResponse[any]{
		Code:    "SHOWING_CREATED",
		Message: "showing created successfully",
	}

	response.WriteJSON(w, http.StatusCreated, apiResponse)
}
