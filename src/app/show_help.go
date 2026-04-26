package app

import (
	"autossh/src/i18n"
	"autossh/src/utils"
	"flag"
)

func showHelp() {
	flag.Usage()
}

func usage() {
	utils.Logln(i18n.T("help_usage"))
}
