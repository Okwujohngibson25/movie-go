package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Okwujohngibson25/movie-go/internal/apperrors"
	"github.com/Okwujohngibson25/movie-go/internal/config"
	"github.com/Okwujohngibson25/movie-go/internal/domains"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userData domains.User) error
	FindByEmail(ctx context.Context, email string) (domains.User, error)
}

type UserService struct {
	repository UserRepository
	config     *config.Config
}

func NewUserService(repository UserRepository, cfg *config.Config) *UserService {
	return &UserService{
		repository: repository,
		config:     cfg,
	}
}

type CustomClaims struct {
	Email  string
	UserID uuid.UUID
	Role   string
	jwt.RegisteredClaims
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserInput) error {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	passwordString := string(hasedPassword)
	userData := domains.User{
		ID:           uuid.New(),
		Fullname:     req.Fullname,
		Email:        req.Email,
		Role:         "user",
		PasswordHash: passwordString,
	}

	if err := s.repository.CreateUser(ctx, userData); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *UserService) LoginUser(ctx context.Context, req LoginUserInput) (signedAccess string, rawRefresh string, errordata error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	// Fetch user
	user, err := s.repository.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrEmailDoesNotExist) {
			return "", "", apperrors.ErrInvalidLoginCredentials
		}
		return "", "", err
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", "", apperrors.ErrInvalidLoginCredentials
	}

	// create jwt token
	// create access token
	accessClaims := CustomClaims{
		Email:  user.Email,
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	// sign access token with secret key
	signedAccess, err = accessToken.SignedString([]byte(s.config.App.JWTSecret))
	if err != nil {
		log.Printf("unable to sign Jwt with secretKey: %v", err)
		return "", "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	// create refresh token
	refreshClaims := CustomClaims{
		Email:  user.Email,
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	// sign refresh token with secret key
	signedRefresh, err := refreshToken.SignedString([]byte(s.config.App.JWTSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	// return accessToken and refreshToken
	return signedAccess, signedRefresh, nil
}

func (s *UserService) SeedAdmin(cfg *config.Config) {
	_, err := s.repository.FindByEmail(context.Background(), cfg.App.AdminEmail)
	if err != nil {
		// if Admin account does not exist... Seed Admin account
		if errors.Is(err, apperrors.ErrEmailDoesNotExist) {
			password, err := bcrypt.GenerateFromPassword([]byte(cfg.App.AdminPassword), 10)
			if err != nil {
				log.Println("failed to hash password: %v", err)
				return
			}
			PasswordString := string(password)
			email := strings.ToLower(cfg.App.AdminEmail)

			adminData := domains.User{
				ID:           uuid.New(),
				Fullname:     "Administrator",
				Email:        email,
				Role:         "admin",
				PasswordHash: PasswordString,
			}

			if err := s.repository.CreateUser(context.Background(), adminData); err != nil {
				log.Println("failed to create admin account: %v", err)
			}
			log.Println("Admin account created successfully")
			return
		}
		log.Printf("failed to find admin by email: %v", err)
		return
	}

	log.Println("admin account exist, skipping seeding")
}
