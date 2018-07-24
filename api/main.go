package main

func main() {
	exitIfError(createLogger())
	exitIfError(configParse())
	exitIfError(connectDB(5))
	exitIfError(performMigrations())
	exitIfError(smtpConfigure())
	exitIfError(smtpTemplatesLoad())
	exitIfError(oauthConfigure())
	exitIfError(createMarkdownRenderer())
	exitIfError(setupSigintCleanup())

	exitIfError(serveRoutes())
}
