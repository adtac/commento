package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "")
}

type resultContainer struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	Comments []Comment `json:"comments,omitempty"`
	Count    *int       `json:"count,omitempty"`
}

func (res *resultContainer) render(w http.ResponseWriter) {
	if res == nil {
		res = &resultContainer{false, "Some internal error occurred", nil, nil}
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json, err := json.Marshal(res)
	if err != nil {
		w.Write([]byte(`{"Success":false,"Message":"Internal Server Error"}`))
		return
	}
	w.Write(json)
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	result := &resultContainer{}
	if r.Method != "POST" {
		result.Message = "This request must be a POST request."
		result.render(w)
		return
	}

	parent, err := strconv.Atoi(r.PostFormValue("parent"))
	if err != nil {
		Emit(err)
		result.Message = "Invalid parent ID."
		result.render(w)
		return
	}

	name := template.HTMLEscapeString(r.PostFormValue("name"))
	comment := template.HTMLEscapeString(r.PostFormValue("comment"))

	if r.PostFormValue("gotcha") != "" {
		// If a value has been set, we just silently ignore the submission and return
		// a success message. This prevents spammers from cottoning-on that the submission
		// did not work.
		result.render(w)
		return
	}

	err = createComment(r.PostFormValue("url"), name, comment, parent)
	if err != nil {
		Emit(err)
		result.Message = "Some internal error occurred."
		result.render(w)
		return
	}
	result.Success = true
	result.Message = "Comment successfully created"
	result.render(w)
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	comments := []Comment{}
	var err error

	result := &resultContainer{Success: true}
	comments, err = getComments(r.PostFormValue("url"))
	if err != nil {
		Emit(err)
	}
	result.Comments = comments
	count := int(len(comments))
	result.Count = &count
	result.render(w)
}

func CountCommentsHandler(w http.ResponseWriter, r *http.Request) {
	result := &resultContainer{}

	url := r.URL.Query().Get("url")
	if url == "" {
		result.Success = false
		result.Message = "No URL provided."
		result.render(w)
		return
	}

	var count int
	var err error
	count, err = countComments(url)
	if err != nil {
		result.Success = false
		result.Message = "Some internal error occurred."
		Emit(err)
		return
	}

	result.Success = true
	count = int(count)
	result.Count = &count
	result.render(w)
}
