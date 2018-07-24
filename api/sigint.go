package main

import (
	"os"
	"os/signal"
	"syscall"
)

func sigintCleanup() int {
	// TODO: close the database connection and do other cleanup jobs
	return 0
}

func sigintCleanupSetup() error {
	logger.Infof("setting up SIGINT cleanup")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		<-c
		os.Exit(sigintCleanup())
	}()

	return nil
}
