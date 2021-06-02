package model

import (
	"time"
)

const (
	UserTableName   = "users"
	FollowTableName = "follows"
)

// User represents database model for users.
type User struct {
	ID        uint      `gorm:"column:user_id" json:"-"`
	Email     string    `gorm:"column:email" json:"email"`
	Name      string    `gorm:"column:name" json:"name"`
	Password  string    `gorm:"column:password" json:"-"`
	Bio       string    `gorm:"column:bio" json:"bio"`
	Image     string    `gorm:"column:image" json:"image"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	Disabled  bool      `gorm:"column:disabled" json:"-"`
}

func (u User) TableName() string {
	return UserTableName
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
	return FollowTableName
}
