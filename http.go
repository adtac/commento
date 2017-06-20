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
	Status   int       `json:"-"`
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	Comments []Comment `json:"comments,omitempty"`
	Count    int       `json:"count"`
}

func (res *resultContainer) render(w http.ResponseWriter) {
	if res == nil {
		res = &resultContainer{
			Status:   http.StatusInternalServerError,
			Success:  false,
			Message:  "Some internal error occurred",
			Comments: nil,
		}
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if res.Status == 0 {
		res.Status = 200
	}
	w.WriteHeader(res.Status)

	json, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Success":false,"Message":"Internal Server Error"}`))
		return
	}
	w.Write(json)
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	result := &resultContainer{}
	if r.Method != "POST" {
		result.Status = http.StatusMethodNotAllowed
		result.Message = "This request must be a POST request."
		result.render(w)
		return
	}

	parent, err := strconv.Atoi(r.PostFormValue("parent"))
	if err != nil {
		Emit(err)
		result.Status = http.StatusBadRequest
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
		result.Status = http.StatusInternalServerError
		result.Message = "Some internal error occurred."
		result.render(w)
		return
	}

	result.Message = "Comment successfully created"
	result.Success = true

	var count int
	var countErr error
	count, countErr = countComments(r.PostFormValue("url"))
	if countErr != nil {
		result.Message = "Comment successfully created, however could not retrieve updated comment count."
		Emit(err)
	} else {
		result.Count = count
	}

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
	result.Count = len(comments)
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
	result.Count = count
	result.render(w)
}
