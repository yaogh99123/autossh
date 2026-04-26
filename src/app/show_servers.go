package app

import (
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
	utils.Blueln(utils.FormatSeparator(" 欢迎使用 Auto SSH ", "=", maxlen))
	for i, server := range cfg.Servers {
		utils.Logln(server.FormatPrint(strconv.Itoa(i+1), cfg.ShowDetail))
	}

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
				utils.Logln(server.FormatPrint(group.Prefix+strconv.Itoa(i+1), cfg.ShowDetail))
			}
		}
	}

	utils.Logln()
	utils.Blueln(utils.FormatSeparator("", "=", maxlen))

	// 常用提示 (对齐 dcli 风格)
	utils.Log(utils.Colored("常用提示: ", utils.ColorYellow))
	utils.Logln("add.添加, edit.编辑, remove.删除, exit.退出")

	// 快捷指令
	utils.Log(utils.Colored("快捷指令: ", utils.ColorYellow))
	utils.Logln("[s]搜索, [menu]菜单")

	utils.Blueln(utils.FormatSeparator("", "=", maxlen))
	utils.Cyanln("请选择功能 [序号, 别名, s]: ")
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
