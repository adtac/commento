package main

import (
	"github.com/op/go-logging"
)

var logger *logging.Logger

func createLogger() error {
	format := logging.MustStringFormatter("[%{level}] %{shortfile} %{shortfunc}(): %{message}")
	logging.SetFormatter(format)
	logger = logging.MustGetLogger("commento")

	return nil
}
