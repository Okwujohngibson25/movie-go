package movie

import "github.com/Okwujohngibson25/movie-go/internal/domains"

func toMovieResponse(movies []domains.Movie) []MovieResponseDTO {
	moviesData := make([]MovieResponseDTO, 0, len(movies))

	for _, value := range movies {
		moviesData = append(moviesData, MovieResponseDTO{
			ID:              value.ID,
			Title:           value.Title,
			Description:     value.Description,
			DurationMinutes: value.DurationMinutes,
			Rating:          value.Rating,
		})
	}

	return moviesData
}

func toSingleMovieResponse(movie domains.Movie) MovieResponseDTO {
	return MovieResponseDTO{
		ID:              movie.ID,
		Title:           movie.Title,
		Description:     movie.Description,
		DurationMinutes: movie.DurationMinutes,
		Rating:          movie.Rating,
	}
}
