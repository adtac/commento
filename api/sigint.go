package main

import (
	"os"
	"os/signal"
	"syscall"
)

func sigintCleanup() int {
	if db != nil {
		err := db.Close()
		if err == nil {
			logger.Errorf("cannot close database connection: %v", err)
			return 1
		}
	}

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
