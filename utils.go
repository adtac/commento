package main

import (
	"fmt"

	"github.com/op/go-logging"
)

var Logger = logging.MustGetLogger("commento")

func Emit(err error) {
	fmt.Printf("%s", fmt.Errorf("%s", err))
}

func Die(err error) {
	Logger.Panic(fmt.Errorf("%s", err))
}
