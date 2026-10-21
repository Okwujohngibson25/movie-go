package user

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strings"
	"testing"

	"github.com/Okwujohngibson25/movie-go/internal/apperrors"
	"github.com/Okwujohngibson25/movie-go/internal/config"
	"github.com/Okwujohngibson25/movie-go/internal/domains"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const testSecret = "test-secret"

type fakeRepo struct {
	user      domains.User
	findErr   error
	createErr error
	created   []domains.User
}

func (f *fakeRepo) CreateUser(ctx context.Context, userData domains.User) error {
	f.created = append(f.created, userData)
	return f.createErr
}

func (f *fakeRepo) FindByEmail(ctx context.Context, email string) (domains.User, error) {
	if f.findErr != nil {
		return domains.User{}, f.findErr
	}
	return f.user, nil
}

func testConfig(adminEmail, adminPassword string) *config.Config {
	return &config.Config{
		App: config.AppConfig{
			JWTSecret:     testSecret,
			AdminEmail:    adminEmail,
			AdminPassword: adminPassword,
		},
	}
}

func tokenType(t *testing.T, signed string) any {
	t.Helper()
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(signed, claims, func(token *jwt.Token) (any, error) {
		return []byte(testSecret), nil
	})
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	return claims["typ"]
}

func TestLoginUser_TokenTypes(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	repo := &fakeRepo{user: domains.User{
		ID:           uuid.New(),
		Fullname:     "Test User",
		Email:        "user@example.com",
		Role:         "user",
		PasswordHash: string(hash),
	}}
	svc := NewUserService(repo, testConfig("", ""))

	access, refresh, err := svc.LoginUser(context.Background(), LoginUserInput{
		Email:    "user@example.com",
		Password: "secret-password",
	})
	if err != nil {
		t.Fatalf("LoginUser: %v", err)
	}

	if got := tokenType(t, access); got != "access" {
		t.Errorf("access token typ = %v, want \"access\"", got)
	}
	if got := tokenType(t, refresh); got != "refresh" {
		t.Errorf("refresh token typ = %v, want \"refresh\"", got)
	}
}

// captureLog redirects the standard logger into a buffer for the duration of the test.
func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prevOut, prevFlags := log.Writer(), log.Flags()
	log.SetOutput(&buf)
	t.Cleanup(func() {
		log.SetOutput(prevOut)
		log.SetFlags(prevFlags)
	})
	return &buf
}

func TestSeedAdmin_CreateFailsDoesNotLogSuccess(t *testing.T) {
	buf := captureLog(t)
	repo := &fakeRepo{
		findErr:   apperrors.ErrEmailDoesNotExist,
		createErr: errors.New("insert failed"),
	}
	svc := NewUserService(repo, testConfig("admin@example.com", "admin-password"))

	svc.SeedAdmin(testConfig("admin@example.com", "admin-password"))

	if strings.Contains(buf.String(), "created successfully") {
		t.Fatalf("log reports success after CreateUser failed:\n%s", buf.String())
	}
}

func TestSeedAdmin_SkipsWhenCredentialsMissing(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
	}{
		{name: "empty email", email: "", password: "admin-password"},
		{name: "empty password", email: "admin@example.com", password: ""},
		{name: "both empty", email: "", password: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := captureLog(t)
			repo := &fakeRepo{findErr: apperrors.ErrEmailDoesNotExist}
			cfg := testConfig(tt.email, tt.password)
			svc := NewUserService(repo, cfg)

			svc.SeedAdmin(cfg)

			if len(repo.created) != 0 {
				t.Fatalf("CreateUser called %d times, want 0", len(repo.created))
			}
			if !strings.Contains(strings.ToLower(buf.String()), "skip") {
				t.Fatalf("log output %q does not contain a skip warning", buf.String())
			}
		})
	}
}

func TestSeedAdmin_AdminExistsDoesNotCreate(t *testing.T) {
	captureLog(t)
	repo := &fakeRepo{user: domains.User{Email: "admin@example.com", Role: "admin"}}
	cfg := testConfig("admin@example.com", "admin-password")
	svc := NewUserService(repo, cfg)

	svc.SeedAdmin(cfg)

	if len(repo.created) != 0 {
		t.Fatalf("CreateUser called %d times, want 0", len(repo.created))
	}
}

func TestSeedAdmin_CreatesAdminWhenMissing(t *testing.T) {
	captureLog(t)
	repo := &fakeRepo{findErr: apperrors.ErrEmailDoesNotExist}
	cfg := testConfig("Admin@Example.com", "admin-password")
	svc := NewUserService(repo, cfg)

	svc.SeedAdmin(cfg)

	if len(repo.created) != 1 {
		t.Fatalf("CreateUser called %d times, want 1", len(repo.created))
	}
	admin := repo.created[0]
	if admin.Role != "admin" {
		t.Errorf("role = %q, want \"admin\"", admin.Role)
	}
	if admin.Email != "admin@example.com" {
		t.Errorf("email = %q, want lowercased \"admin@example.com\"", admin.Email)
	}
	if admin.PasswordHash == "admin-password" {
		t.Fatal("password stored in plain text")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte("admin-password")); err != nil {
		t.Errorf("password hash does not match configured password: %v", err)
	}
}
