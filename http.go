package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"html/template"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Foosh")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
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
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	json, err := json.Marshal(res)
	if err != nil {
		w.Write([]byte(`{"Success":false,"Message":"Internal Server Error"}`))
		return
	}
	w.Write(json)
}

func createCommentHandler(w http.ResponseWriter, r *http.Request) {
	result := &resultContainer{}

	fmt.Println("parent:", 	r.PostFormValue("parent"))
	fmt.Println("name:", 	r.PostFormValue("name"))
	fmt.Println("url:", 	r.PostFormValue("url"))
	fmt.Println("comment:", r.PostFormValue("comment"))

	fmt.Println("parent:", 	r.PostFormValue("parent"))
	fmt.Println("name:", 	r.PostFormValue("name"))
	fmt.Println("url:", 	r.PostFormValue("url"))
	fmt.Println("comment:", r.PostFormValue("comment"))

	parent, err := strconv.Atoi(r.PostFormValue("parent"))
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
	fmt.Println("getCommentsHandler")
	// // JwtHandler validates the signature of a JWT and decodes the claims
	// func (mw Middlewares) JwtHandler(next http.Handler) http.Handler {

	// 	fn := func(w http.ResponseWriter, r *http.Request) {

	// 		var jwtString string

	// 		// There was an error getting the JWT from the cookie

	// 		// Getting the JWT from the cookie
	// 		if jwtCookie, err := r.Cookie("TSIJwt"); err == nil {
	// 			jwtString = jwtCookie.Value
	// 			// Could not get the JWT from the cookie
	// 		} else

	fmt.Println("here1")
	// w.Header().Set("Access-Control-Allow-Origin", "GET")
	// origin := r.Header.Get("Origin")
	// fmt.Println("Origin:", origin)
	// w.Header().Set("Access-Control-Allow-Credentials", "true")
	// w.Header().Set("Access-Control-Allow-Origin", origin)

	result := &resultContainer{}

	fmt.Println("here1")

	for _, cookie := range r.Cookies() {
		fmt.Fprint(w, cookie.Name)
	}


	var userID string
	if userIDCookie, err := r.Cookie("userID"); err != nil {
		emit(err)
		result.Message = "Cookie was not set"
		result.render(w)
		return
	} else {
		userID = userIDCookie.Value
	}
	fmt.Println("here1")
	comments := []Comment{}
	result.Success = true
	// comments, err = getComments(r.PostFormValue("url"))
	fmt.Println("here1")
	comments, err := getComments(userID)
	if err != nil {
		emit(err)
	}
	result.Comments = comments
	fmt.Println("result:", result)

	result.render(w)
}
