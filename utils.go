package main

import (
	"fmt"
)

func emit(err error) {
	fmt.Printf("%s", fmt.Errorf("%s", err))
}

func die(err error) {
	logger.Panic(fmt.Errorf("%s", err))
}

func alphaOnly(s string) bool {
	alpha := "abcdefghijklmnopqrstuvwxyz0123456789"
	for _, char := range s {
		if !strings.Contains(alpha, strings.ToLower(string(char))) {
			return false
		}
	}
	return true
}
