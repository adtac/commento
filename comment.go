package main

import (
	"time"
)

// A Comment represents a JSON structure for comments on commento
type Comment struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	Comment   string    `json:"comment"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Parent    int       `json:"parent"`
}

// A CommentService handles the CRUD operations on comments
type CommentService interface {
	CreateComment(comment *Comment) error
	DeleteComment(commentId int) error
	GetComments(url string) ([]Comment, error)
}
