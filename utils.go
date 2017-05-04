package commento

import "log"

// Emit outputs an error message to the logger in use
func Emit(err error) {
	log.Print(err)
}

// Die logs a fatal error and exits the application
func Die(err error) {
	log.Fatal(err)
}
