package auth

import "time"

type LoginResponse struct {
	Status string    `json:"status"`
	Msg    string    `json:"msg"`
	Data   LogedUser `json:"data"`
}

type LogedUser struct {
	UserId   uint   `json:"id"`
	UserName string `json:"username"`
	IsActive bool   `json:"isactive"`
	Role     string `json:"role"`
}

type Response struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

type User struct {
	ID           uint        `json:"id"`
	Email        string      `json:"email"`
	PasswordHash string      `json:"password_hash"`
	Username     string      `json:"username"`
	Role         string      `json:"role"`
	IsActive     bool        `json:"is_active"`
	LastLogin    time.Time   `json:"last_login"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Profile      UserProfile `json:"profile"`
}

type UserProfile struct {
	ID              uint       `json:"id"`
	UserID          uint       `json:"user_id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Phone           string     `json:"phone"`
	Dob             *time.Time `json:"dob"`
	Gender          string     `json:"gender"`
	Address         string     `json:"address"`
	City            string     `json:"city"`
	Country         string     `json:"country"`
	ProfileImageURL string     `json:"profile_image_url"`
	Bio             string     `json:"bio"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
