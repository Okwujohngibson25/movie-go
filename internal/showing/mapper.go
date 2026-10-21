package showing

import "github.com/Okwujohngibson25/movie-go/internal/domains"

func toShowingResponse(showingData []domains.Showing) []ShowingResponseDTO {
	showingsData := make([]ShowingResponseDTO, 0, len(showingData))
	for _, value := range showingData {
		showingsData = append(showingsData, ShowingResponseDTO{
			ID:        value.ID,
			Movie:     value.Movie,
			Screen:    value.Screen,
			StartTime: value.StartTime,
			EndTime:   value.EndTime,
		})
	}

	return showingsData
}
