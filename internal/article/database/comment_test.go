package database

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/article/model"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/database"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/internal/user/model"
	"math"
	"testing"
	"time"
)

func (s *Suite) TestSaveComment() {
	a := newArticle("article1", "", "", *s.u1, nil)
	s.NoError(s.db.Save(context.TODO(), a))
	c := newComment("comment", *s.u1, *a)

	err := s.db.SaveComment(context.TODO(), c)

	s.NoError(err)
}

func (s *Suite) TestSaveCommentFail() {
	a := newArticle("article1", "", "", *s.u1, nil)
	s.NoError(s.db.Save(context.TODO(), a))

	cases := []struct {
		name    string
		comment *model.Comment
		msg     string
	}{
		{
			name: "empty article commentID",
			comment: &model.Comment{
				Body:     "body",
				AuthorID: s.u1.ID,
			},
			msg: "require article id and author id",
		}, {
			name: "empty author commentID",
			comment: &model.Comment{
				Body:      "body",
				ArticleID: a.ID,
			},
			msg: "require article id and author id",
		}, {
			name: "not exist article commentID",
			comment: &model.Comment{
				Body:      "body",
				ArticleID: math.MaxInt16,
				AuthorID:  s.u1.ID,
			},
			msg: database.ErrFKConstraint.Error(),
		}, {
			name: "not exist author commentID",
			comment: &model.Comment{
				Body:      "body",
				ArticleID: a.ID,
				AuthorID:  math.MaxInt16,
			},
			msg: database.ErrFKConstraint.Error(),
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.db.SaveComment(context.TODO(), tc.comment)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.msg)
		})
	}
}

func (s *Suite) TestFindCommentsByArticleID() {
	a := newArticle("article1", "", "", *s.u1, nil)
	s.NoError(s.db.Save(context.TODO(), a))
	comments := []*model.Comment{
		newComment("comment1", *s.u1, *a),
		newComment("comment2", *s.u2, *a),
	}
	s.NoError(s.originDB.Save(comments).Error)

	find, err := s.db.FindCommentsByArticleID(context.TODO(), a.ID)

	s.NoError(err)
	s.Len(find, 2)
	c1Idx, c2Idx := 0, 1
	if find[0].Body == "comment2" {
		c1Idx, c2Idx = c2Idx, c1Idx
	}
	s.assertComment(comments[0], find[c1Idx])
	s.assertComment(comments[1], find[c2Idx])
}

func (s *Suite) TestFindCommentsByArticleID_Empty() {
	find, err := s.db.FindCommentsByArticleID(context.TODO(), 1)

	s.NoError(err)
	s.Empty(find)
}

func (s *Suite) TestDeleteCommentByID() {
	a := newArticle("article1", "", "", *s.u1, nil)
	s.NoError(s.db.Save(context.TODO(), a))
	c := newComment("comment1", *s.u1, *a)
	s.NoError(s.db.SaveComment(context.TODO(), c))

	err := s.db.DeleteCommentByID(context.TODO(), s.u1, a.ID, c.ID)

	s.NoError(err)
}

func (s *Suite) TestDeleteCommentByID_Fail() {
	a := newArticle("article1", "", "", *s.u1, nil)
	s.NoError(s.db.Save(context.TODO(), a))
	c := newComment("comment1", *s.u1, *a)
	s.NoError(s.db.SaveComment(context.TODO(), c))

	cases := []struct {
		name      string
		user      *userModel.User
		commentID uint
		articleID uint
		msg       string
	}{
		{
			name:      "not exist comment commentID",
			user:      s.u1,
			commentID: math.MaxInt16,
			articleID: a.ID,
			msg:       database.ErrRecordNotFound.Error(),
		}, {
			name:      "mismatch user",
			user:      s.u2,
			commentID: c.ID,
			articleID: a.ID,
			msg:       database.ErrRecordNotFound.Error(),
		}, {
			name:      "mismatch article id",
			user:      s.u1,
			articleID: a.ID + 1,
			commentID: c.ID,
			msg:       database.ErrRecordNotFound.Error(),
		}, {
			name:      "mismatch user and article id",
			user:      s.u2,
			articleID: a.ID + 1,
			commentID: c.ID,
			msg:       database.ErrRecordNotFound.Error(),
		}, {
			name:      "no user provided",
			commentID: math.MaxInt16,
			msg:       database.ErrRecordNotFound.Error(),
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.db.DeleteCommentByID(context.TODO(), tc.user, tc.articleID, tc.commentID)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.msg)
		})
	}
}

func (s *Suite) assertComment(expected, actual *model.Comment) {
	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Body, actual.Body)
	s.WithinDuration(actual.CreatedAt, time.Now(), time.Minute)
	s.WithinDuration(actual.UpdatedAt, time.Now(), time.Minute)
	s.False(actual.DeletedAt.Valid)
	s.Equal(expected.ArticleID, actual.ArticleID)
	s.Equal(expected.AuthorID, actual.AuthorID)
	s.Equal(expected.Author.ID, actual.Author.ID)
	s.Equal(expected.Author.Name, actual.Author.Name)
	s.Equal(expected.Author.Bio, actual.Author.Bio)
	s.Equal(expected.Author.Image, actual.Author.Image)
}

func newComment(body string, author userModel.User, a model.Article) *model.Comment {
	return &model.Comment{
		Body:      body,
		ArticleID: a.ID,
		Author:    author,
		AuthorID:  author.ID,
	}
}
