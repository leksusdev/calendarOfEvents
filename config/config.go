package config

const (
	DataDir         = "data/"
	DataFileName    = DataDir + "calendar.json"
	LogFileName     = DataDir + "app.log"
	LogArchiveName  = DataDir + "console-log.zip"
	ZipLogEntryName = "console.log"

	PromptPrefix         = "> "
	PromptMaxSuggestions = 3

	ListColWidthID     = 37
	ListColWidthTitle  = 51
	ListColWidthDate   = 17
	ListColWidthStatus = 7
	ListTitlePad       = 50

	PrettyJSON = true
)
