package domains

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Fullname     string
	Email        string
	Role         string
	PasswordHash string
}
