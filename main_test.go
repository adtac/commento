package main

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testsSetup()
	exitStatus := m.Run()
	testsTeardown()
	os.Exit(exitStatus)
}

const testdb = "sqlite3_test.db"

func testsSetup() {
	if err := LoadDatabase(testdb); err != nil {
		log.Fatalf("unable to create test db: %v", err)
	}
}
func testsTeardown() {
	if err := os.Remove(testdb); err != nil {
		log.Fatalf("unable to remove test db: %v", err)
	}
}
