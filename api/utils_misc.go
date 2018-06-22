package main

import (
	"bytes"
	"fmt"
	"os"
)

func concat(a bytes.Buffer, b bytes.Buffer) []byte {
	return append(a.Bytes(), b.Bytes()...)
}

func exitIfError(err error) {
	if err != nil {
		fmt.Printf("fatal error: %v\n", err)
		os.Exit(1)
	}
}
