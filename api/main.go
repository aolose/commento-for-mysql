package main

func main() {
	exitIfError(loggerCreate())
	exitIfError(configParse())
	exitIfError(dbConnect(5))
	exitIfError(smtpConfigure())
	exitIfError(smtpTemplatesLoad())
	exitIfError(oauthConfigure())
	exitIfError(markdownRendererCreate())
	exitIfError(emailNotificationPendingResetAll())
	exitIfError(emailNotificationBegin())
	exitIfError(sigintCleanupSetup())
	exitIfError(versionCheckStart())
	exitIfError(domainExportCleanupBegin())
	exitIfError(viewsCleanupBegin())
	exitIfError(routesServe())
}
