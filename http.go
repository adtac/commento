package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"html/template"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "")
}

type resultContainer struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	Comments []Comment `json:"comments"`
}

func (res *resultContainer) render(w http.ResponseWriter) {
	if res == nil {
		res = &resultContainer{false, "Some internal error occurred", nil}
	}

	json, err := json.Marshal(res)
	if err != nil {
		w.Write([]byte(`{"Success":false,"Message":"Internal Server Error"}`))
		return
	}
	w.Write(json)
}

func createCommentHandler(w http.ResponseWriter, r *http.Request) {
	result := &resultContainer{}

	parent, err = strconv.Atoi(r.PostFormValue("parent"))
	if err != nil {
		emit(err)
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

	result := &resultContainer{}

	result.Success = true
	if comments, err := getComments(r.PostFormValue("url")); err != nil {
		emit(err)
		result.Message = "Some internal error occurred."
		result.render(w)
	} else {
		result.Comments = comments
		result.render(w)
	}
}
