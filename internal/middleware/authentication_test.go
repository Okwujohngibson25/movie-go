package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const testSecret = "test-secret"

func sign(t *testing.T, method jwt.SigningMethod, key any, claims jwt.MapClaims) string {
	t.Helper()
	token, err := jwt.NewWithClaims(method, claims).SignedString(key)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return token
}

func claimsWith(typ string, role string, exp time.Time) jwt.MapClaims {
	c := jwt.MapClaims{
		"email":  "user@example.com",
		"userid": uuid.NewString(),
		"Role":   role,
		"exp":    exp.Unix(),
	}
	if typ != "" {
		c["typ"] = typ
	}
	return c
}

func responseCode(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response body %q: %v", rec.Body.String(), err)
	}
	return body.Code
}

func TestRequireAuth(t *testing.T) {
	future := time.Now().Add(time.Hour)
	past := time.Now().Add(-time.Hour)

	tests := []struct {
		name       string
		token      string
		wantNext   bool
		wantStatus int
		wantCode   string
	}{
		{
			name:       "missing header",
			token:      "",
			wantStatus: http.StatusUnauthorized,
			wantCode:   "AUTH_TOKEN_MISSING",
		},
		{
			name:       "bad signature",
			token:      sign(t, jwt.SigningMethodHS256, []byte("wrong-secret"), claimsWith("access", "admin", future)),
			wantStatus: http.StatusUnauthorized,
			wantCode:   "INVALID_AUTH_TOKEN",
		},
		{
			name:       "expired",
			token:      sign(t, jwt.SigningMethodHS256, []byte(testSecret), claimsWith("access", "admin", past)),
			wantStatus: http.StatusUnauthorized,
			wantCode:   "INVALID_AUTH_TOKEN",
		},
		{
			name:       "non-HMAC signing method",
			token:      sign(t, jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, claimsWith("access", "admin", future)),
			wantStatus: http.StatusUnauthorized,
			wantCode:   "INVALID_AUTH_TOKEN",
		},
		{
			name:       "refresh token",
			token:      sign(t, jwt.SigningMethodHS256, []byte(testSecret), claimsWith("refresh", "admin", future)),
			wantStatus: http.StatusUnauthorized,
			wantCode:   "INVALID_AUTH_TOKEN",
		},
		{
			// tokens issued before the typ claim existed must not be accepted
			name:       "no typ claim",
			token:      sign(t, jwt.SigningMethodHS256, []byte(testSecret), claimsWith("", "admin", future)),
			wantStatus: http.StatusUnauthorized,
			wantCode:   "INVALID_AUTH_TOKEN",
		},
		{
			name:       "access token",
			token:      sign(t, jwt.SigningMethodHS256, []byte(testSecret), claimsWith("access", "admin", future)),
			wantNext:   true,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nextCalled bool
			var gotClaims *CustomClaims
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				gotClaims, _ = r.Context().Value(claimsKey).(*CustomClaims)
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			rec := httptest.NewRecorder()

			RequireAuth(testSecret)(next).ServeHTTP(rec, req)

			if nextCalled != tt.wantNext {
				t.Fatalf("next called = %v, want %v", nextCalled, tt.wantNext)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantNext {
				if gotClaims == nil || gotClaims.Email != "user@example.com" {
					t.Fatalf("claims in context = %+v, want email user@example.com", gotClaims)
				}
				return
			}
			if code := responseCode(t, rec); code != tt.wantCode {
				t.Fatalf("code = %q, want %q", code, tt.wantCode)
			}
		})
	}
}

func TestRequireAuthThenRequireAdmin(t *testing.T) {
	future := time.Now().Add(time.Hour)

	tests := []struct {
		name       string
		role       string
		wantStatus int
		wantCode   string
	}{
		{name: "admin", role: "admin", wantStatus: http.StatusOK},
		{name: "regular user", role: "user", wantStatus: http.StatusForbidden, wantCode: "UNAUTHORIZED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			handler := RequireAuth(testSecret)(RequireAdmin(ok))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", sign(t, jwt.SigningMethodHS256, []byte(testSecret), claimsWith("access", tt.role, future)))
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantCode != "" {
				if code := responseCode(t, rec); code != tt.wantCode {
					t.Fatalf("code = %q, want %q", code, tt.wantCode)
				}
			}
		})
	}
}
