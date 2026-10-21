package movie

import (
	"fmt"
	"net/http"

	"github.com/Okwujohngibson25/movie-go/internal/apperrors"
	"github.com/Okwujohngibson25/movie-go/internal/helper"
	"github.com/Okwujohngibson25/movie-go/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type MovieHandler struct {
	service  *MovieService
	validate *validator.Validate
}

func NewMovieHandler(service *MovieService, validate *validator.Validate) *MovieHandler {
	return &MovieHandler{
		service:  service,
		validate: validate,
	}
}

func (s *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var req CreateMovieRequest

	if err := helper.DecodeJSON(w, r, &req); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		fmt.Printf("validation error: %+v\n", err)
		apperrors.HandleError(w, err)
		return
	}

	err := s.service.CreateMovie(r.Context(), MovieInput{
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		Rating:          req.Rating,
	})
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	apiResponse := response.APIResponse[any]{
		Code:    "MOVIE_CREATED_SUCCESSFULLY",
		Message: "movie created successfully",
	}

	response.WriteJSON(w, http.StatusCreated, apiResponse)
}

func (s *MovieHandler) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	movieData, err := s.service.GetAllMovies(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	apiResponse := response.APIResponse[any]{
		Code:    "MOVIES_FETCHED",
		Message: "Movies fetched successfully",
		Data:    toMovieResponse(movieData),
	}

	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (s *MovieHandler) GetMovieByID(w http.ResponseWriter, r *http.Request) {
	movieIDStr := chi.URLParam(r, "movieID")

	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		apperrors.HandleError(w, apperrors.ErrInvalidID)
		return
	}

	movieData, err := s.service.GetMovieByID(r.Context(), movieID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	apiResponse := response.APIResponse[any]{
		Code:    "MOVIE_FETCHED",
		Message: "Movie fetched successfully",
		Data:    toSingleMovieResponse(movieData),
	}

	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (s *MovieHandler) UpdateMovieByID(w http.ResponseWriter, r *http.Request) {
	movieIDStr := chi.URLParam(r, "movieID")

	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		apperrors.HandleError(w, apperrors.ErrInvalidID)
		return
	}

	var req UpdateMovieRequest

	if err := helper.DecodeJSON(w, r, &req); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	err = s.service.UpdateMovieByID(r.Context(), movieID, MovieInput{
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		Rating:          req.Rating,
	})
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	apiResponse := response.APIResponse[any]{
		Code:    "MOVIE_UPDATED",
		Message: "movie data updated successfully",
	}
	response.WriteJSON(w, http.StatusOK, apiResponse)
}

func (s *MovieHandler) DeleteMovieByID(w http.ResponseWriter, r *http.Request) {
	movieIdStr := chi.URLParam(r, "movieID")

	movieID, err := uuid.Parse(movieIdStr)
	if err != nil {
		apperrors.HandleError(w, apperrors.ErrInvalidID)
	}

	if err := s.service.DeleteMovieByID(r.Context(), movieID); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	apiResponse := response.APIResponse[any]{
		Code:    "MOVIE_DELETED_SUCCESSFULLY",
		Message: "movie deleted successfully",
	}

	response.WriteJSON(w, http.StatusOK, apiResponse)
}
