package main

import (
	"bytes"
	"fmt"
	"os"
)

func concat(a bytes.Buffer, b bytes.Buffer) []byte {
	return append(a.Bytes(), b.Bytes()...)
}

func nameFromEmail(email string) string {
	for i, c := range email {
		if c == '@' {
			return email[:i]
		}
	}

	return email
}

func exitIfError(err error) {
	if err != nil {
		fmt.Printf("fatal error: %v\n", err)
		os.Exit(1)
	}
}
