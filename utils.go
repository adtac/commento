package main

import (
	"fmt"

	"github.com/op/go-logging"
)

// Logger specifies the logger to use
var Logger = logging.MustGetLogger("commento")

// Emit prints an error message to
func Emit(err error) {
	fmt.Printf("%s", fmt.Errorf("%s", err))
}

// Die logs an error and panics
func Die(err error) {
	Logger.Panic(fmt.Errorf("%s", err))
}
