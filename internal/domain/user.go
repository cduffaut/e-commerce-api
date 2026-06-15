package domain

import "time"

// user in db
type User struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// info given by user (subscription)
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required, email"`
	Password string `json:"password" binding:"required"`
}

// user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required, email"`
	Password string `json:"password" binding:"required"`
}

// response after a register/login request
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// what we send back to user w/o password
type UserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
