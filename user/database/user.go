package database

import (
	"context"
	"github.com/zacscoding/echo-gorm-realworld-app/database"
	"github.com/zacscoding/echo-gorm-realworld-app/logging"
	"github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"gorm.io/gorm"
	"time"
)

//go:generate mockery --name UserDB --filename user_mock.go
type UserDB interface {
	// Save saves a given user usr.
	// database.ErrKeyConflict will be returned if duplicate emails.
	Save(ctx context.Context, u *model.User) error

	// Update updates given model.User from id.
	// database.ErrRecordNotFound will be returned if not exists.
	Update(ctx context.Context, u *model.User) error

	// FindByID returns a model.User if exists with given userID.
	// database.ErrRecordNotFound will be returned if not exists.
	FindByID(ctx context.Context, userID uint) (*model.User, error)

	// FindByName returns a model.User if exists with given username.
	// database.ErrRecordNotFound will be returned if not exists.
	FindByName(ctx context.Context, username string) (*model.User, error)

	// FindByEmail returns a model.User if exists with given email.
	// database.ErrRecordNotFound will be returned if not exists.
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// Follow follows given userID to followerID
	// database.ErrKeyConflict will be returned if already followed.
	// database.ErrFKConstraint will be returned if not exist followerID.
	Follow(ctx context.Context, userID, followerID uint) error

	// IsFollow returns a true if userID follows followerID, otherwise false.
	IsFollow(ctx context.Context, userID, followerID uint) (bool, error)

	// UnFollow unfollows given userID to followerID.
	// database.ErrRecordNotFound will be returned if user does not follow.
	UnFollow(ctx context.Context, userID, followerID uint) error
}

// NewUserDB creates a new UserDB with given gorm.DB
func NewUserDB(db *gorm.DB) UserDB {
	return &userDB{
		db: db,
	}
}

type userDB struct {
	db *gorm.DB
}

func (db *userDB) Save(ctx context.Context, u *model.User) error {
	logger := logging.FromContext(ctx)
	logger.Debugw("UserDB_Save try to save an user", "user", u)

	if err := db.db.WithContext(ctx).Create(u).Error; err != nil {
		logger.Errorw("UserDB_Save failed to save an user", "err", err)
		return database.WrapError(err)
	}
	return nil
}

func (db *userDB) Update(ctx context.Context, u *model.User) error {
	logger := logging.FromContext(ctx)
	logger.Debugw("UserDB_Update try to update an user", "user", u)

	result := db.db.WithContext(ctx).
		Model(new(model.User)).
		Where("user_id = ?", u.ID).
		Updates(model.User{
			Email:     u.Email,
			Name:      u.Name,
			Password:  u.Password,
			Bio:       u.Bio,
			Image:     u.Image,
			UpdatedAt: time.Now(),
			Disabled:  u.Disabled,
		})
	if result.Error != nil {
		logger.Errorw("UserDB_Update failed to update an user", "err", result.Error)
		return database.WrapError(result.Error)
	}
	if result.RowsAffected != 1 {
		logger.Error("UserDB_Update failed to update an user. zero rows affected")
		return database.WrapError(gorm.ErrRecordNotFound)
	}
	return nil
}

func (db *userDB) FindByID(ctx context.Context, userID uint) (*model.User, error) {
	logger := logging.FromContext(ctx)
	logger.Debugw("UserDB_FindByID try to find an user", "userID", userID)

	var u model.User
	if err := db.db.WithContext(ctx).First(&u, "user_id = ?", userID).Error; err != nil {
		logger.Errorw("UserDB_FindByID failed to find an user", "err", err)
		return nil, database.WrapError(err)
	}
	if u.Disabled {
		return nil, database.WrapError(gorm.ErrRecordNotFound)
	}
	return &u, nil
}

func (db *userDB) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	logger := logging.FromContext(ctx)
	logger.Debugw("UserDB_FindByEmail try to find an user", "email", email)

	var u model.User
	if err := db.db.WithContext(ctx).First(&u, "email = ?", email).Error; err != nil {
		logger.Errorw("UserDB_FindByEmail failed to find an user", "err", err)
		return nil, database.WrapError(err)
	}
	if u.Disabled {
		return nil, database.WrapError(gorm.ErrRecordNotFound)
	}
	return &u, nil
}

func (db *userDB) FindByName(ctx context.Context, username string) (*model.User, error) {
	logger := logging.FromContext(ctx)
	logger.Debugw("UserDB_FindByName try to find an user", "username", username)

	var u model.User
	if err := db.db.WithContext(ctx).First(&u, "name = ?", username).Error; err != nil {
		logger.Errorw("UserDB_FindByName failed to find an user", "err", err)
		return nil, database.WrapError(err)
	}
	if u.Disabled {
		return nil, database.WrapError(gorm.ErrRecordNotFound)
	}
	return &u, nil
}

func (db *userDB) Follow(ctx context.Context, userID, followerID uint) error {
	logger := logging.FromContext(ctx)
	logger.Debugw("UserDB_Follow try to insert the following relation", "userID", userID, "followerID", followerID)

	if err := db.db.WithContext(ctx).Create(&model.Follow{
		UserID:   userID,
		FollowID: followerID,
	}).Error; err != nil {
		logger.Errorw("UserDB_Follow failed to insert following relation", "err", err)
		return database.WrapError(err)
	}
	return nil
}

func (db *userDB) IsFollow(ctx context.Context, userID, followerID uint) (bool, error) {
	logger := logging.FromContext(ctx)
	logger.Debugw("UserDB_Follow try to check the following relation", "userID", userID, "followerID", followerID)

	var count int64
	if err := db.db.WithContext(ctx).Model(new(model.Follow)).
		Where("user_id = ? AND follow_id = ?", userID, followerID).
		Count(&count).Error; err != nil {
		logger.Errorw("UserDB_Follow failed to find the following relation", "err", err)
		return false, database.WrapError(err)
	}
	return count == 1, nil
}

func (db *userDB) UnFollow(ctx context.Context, userID, followerID uint) error {
	logger := logging.FromContext(ctx)
	logger.Debugw("UserDB_UnFollow try to delete the following relation", "userID", userID, "followerID", followerID)

	result := db.db.WithContext(ctx).Unscoped().Where("user_id = ? AND follow_id = ?", userID, followerID).Delete(new(model.Follow))
	if result.Error != nil {
		logger.Errorw("UserDB_UnFollow failed to delete", "err", result.Error)
		return database.WrapError(result.Error)
	}
	if result.RowsAffected != 1 {
		logger.Error("UserDB_UnFollow failed to delete the relation. zero rows affected")
		return database.WrapError(gorm.ErrRecordNotFound)
	}
	return nil
}
