package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Okwujohngibson25/movie-go/internal/apperrors"
)

const maxBodyBytes = 1 << 20 // 1 MiB; tune per endpoint

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return apperrors.ErrBodyTooLarge
		}
		return fmt.Errorf("%w: %v", apperrors.ErrInvalidJSON, err)
	}

	// Reject trailing garbage like `{"a":1}{"b":2}`
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w: body must contain a single JSON value", apperrors.ErrInvalidJSON)
	}
	return nil
}
