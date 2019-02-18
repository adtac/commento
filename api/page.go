package main

import ()

type page struct {
	Domain           string `json:"domain"`
	Path             string `json:"path"`
	IsLocked         bool   `json:"isLocked"`
	CommentCount     int    `json:"commentCount"`
	StickyCommentHex string `json:"stickyCommentHex"`
	Title            string `json:"title"`
}
