package types

import (
	userModel "github.com/zacscoding/echo-gorm-realworld-app/internal/user/model"
)

// UserResponse represents User resource.
type UserResponse struct {
	User struct {
		Email    string `json:"email"`
		Token    string `json:"token"`
		Username string `json:"username"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"user"`
}

// ToUserResponse converts given model.User to UserResponse with JWT Token.
func ToUserResponse(u *userModel.User, token string) *UserResponse {
	user := new(UserResponse)
	user.User.Email = u.Email
	user.User.Token = token
	user.User.Username = u.Name
	user.User.Bio = u.Bio
	user.User.Image = u.Image
	return user
}

// UserProfile represents UserProfile resource.
type UserProfile struct {
	Profile struct {
		Username  string `json:"username"`
		Bio       string `json:"bio"`
		Image     string `json:"image"`
		Following bool   `json:"following"`
	} `json:"profile"`
}

// ToUserProfile converts given model.User to UserProfile.
func ToUserProfile(u *userModel.User) *UserProfile {
	up := new(UserProfile)
	up.Profile.Username = u.Name
	up.Profile.Bio = u.Bio
	up.Profile.Image = u.Image
	up.Profile.Following = u.Following
	return up
}
