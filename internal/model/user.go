package model

import "time"

type UserRole uint8

const (
	UserRoleUser UserRole = iota
	UserRoleAdmin
)

func (u UserRole) String() string {
	switch u {
	case UserRoleUser:
		return "User"
	case UserRoleAdmin:
		return "Admin"
	default:
		return "Unknown"
	}
}

type User struct {
	ID           uint64    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Role         UserRole  `json:"role" db:"role"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserRecord struct {
	Username     string
	Role         UserRole
	PasswordHash string
	CreatedAt    time.Time
}

type UpdateUserRecord struct {
	Username     *string
	Role         *UserRole
	PasswordHash *string
}
