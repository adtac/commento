package main

import (
	"bytes"
	"os"
)

func concat(a bytes.Buffer, b bytes.Buffer) []byte {
	return append(a.Bytes(), b.Bytes()...)
}

func exitIfError(err error) {
	if err != nil {
		os.Exit(1)
	}
}
