package app

import (
	"autossh/src/i18n"
	"autossh/src/utils"
	"strconv"
)

func showServers(configFile string) {
	cfg, err := loadConfig(configFile)
	if err != nil {
		utils.Errorln(err)
		return
	}

	// 清屏
	_ = utils.Clear()

	show(cfg)

	for {
		loop, clear, reload := scanInput(cfg)
		if !loop {
			break
		}

		if reload {
			cfg, err = loadConfig(configFile)
		}

		if clear {
			_ = utils.Clear()
		}

		show(cfg)
	}
}

// 显示服务
func show(cfg *Config) {
	maxlen := separatorLength(*cfg)
	flagWidth, nameWidth := calculateWidths(cfg)
	utils.Blueln(utils.FormatSeparator(i18n.T("welcome_autossh"), "=", maxlen))

	count := 0
	limit := 8
	hasMore := false

	for i, server := range cfg.Servers {
		if !cfg.ShowAll && count >= limit {
			hasMore = true
			break
		}
		utils.Logln(server.FormatPrint(strconv.Itoa(i+1), cfg.ShowDetail, flagWidth, nameWidth))
		count++
	}

	if !(!cfg.ShowAll && count >= limit) {
		for _, group := range cfg.Groups {
			if len(group.Servers) == 0 {
				continue
			}

			var collapseNotice = ""
			if group.Collapse {
				collapseNotice = "[" + group.Prefix + " ↓]"
			} else {
				collapseNotice = "[" + group.Prefix + " ↑]"
			}

			utils.Logln()
			utils.Yellowln(utils.FormatSeparator(" "+group.GroupName+" "+collapseNotice+" ", "_", maxlen))
			if !group.Collapse {
				for i, server := range group.Servers {
					if !cfg.ShowAll && count >= limit {
						hasMore = true
						break
					}
					utils.Logln(server.FormatPrint(group.Prefix+strconv.Itoa(i+1), cfg.ShowDetail, flagWidth, nameWidth))
					count++
				}
			}
			if !cfg.ShowAll && count >= limit {
				break
			}
		}
	}

	if hasMore {
		utils.Yellowln(i18n.T("more_servers_hidden"))
	}

	utils.Logln()
	utils.Blueln(utils.FormatSeparator("", "=", maxlen))

	// 常用提示 (对齐 dcli 风格)
	utils.Log(utils.Colored(i18n.T("common_tips_prefix"), utils.ColorYellow))
	utils.Logln(i18n.T("common_tips_content"))

	// 快捷指令
	utils.Log(utils.Colored(i18n.T("shortcut_prefix"), utils.ColorYellow))
	utils.Logln(i18n.T("shortcut_content"))

	// 传输指令 (支持多线程)
	// utils.Log(utils.Colored("传输指令: ", utils.ColorYellow))
	// utils.Logln("autossh up/down [-r] [-j 并发数] 源 目标")
	// utils.Log(utils.Colored("      示例: ", utils.ColorBlue))
	// utils.Logln("autossh up -r -j 10 ./dist/ server:/var/www/html")
	// utils.Log(utils.Colored("            ", utils.ColorBlue))
	// utils.Logln("autossh down server:/logs/app.log ./local_logs/")

	utils.Blueln(utils.FormatSeparator("", "=", maxlen))
	utils.Cyanln(i18n.T("select_function"))
}

// 计算分隔符长度
func separatorLength(cfg Config) int {
	maxlength := 60
	for _, group := range cfg.Groups {
		length := utils.ZhLen(group.GroupName)
		if length > maxlength {
			maxlength = length + 10
		}
	}

	return maxlength
}

// 预计算所有列的最大宽度
func calculateWidths(cfg *Config) (int, int) {
	flagWidth := 0
	nameWidth := 0

	// 辅助计算函数
	update := func(s *Server, flag string) {
		alias := ""
		if s.Alias != "" {
			alias = "|" + s.Alias
		}
		fullFlag := "[" + flag + alias + "]"
		fw := utils.ZhLen(fullFlag)
		if fw > flagWidth {
			flagWidth = fw
		}
		nw := utils.ZhLen(s.Name)
		if nw > nameWidth {
			nameWidth = nw
		}
	}

	for i, server := range cfg.Servers {
		update(server, strconv.Itoa(i+1))
	}

	for _, group := range cfg.Groups {
		for i := range group.Servers {
			update(&group.Servers[i], group.Prefix+strconv.Itoa(i+1))
		}
	}

	return flagWidth, nameWidth
}
