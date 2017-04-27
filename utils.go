package main

import (
	"fmt"
	"strings"
)

func emit(err error) {
	fmt.Printf("%s", fmt.Errorf("%s", err))
}

func die(err error) {
	logger.Panic(fmt.Errorf("%s", err))
}

func alphaNumericOnly(s string) string {
	alpha := "abcdefghijklmnopqrstuvwxyz0123456789 _-"
	output := ""
	for _, char := range s {
		if strings.Contains(alpha, strings.ToLower(string(char))) {
			output += string(char)
		}
	}
	return output
}
