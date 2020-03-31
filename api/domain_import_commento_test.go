package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestImportCommento(t *testing.T) {
	failTestOnError(t, setupTestEnv())

	// Create JSON data
	data := commentoExportV1{
		Version: 1,
		Comments: []comment{
			{
				CommentHex:   "5a349182b3b8e25107ab2b12e514f40fe0b69160a334019491d7c204aff6fdc2",
				Domain:       "localhost:1313",
				Path:         "/post/first-post/",
				CommenterHex: "anonymous",
				Markdown:     "This is a reply!",
				Html:         "",
				ParentHex:    "7ed60b1227f6c4850258a2ac0304e1936770117d6f3a379655f775c46b9f13cd",
				Score:        0,
				State:        "approved",
				CreationDate: timeParse(t, "2020-01-27T14:08:44.061525Z"),
				Direction:    0,
				Deleted:      false,
			},
			{
				CommentHex:   "7ed60b1227f6c4850258a2ac0304e1936770117d6f3a379655f775c46b9f13cd",
				Domain:       "localhost:1313",
				Path:         "/post/first-post/",
				CommenterHex: "anonymous",
				Markdown:     "This is a comment!",
				Html:         "",
				ParentHex:    "root",
				Score:        0,
				State:        "approved",
				CreationDate: timeParse(t, "2020-01-27T14:07:49.244432Z"),
				Direction:    0,
				Deleted:      false,
			},
			{
				CommentHex:   "a7c84f251b5a09d5b65e902cbe90633646437acefa3a52b761fee94002ac54c7",
				Domain:       "localhost:1313",
				Path:         "/post/first-post/",
				CommenterHex: "4629a8216538b73987597d66f266c1a1801b0451f99cf066e7122aa104ef3b07",
				Markdown:     "This is a test comment, bar foo\n\n#Here is something big\n\n```\nhere code();\n```",
				Html:         "",
				ParentHex:    "root",
				Score:        0,
				State:        "approved",
				CreationDate: timeParse(t, "2020-01-27T14:20:21.101653Z"),
				Direction:    0,
				Deleted:      false,
			},
		},
		Commenters: []commenter{
			{
				CommenterHex: "4629a8216538b73987597d66f266c1a1801b0451f99cf066e7122aa104ef3b07",
				Email:        "john@doe.com",
				Name:         "John Doe",
				Link:         "https://john.doe",
				Photo:        "undefined",
				Provider:     "commento",
				JoinDate:     timeParse(t, "2020-01-27T14:17:59.298737Z"),
				IsModerator:  false,
			},
		},
	}

	// Create listener with random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Errorf("couldn't create listener: %v", err)
		return
	}
	defer func() {
		_ = listener.Close()
	}()
	port := listener.Addr().(*net.TCPAddr).Port

	// Launch http server serving commento json gzipped data
	go func() {
		http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gzipper := gzip.NewWriter(w)
			defer func() {
				_ = gzipper.Close()
			}()
			encoder := json.NewEncoder(gzipper)
			if err := encoder.Encode(data); err != nil {
				t.Errorf("couldn't write data: %v", err)
			}
		}))
	}()
	url := fmt.Sprintf("http://127.0.0.1:%d", port)

	domainNew("temp-owner-hex", "Example", "example.com")

	n, err := domainImportCommento("example.com", url)
	if err != nil {
		t.Errorf("unexpected error importing comments: %v", err)
		return
	}
	if n != len(data.Comments) {
		t.Errorf("imported comments missmatch (got %d, want %d)", n, len(data.Comments))
	}
}

func timeParse(t *testing.T, s string) time.Time {
	time, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t.Errorf("couldn't parse time: %v", err)
	}
	return time
}
