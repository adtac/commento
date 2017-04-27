package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "")
}

type resultContainer struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	Comments []Comment `json:"comments,omitempty"`
}

func (res *resultContainer) render(w http.ResponseWriter) {
	if res == nil {
		res = &resultContainer{false, "Some internal error occurred", nil}
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json, err := json.Marshal(res)
	if err != nil {
		w.Write([]byte(`{"Success":false,"Message":"Internal Server Error"}`))
		return
	}
	w.Write(json)
}

func createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var parent int
	var err error
	var name string

	result := &resultContainer{}
	if r.Method != "POST" {
		result.Message = "This request must be a POST request."
		result.render(w)
		return
	}

	parent, err = strconv.Atoi(r.PostFormValue("parent"))
	if err != nil {
		emit(err)
		result.Message = "Invalid parent ID."
		result.render(w)
		return
	}

	name = alphaNumericOnly(r.PostFormValue("name"))

	if r.PostFormValue("gotcha") != "" {
		// If a value has been set, we just silently ignore the submission and return
		// a success message. This prevents spammers from cottoning-on that the submission
		// did not work.
		result.render(w)
		return
	}

	err = createComment(r.PostFormValue("url"), name, r.PostFormValue("comment"), parent)
	if err != nil {
		emit(err)
		result.Message = "Some internal error occurred."
		result.render(w)
		return
	}
	result.Success = true
	result.Message = "Comment successfully created"
	result.render(w)
}

func getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	comments := []Comment{}
	var err error

	result := &resultContainer{Success: true}
	if r.Method != "POST" {
		result.Success = false
		result.Message = "This request must be a POST request."
		result.render(w)
		return
	}

	comments, err = getComments(r.PostFormValue("url"))
	if err != nil {
		emit(err)
	}
	result.Comments = comments
	result.render(w)
}
