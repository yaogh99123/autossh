package app

import (
	"autossh/src/i18n"
	"autossh/src/utils"
	"io"
)

func handleEdit(cfg *Config, args []string) error {
	id, err := utils.ReadLine(utils.Colored(i18n.T("edit_enter_index"), utils.ColorGreen))
	if err == io.EOF || (err != nil && err.Error() == "Interrupt") {
		return nil
	}

	serverIndex, ok := cfg.serverIndex[id]
	if !ok {
		utils.Errorln(i18n.T("edit_index_not_found"))
		return handleEdit(cfg, args)
	}

	if err := serverIndex.server.Edit(); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return cfg.saveConfig(true)
}
