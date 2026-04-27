package app

import (
	"autossh/src/i18n"
	"autossh/src/utils"
	"io"
)

func handleRemove(cfg *Config, args []string) error {
	id, err := utils.ReadLine(utils.Colored(i18n.T("remove_enter_index"), utils.ColorGreen))
	if err == io.EOF || (err != nil && err.Error() == "Interrupt") {
		return nil
	}

	serverIndex, ok := cfg.serverIndex[id]
	if !ok {
		utils.Errorln(i18n.T("remove_index_not_found"))
		return handleRemove(cfg, args)
	}

	if serverIndex.indexType == IndexTypeServer {
		servers := cfg.Servers
		cfg.Servers = append(servers[:serverIndex.serverIndex], servers[serverIndex.serverIndex+1:]...)
	} else {
		servers := cfg.Groups[serverIndex.groupIndex].Servers
		servers = append(servers[:serverIndex.serverIndex], servers[serverIndex.serverIndex+1:]...)
		cfg.Groups[serverIndex.groupIndex].Servers = servers
	}

	return cfg.saveConfig(true)
}
