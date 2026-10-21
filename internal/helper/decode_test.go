package helper

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Okwujohngibson25/movie-go/internal/apperrors"
)

func TestDecodeJSON_RejectsUnknownFieldAndOversize(t *testing.T) {
	type req struct {
		DurationMinutes int `json:"duration_minutes"`
	}

	// typo'd field -> error, not a silent zero
	r := httptest.NewRequest("POST", "/", strings.NewReader(`{"duration_mintes": 30}`))
	if err := DecodeJSON(httptest.NewRecorder(), r, &req{}); !errors.Is(err, apperrors.ErrInvalidJSON) {
		t.Fatalf("want ErrInvalidJSON, got %v", err)
	}

	// oversize -> ErrBodyTooLarge
	big := strings.NewReader(`{"duration_minutes":1,"pad":"` + strings.Repeat("x", maxBodyBytes) + `"}`)
	r = httptest.NewRequest("POST", "/", big)
	if err := DecodeJSON(httptest.NewRecorder(), r, &req{}); !errors.Is(err, apperrors.ErrBodyTooLarge) {
		t.Fatalf("want ErrBodyTooLarge, got %v", err)
	}
}
