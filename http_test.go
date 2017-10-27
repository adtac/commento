package main

import (
	"os"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type testCaseHTTP struct {
	name         string
	handler      func(w http.ResponseWriter, r *http.Request)
	method       string
	path         string
	postBody     map[string]string
	statusCode   int
	bodyContains string
}

func TestHTTP(t *testing.T) {
	tests := []func(*testing.T){
		testIndexHandler,
		testCreateCommentHandler,
		testGetCommentsHandler,
	}

	cleanupHTTPTest(t)
	for _, test := range tests {
		setupHTTPTest(t)
		test(t)
		cleanupHTTPTest(t)
	}
}

func setupHTTPTest(t *testing.T) {
	err := LoadDatabase("sqlite:file=sqlite3.db")
	if err != nil {
		t.Errorf("Unable to create test sqlite DB: %v", err)
	}
}

func cleanupHTTPTest(t *testing.T) {
	err := os.Remove("sqlite3.db")
	if err != nil && !os.IsNotExist(err) {
		t.Logf("Unable to remove the test sqlite file: %v", err)
	}
}

func (tc testCaseHTTP) Test(t *testing.T, funcName string) {
	form := url.Values{}
	for key, value := range tc.postBody {
		form.Add(key, value)
	}

	request, err := http.NewRequest(tc.method, tc.path, strings.NewReader(form.Encode()))
	if err != nil {
		t.Errorf("%s: %s: Cannot prepare the message: %v", funcName, tc.name, err)
		return
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	tc.handler(recorder, request)
	response := recorder.Result()

	if tc.statusCode != response.StatusCode {
		t.Errorf("%s: %s: Incorrect status code: expected %d, got %d\n", funcName, tc.name, tc.statusCode, response.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("%s: %s: Cannot read the response body: %v\n", funcName, tc.name, err)
		return
	}

	if !strings.Contains(string(body), tc.bodyContains) {
		t.Errorf("%s: %s: Response body does not contain '%s'; instead got '%s'\n", funcName, tc.name, tc.bodyContains, string(body))
		return
	}
}

func runTestCasesHTTP(t *testing.T, funcName string, testCasesHTTP []testCaseHTTP) {
	for _, tc := range testCasesHTTP {
		tc.Test(t, funcName)
	}
}

func testIndexHandler(t *testing.T) {
	testCasesHTTP := []testCaseHTTP{
		testCaseHTTP{
			"Index should return 200 OK",
			IndexHandler, "GET", "/",
			nil,
			http.StatusOK, "",
		},
	}

	runTestCasesHTTP(t, "IndexHandler", testCasesHTTP)
}

func testCreateCommentHandler(t *testing.T) {
	testCasesHTTP := []testCaseHTTP{
		testCaseHTTP{
			"non-POST requests should be rejected",
			CreateCommentHandler, "GET", "/create",
			nil,
			http.StatusMethodNotAllowed, errorList["err.request.method.invalid"].Error(),
		},

		testCaseHTTP{
			"POST request with no body should be rejected",
			CreateCommentHandler, "POST", "/create",
			nil,
			http.StatusBadRequest, errorList["err.request.field.missing"].Error(),
		},

		// this case is sufficient to test all missing fields scenarios
		testCaseHTTP{
			"Empty 'parent' should be rejected",
			CreateCommentHandler, "POST", "/create",
			map[string]string{
				"parent": "",
				"name": "name",
				"url": "url",
				"comment": "comment",
			},
			http.StatusBadRequest, errorList["err.request.field.missing"].Error(),
		},

		testCaseHTTP{
			"Whitespace fields should be rejected",
			CreateCommentHandler, "POST", "/create",
			map[string]string{
				"parent": "1",
				"name": " \t",
				"url": "url",
				"comment": "comment",
			},
			http.StatusBadRequest, errorList["err.request.field.missing"].Error(),
		},

		testCaseHTTP{
			"Non-integral 'parent' should be rejected",
			CreateCommentHandler, "POST", "/create",
			map[string]string{
				"parent": "abcd",
				"name": "name",
				"url": "url",
				"comment": "comment",
			},
			http.StatusBadRequest, errorList["err.request.field.invalid"].Error(),
		},

		testCaseHTTP{
			"Negative 'parent' (except -1) should be rejected",
			CreateCommentHandler, "POST", "/create",
			map[string]string{
				"parent": "-12",
				"name": "name",
				"url": "url",
				"comment": "comment",
			},
			http.StatusBadRequest, errorList["err.request.field.invalid"].Error(),
		},

		testCaseHTTP{
			"A good message should be accepted",
			CreateCommentHandler, "POST", "/create",
			map[string]string{
				"parent": "2",
				"name": "name",
				"url": "url",
				"comment": "comment",
			},
			http.StatusOK, "",
		},

		testCaseHTTP{
			"Comments with honeypot filled should be rejected",
			CreateCommentHandler, "POST", "/create",
			map[string]string{
				"parent": "2",
				"name": "name",
				"url": "url2",
				"comment": "comment from a spammer",
				"gotcha": "random autofilled text",
			},
			http.StatusOK, "",
		},

		testCaseHTTP{
			"The previous comment should not be present",
			GetCommentsHandler, "POST", "/get",
			map[string]string{
				"url": "url2",
			},
			http.StatusOK, `{"success":true,"message":""}`,
		},
	}

	runTestCasesHTTP(t, "CreateCommentHandler", testCasesHTTP)
}

func testGetCommentsHandler(t *testing.T) {
	testCasesHTTP := []testCaseHTTP{
		testCaseHTTP{
			"non-POST requests should be rejected",
			GetCommentsHandler, "GET", "/get",
			nil,
			http.StatusMethodNotAllowed, errorList["err.request.method.invalid"].Error(),
		},

		testCaseHTTP{
			"POST request with no body should be rejected",
			GetCommentsHandler, "POST", "/get",
			nil,
			http.StatusBadRequest, errorList["err.request.field.missing"].Error(),
		},

		testCaseHTTP{
			"Empty 'url' should be rejected",
			GetCommentsHandler, "POST", "/get",
			map[string]string{
				"url": "",
			},
			http.StatusBadRequest, errorList["err.request.field.missing"].Error(),
		},

		testCaseHTTP{
			"Whitespace fields should be rejected",
			CreateCommentHandler, "POST", "/create",
			map[string]string{
				"url": " \t",
			},
			http.StatusBadRequest, errorList["err.request.field.missing"].Error(),
		},

		testCaseHTTP{
			"Seed comment for retrieval",
			CreateCommentHandler, "POST", "/create",
			map[string]string{
				"parent": "2",
				"name": "name",
				"url": "url1",
				"comment": "some unique comment",
			},
			http.StatusOK, "",
		},

		testCaseHTTP{
			"Comment retrieval should return all comments",
			GetCommentsHandler, "POST", "/get",
			map[string]string{
				"url": "url1",
			},
			http.StatusOK, "some unique comment",
		},

		testCaseHTTP{
			"Retrieval for a non-existant URL should return no comments",
			GetCommentsHandler, "POST", "/get",
			map[string]string{
				"url": "url2",
			},
			http.StatusOK, `{"success":true,"message":""}`,
		},
	}

	runTestCasesHTTP(t, "GetCommentHandler", testCasesHTTP)
}
