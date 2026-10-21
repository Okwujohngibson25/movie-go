package user

import (
	"context"
	"testing"

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
