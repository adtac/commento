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
