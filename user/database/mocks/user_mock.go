// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/zacscoding/echo-gorm-realworld-app/user/model"
)

// UserDB is an autogenerated mock type for the UserDB type
type UserDB struct {
	mock.Mock
}

// FindByEmail provides a mock function with given fields: ctx, email
func (_m *UserDB) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	ret := _m.Called(ctx, email)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByID provides a mock function with given fields: ctx, userID
func (_m *UserDB) FindByID(ctx context.Context, userID uint) (*model.User, error) {
	ret := _m.Called(ctx, userID)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(context.Context, uint) *model.User); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByName provides a mock function with given fields: ctx, username
func (_m *UserDB) FindByName(ctx context.Context, username string) (*model.User, error) {
	ret := _m.Called(ctx, username)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Follow provides a mock function with given fields: ctx, userID, followerID
func (_m *UserDB) Follow(ctx context.Context, userID uint, followerID uint) error {
	ret := _m.Called(ctx, userID, followerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) error); ok {
		r0 = rf(ctx, userID, followerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsFollow provides a mock function with given fields: ctx, userID, followerID
func (_m *UserDB) IsFollow(ctx context.Context, userID uint, followerID uint) (bool, error) {
	ret := _m.Called(ctx, userID, followerID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) bool); ok {
		r0 = rf(ctx, userID, followerID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint, uint) error); ok {
		r1 = rf(ctx, userID, followerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, u
func (_m *UserDB) Save(ctx context.Context, u *model.User) error {
	ret := _m.Called(ctx, u)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) error); ok {
		r0 = rf(ctx, u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnFollow provides a mock function with given fields: ctx, userID, followerID
func (_m *UserDB) UnFollow(ctx context.Context, userID uint, followerID uint) error {
	ret := _m.Called(ctx, userID, followerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) error); ok {
		r0 = rf(ctx, userID, followerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, u
func (_m *UserDB) Update(ctx context.Context, u *model.User) error {
	ret := _m.Called(ctx, u)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) error); ok {
		r0 = rf(ctx, u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
