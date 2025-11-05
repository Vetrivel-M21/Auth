package auth

type UserInfo struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password_hash string `json:"password"`
	IsActive bool   `json:"is_active"`
	LastLogin string `json:"last_login"`
	Role	 string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

}

type Profile struct {
	UserName    string `json:"user_name"`
	Email       string `json:"email"`
	Bio        string `json:"bio"`
	DateofBirth string `json:"date_of_birth"`
	Location    string `json:"location"`
	ProfilePicture string `json:"profile_picture"`
	Description string `json:"description"`
}

type Response struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}
