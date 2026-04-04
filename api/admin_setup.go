package main

import (
	"os"
)

func superAdminSetup() error {
	email := os.Getenv("SUPER_ADMIN_EMAIL")
	if email == "" {
		return nil
	}

	statement := `
		UPDATE owners
		SET isSuperAdmin = true
		WHERE email = $1;
	`
	_, err := db.Exec(statement, email)
	if err != nil {
		logger.Errorf("cannot set super admin: %v", err)
		return err
	}

	logger.Infof("super admin configured for %s", email)
	return nil
}
