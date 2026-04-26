package app

import (
	"autossh/src/i18n"
	"autossh/src/utils"
)

func showVersion() {
	utils.Logln(i18n.T("version_info", Version, Build))
	utils.Logln(i18n.T("author_info"))
}
