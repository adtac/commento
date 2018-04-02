package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// resultContainer stores the results of a request
type resultContainer struct {
	Status   int       `json:"-"`
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	Comments []Comment `json:"comments,omitempty"`
}

// IndexHandler handles GET requests to the root path '/'
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "")
}

// render writes a resultContainer to a response stream
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

// CreateCommentHandler handles the '/create' endpoint that is used to create a
// new comment. It requires the following POST request body values:
//	   - name: the name of the comment author
//	   - parent: ID of the parent comment
//     - comment: the comment text itself
//     - url: the URL associated with this comment
func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	result := &resultContainer{}
	var err error

	if r.Method != "POST" {
		result.Status = http.StatusMethodNotAllowed
		result.Message = errorList["err.request.method.invalid"].Error()
		result.render(w)
		return
	}

	requiredFields := []string{"name", "parent", "comment", "url"}
	for _, field := range requiredFields {
		if strings.TrimSpace(r.PostFormValue(field)) == "" {
			result.Status = http.StatusBadRequest
			result.Message = errorList["err.request.field.missing"].Error()
			result.render(w)
			return
		}
	}

	if r.PostFormValue("gotcha") != "" {
		result.Success = true
		result.Message = "Comment successfully created"
		result.render(w)
		return
	}

	comment := Comment{}

	comment.Name = template.HTMLEscapeString(r.PostFormValue("name"))

	comment.Comment = template.HTMLEscapeString(r.PostFormValue("comment"))

	comment.URL = r.PostFormValue("url")

	comment.Parent, err = strconv.Atoi(r.PostFormValue("parent"))
	if err != nil || comment.Parent < -1 {
		result.Status = http.StatusBadRequest
		result.Message = errorList["err.request.field.invalid"].Error()
		result.render(w)
		return
	}

	if isSpam := checkSpam(r, comment.URL, comment.Name, comment.Comment); isSpam {
		// Silently fail. Don't tell the spammer we detected their comment.
		result.Success = true
		result.Message = "Comment successfully created"
		result.render(w)
		return
	}

	err = db.CreateComment(&comment)
	if err != nil {
		result.Status = http.StatusInternalServerError
		result.Message = errorList["err.internal"].Error()
		result.render(w)
		return
	}

	result.Success = true
	result.Message = "Comment successfully created"
	result.render(w)
}

// GetCommentsHandler handles the '/get' endpoint that is used to retrieve
// all the comments for a particular URL. It takes one value:
//     - url: the URL associated with this comment
func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	result := &resultContainer{}
	comments := []Comment{}
	var err error

	if r.Method != "POST" {
		result.Status = http.StatusMethodNotAllowed
		result.Message = errorList["err.request.method.invalid"].Error()
		result.render(w)
		return
	}

	requiredFields := []string{"url"}
	for _, field := range requiredFields {
		if strings.TrimSpace(r.PostFormValue(field)) == "" {
			result.Status = http.StatusBadRequest
			result.Message = errorList["err.request.field.missing"].Error()
			result.render(w)
			return
		}
	}

	comments, err = db.GetComments(r.PostFormValue("url"))
	if err != nil {
		result.Status = http.StatusInternalServerError
		result.Message = errorList["err.internal"].Error()
		result.render(w)
		return
	}

	result.Success = true
	result.Comments = comments
	result.render(w)
}
