package user

import (
	"github.com/labstack/echo/v4"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/hashutils"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/httputils"
)

// SignUpRequest represents request body data of an user registration.
type SignUpRequest struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"user" validate:"required"`
}

func (r *SignUpRequest) Bind(ctx echo.Context, u *userModel.User) error {
	if err := httputils.BindAndValidate(ctx, r); err != nil {
		return err
	}
	password, err := hashutils.EncodePassword(r.User.Password)
	if err != nil {
		return err
	}
	u.Name = r.User.Username
	u.Email = r.User.Email
	u.Password = password
	return nil
}

// SignInRequest represents request body data of sign in.
type SignInRequest struct {
	User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"user" validate:"required"`
}

func (r *SignInRequest) Bind(ctx echo.Context) error {
	return httputils.BindAndValidate(ctx, r)
}

type UpdateUserRequest struct {
	User struct {
		Username string `json:"username" validate:"omitempty"`
		Email    string `json:"email" validate:"omitempty"`
		Password string `json:"password" validate:"omitempty"`
		Bio      string `json:"bio" validate:"omitempty"`
		Image    string `json:"image" validate:"omitempty"`
	} `json:"user" validate:"required"`
}

func (r *UpdateUserRequest) Bind(ctx echo.Context, u *userModel.User) error {
	if err := httputils.BindAndValidate(ctx, r); err != nil {
		return err
	}
	if r.User.Username != "" {
		u.Name = r.User.Username
	}
	if r.User.Email != "" {
		u.Name = r.User.Email
	}
	if r.User.Password != "" {
		password, err := hashutils.EncodePassword(r.User.Password)
		if err != nil {
			return err
		}
		u.Password = password
	}
	if r.User.Bio != "" {
		u.Bio = r.User.Bio
	}
	if r.User.Image != "" {
		u.Image = r.User.Image
	}
	return nil
}
