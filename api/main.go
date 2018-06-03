package main

func main() {
	exitIfError(createLogger())
	exitIfError(parseConfig())
	exitIfError(connectDB())
	exitIfError(performMigrations())
	exitIfError(smtpConfigure())
	exitIfError(oauthConfigure())
	exitIfError(createMarkdownRenderer())
	exitIfError(setupSigintCleanup())

	exitIfError(serveRoutes())
}
