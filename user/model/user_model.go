package model

import (
	"time"
)

const (
	TableNameUser   = "users"
	TableNameFollow = "follows"
)

// User represents database model for users.
type User struct {
	ID        uint      `gorm:"column:user_id" json:"-"`
	Email     string    `gorm:"column:email" json:"email"`
	Name      string    `gorm:"column:name" json:"username"`
	Password  string    `gorm:"column:password" json:"-"`
	Bio       string    `gorm:"column:bio" json:"bio"`
	Image     string    `gorm:"column:image" json:"image"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	Disabled  bool      `gorm:"column:disabled" json:"-"`

	// Following is used for profile not database field.
	Following bool `gorm:"-"`
}

func (u User) TableName() string {
	return TableNameUser
}

// ToProfile converts current user to Profile.
func (u *User) ToProfile() *Profile {
	return &Profile{
		Username:  u.Name,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: u.Following,
	}
}

// Follow represents a database model for following relation between users.
type Follow struct {
	User      User `gorm:"foreignkey:UserID"`
	UserID    uint
	Follow    User `gorm:"foreignkey:FollowID"`
	FollowID  uint
	CreatedAt time.Time
}

func (f Follow) TableName() string {
	return TableNameFollow
}

// Profile represents user profile resource not database model.
type Profile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}
