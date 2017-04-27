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

func createCommentHandler(w http.ResponseWriter, r *http.Request) {
	type resultContainer struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	var parent int
	var err error
	var name string

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

	name = alphaNumericOnly(r.PostFormValue("name"))

	err = createComment(r.PostFormValue("url"), name, r.PostFormValue("comment"), parent)
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
		Success  bool      `json:"success"`
		Message  string    `json:"message"`
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
