package main

import (
	"github.com/op/go-logging"
)

var logger *logging.Logger

func loggerCreate() error {
	format := logging.MustStringFormatter("[%{level}] %{shortfile} %{shortfunc}(): %{message}")
	logging.SetFormatter(format)
	logger = logging.MustGetLogger("commento")

	return nil
}
