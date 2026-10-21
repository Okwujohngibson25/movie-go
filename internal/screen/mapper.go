package screen

import "github.com/Okwujohngibson25/movie-go/internal/domains"

func toScreenResponse(screens []domains.Screen) []ScreenResponseDTO {
	screensData := make([]ScreenResponseDTO, 0, len(screens))

	for _, value := range screens {
		screensData = append(screensData, ScreenResponseDTO{
			ID:       value.ID,
			Name:     value.Name,
			Capacity: value.Capacity,
		})
	}

	return screensData
}

func toSingleScreenResponse(screen domains.Screen) ScreenResponseDTO {
	return ScreenResponseDTO{
		ID:       screen.ID,
		Name:     screen.Name,
		Capacity: screen.Capacity,
	}
}

func toSeatsResponse(seats []domains.Seats) []SeatsResponseDTO {
	seatsData := make([]SeatsResponseDTO, 0, len(seats))

	for _, value := range seats {
		seatsData = append(seatsData, SeatsResponseDTO{
			ID:         value.ID,
			ScreenID:   value.ScreenID,
			SeatNumber: value.SeatNumber,
		})
	}

	return seatsData
}
