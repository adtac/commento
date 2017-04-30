package main

import (
	"fmt"
)

func emit(err error) {
	fmt.Println(err)
}

func die(err error) {
	logger.Panic(err)
}
