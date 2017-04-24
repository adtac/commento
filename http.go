package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "")
}

func createCommentHandler(w http.ResponseWriter, r *http.Request) {
	type resultContainer struct {
		Success bool `json:"success"`
		Message string `json:"message"`
	}

	var parent int
	var err error

	result := resultContainer{Success: true}
	if r.Method != "POST" {
		result.Success = false
		result.Message = "This request must be a POST request."
		goto end
	}

	parent, err = strconv.Atoi(r.PostFormValue("parent"))
	if err != nil {
		emit(err)
		result.Success = false
		result.Message = "Invalid parent ID."
		goto end
	}

	err = createComment(r.PostFormValue("url"), r.PostFormValue("name"), r.PostFormValue("comment"), parent)
	if err != nil {
		emit(err)
		result.Success = false
		result.Message = "Some internal error occurred."
		goto end
	}

	end:
    json, _ := json.Marshal(result)
	w.Header().Set("Access-Control-Allow-Origin", "*")
    fmt.Fprintf(w, "%s", string(json))
}

func getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	type resultContainer struct {
		Success bool       `json:"success"`
		Message string     `json:"message"`
		Comments []Comment `json:"comments"`
	}

	comments := []Comment{}
	var err error

	result := resultContainer{Success: true}
	if r.Method != "POST" {
		result.Success = false
		result.Message = "This request must be a POST request."
		goto end
	}

	comments, err = getComments(r.PostFormValue("url"))
	if err != nil {
		emit(err)
	}

	end:
	result.Comments = comments
    json, _ := json.Marshal(result)
	w.Header().Set("Access-Control-Allow-Origin", "*")
    fmt.Fprintf(w, "%s", string(json))
}
