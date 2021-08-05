package types

import (
	articlemodel "github.com/zacscoding/echo-gorm-realworld-app/article/model"
)

// CommentResponse represents a single comment response.
type CommentResponse struct {
	Comment *Comment `json:"comment"`
}

// ToCommentResponse converts given c to CommentResponse.
func ToCommentResponse(c *articlemodel.Comment) *CommentResponse {
	return &CommentResponse{
		Comment: toComment(c),
	}
}

// CommentsResponse represents multiple comments response.
type CommentsResponse struct {
	Comments []*Comment `json:"comments"`
}

// ToCommentsResponse converts given comments to CommentsResponse.
func ToCommentsResponse(comments []*articlemodel.Comment) *CommentsResponse {
	res := new(CommentsResponse)
	res.Comments = make([]*Comment, len(comments))
	for i, c := range comments {
		res.Comments[i] = toComment(c)
	}
	return res
}

type Comment struct {
	ID        uint     `json:"id"`
	CreatedAt JSONTime `json:"createdAt"`
	UpdatedAt JSONTime `json:"updatedAt"`
	Body      string   `json:"body"`
	Author    Author   `json:"author"`
}

func toComment(c *articlemodel.Comment) *Comment {
	return &Comment{
		ID:        c.ID,
		CreatedAt: JSONTime(c.CreatedAt),
		UpdatedAt: JSONTime(c.UpdatedAt),
		Body:      c.Body,
		Author:    toAuthor(&c.Author),
	}
}
