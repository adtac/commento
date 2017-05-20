package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type handlerTestList struct {
	handler      http.Handler
	method       string
	path         string
	reqValues    map[string]string
	code         int
	bodycontains string
}

func (hl handlerTestList) Test(t *testing.T) {
	req, err := http.NewRequest(hl.method, hl.path, nil)
	if err != nil {
		t.Fatal(err)
	}
	// load test form values, if any
	func() {
		if len(hl.reqValues) == 0 {
			return
		}
		form := url.Values{}
		for k, v := range hl.reqValues {
			form.Add(k, v)
		}
		req, err = http.NewRequest(hl.method, hl.path, strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}()
	rr := httptest.NewRecorder()
	hl.handler.ServeHTTP(rr, req)
	resp := rr.Result()
	if resp.StatusCode != hl.code {
		t.Errorf("handler returned bad status code, got %d, want %d", resp.StatusCode, hl.code)
	}
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("handler response header Access-Control-Allow-Origin not set to *")
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(dat), hl.bodycontains) {
		t.Errorf("handler missing content '%s', required in body", hl.bodycontains)
	}
}

func TestIndexHandler(t *testing.T) {
	handler := http.HandlerFunc(IndexHandler)
	l := handlerTestList{handler, "GET", "/", nil, http.StatusOK, ""}
	l.Test(t)
}
func TestCreateCommentHandler(t *testing.T) {
	withBadParent := map[string]string{
		"parent": "missing",
	}
	withGotcha := map[string]string{
		"parent": "1",
		"gotcha": "got me",
	}
	withMissingFields := map[string]string{
		"parent":  "1",
		"comment": "hello commentors!",
	}
	withBlankFields := map[string]string{
		"parent":  "1",
		"url":     "",
		"name":    "",
		"comment": "hello commentors!",
	}
	withSpaceInFields := map[string]string{
		"parent":  "1",
		"url":     " ",
		"name":    " ",
		"comment": "hello commentors!",
	}
	withOKComment := map[string]string{
		"parent":  "2",
		"url":     "/",
		"name":    "Tester",
		"comment": "hello commentors!",
	}
	handler := http.HandlerFunc(CreateCommentHandler)
	list := []handlerTestList{
		{handler, "GET", "/", nil, http.StatusMethodNotAllowed, "must be a POST request"},
		{handler, "POST", "/", nil, http.StatusBadRequest, "Invalid parent ID"},
		{handler, "POST", "/", withBadParent, http.StatusBadRequest, "Invalid parent ID"},
		{handler, "POST", "/", withGotcha, http.StatusOK, ""},
		{handler, "POST", "/", withMissingFields, http.StatusBadRequest, "invalid comment"},
		{handler, "POST", "/", withBlankFields, http.StatusBadRequest, "invalid comment"},
		{handler, "POST", "/", withSpaceInFields, http.StatusBadRequest, "invalid comment"},
		{handler, "POST", "/", withOKComment, http.StatusOK, "successfully created"},
	}
	for _, l := range list {
		l.Test(t)
	}
}

func TestGetCommentsHandler(t *testing.T) {
	withOKComment := map[string]string{
		"parent":  "3",
		"url":     "/",
		"name":    "Tester",
		"comment": "hello commentors!",
	}
	postH := http.HandlerFunc(CreateCommentHandler)
	posted := handlerTestList{postH, "POST", "/", withOKComment, http.StatusOK, "successfully created"}
	posted.Test(t) // add comment, in order not to fetch blanks

	withOKURL := map[string]string{
		"url": "/",
	}
	handler := http.HandlerFunc(GetCommentsHandler)
	list := []handlerTestList{
		{handler, "POST", "/", withOKURL, http.StatusOK, "hello commentors"},
	}
	for _, l := range list {
		l.Test(t)
	}
}
