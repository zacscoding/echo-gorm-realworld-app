// Code generated by mockery (devel). DO NOT EDIT.

package mocks

import (
	context "context"

	articlemodel "github.com/zacscoding/echo-gorm-realworld-app/article/model"

	mock "github.com/stretchr/testify/mock"

	model "github.com/zacscoding/echo-gorm-realworld-app/user/model"
)

// ArticleDB is an autogenerated mock type for the ArticleDB type
type ArticleDB struct {
	mock.Mock
}

// DeleteBySlug provides a mock function with given fields: ctx, user, slug
func (_m *ArticleDB) DeleteBySlug(ctx context.Context, user *model.User, slug string) error {
	ret := _m.Called(ctx, user, slug)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User, string) error); ok {
		r0 = rf(ctx, user, slug)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteCommentByID provides a mock function with given fields: ctx, user, articleID, commentID
func (_m *ArticleDB) DeleteCommentByID(ctx context.Context, user *model.User, articleID uint, commentID uint) error {
	ret := _m.Called(ctx, user, articleID, commentID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User, uint, uint) error); ok {
		r0 = rf(ctx, user, articleID, commentID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindArticlesByAuthors provides a mock function with given fields: ctx, user, authors, offset, limit
func (_m *ArticleDB) FindArticlesByAuthors(ctx context.Context, user *model.User, authors []uint, offset int, limit int) (*articlemodel.Articles, error) {
	ret := _m.Called(ctx, user, authors, offset, limit)

	var r0 *articlemodel.Articles
	if rf, ok := ret.Get(0).(func(context.Context, *model.User, []uint, int, int) *articlemodel.Articles); ok {
		r0 = rf(ctx, user, authors, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*articlemodel.Articles)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.User, []uint, int, int) error); ok {
		r1 = rf(ctx, user, authors, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindArticlesByQuery provides a mock function with given fields: ctx, user, query, offset, limit
func (_m *ArticleDB) FindArticlesByQuery(ctx context.Context, user *model.User, query articlemodel.ArticleQuery, offset int, limit int) (*articlemodel.Articles, error) {
	ret := _m.Called(ctx, user, query, offset, limit)

	var r0 *articlemodel.Articles
	if rf, ok := ret.Get(0).(func(context.Context, *model.User, articlemodel.ArticleQuery, int, int) *articlemodel.Articles); ok {
		r0 = rf(ctx, user, query, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*articlemodel.Articles)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.User, articlemodel.ArticleQuery, int, int) error); ok {
		r1 = rf(ctx, user, query, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindBySlug provides a mock function with given fields: ctx, user, slug
func (_m *ArticleDB) FindBySlug(ctx context.Context, user *model.User, slug string) (*articlemodel.Article, error) {
	ret := _m.Called(ctx, user, slug)

	var r0 *articlemodel.Article
	if rf, ok := ret.Get(0).(func(context.Context, *model.User, string) *articlemodel.Article); ok {
		r0 = rf(ctx, user, slug)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*articlemodel.Article)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.User, string) error); ok {
		r1 = rf(ctx, user, slug)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindCommentsByArticleID provides a mock function with given fields: ctx, articleID
func (_m *ArticleDB) FindCommentsByArticleID(ctx context.Context, articleID uint) ([]*articlemodel.Comment, error) {
	ret := _m.Called(ctx, articleID)

	var r0 []*articlemodel.Comment
	if rf, ok := ret.Get(0).(func(context.Context, uint) []*articlemodel.Comment); ok {
		r0 = rf(ctx, articleID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*articlemodel.Comment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, articleID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, a
func (_m *ArticleDB) Save(ctx context.Context, a *articlemodel.Article) error {
	ret := _m.Called(ctx, a)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *articlemodel.Article) error); ok {
		r0 = rf(ctx, a)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveComment provides a mock function with given fields: ctx, c
func (_m *ArticleDB) SaveComment(ctx context.Context, c *articlemodel.Comment) error {
	ret := _m.Called(ctx, c)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *articlemodel.Comment) error); ok {
		r0 = rf(ctx, c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, user, a
func (_m *ArticleDB) Update(ctx context.Context, user *model.User, a *articlemodel.Article) error {
	ret := _m.Called(ctx, user, a)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User, *articlemodel.Article) error); ok {
		r0 = rf(ctx, user, a)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
