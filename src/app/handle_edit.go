package app

import (
	"autossh/src/i18n"
	"autossh/src/utils"
	"fmt"
	"io"
)

func handleEdit(cfg *Config, args []string) error {
	utils.Info(i18n.T("edit_enter_index"))
	id := ""
	if _, err := fmt.Scanln(&id); err == io.EOF {
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
