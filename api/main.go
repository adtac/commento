package main

func main() {
	exitIfError(loggerCreate())
	exitIfError(versionPrint())
	exitIfError(configParse())
	exitIfError(dbConnect(5))
	exitIfError(migrate())
	exitIfError(smtpConfigure())
	exitIfError(smtpTemplatesLoad())
	exitIfError(emailProviderConfigure())
	exitIfError(oauthConfigure())
	exitIfError(markdownRendererCreate())
	exitIfError(superAdminSetup())
	rateLimiterInit()
	exitIfError(sessionCleanupBegin())
	exitIfError(sigintCleanupSetup())
	exitIfError(versionCheckStart())
	exitIfError(domainExportCleanupBegin())
	exitIfError(viewsCleanupBegin())
	exitIfError(ssoTokenCleanupBegin())

	exitIfError(routesServe())
}
