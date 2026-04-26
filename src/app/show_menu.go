package app

import (
	"autossh/src/i18n"
	"autossh/src/utils"
	"strings"
)

type Operation struct {
	Key     string
	Label   string
	End     bool
	Process func(cfg *Config, args []string) error
}

var menuMap [][]Operation

var operations = make(map[string]Operation)

func init() {
	menuMap = [][]Operation{
		{
			{Key: "add", Label: "menu_add", Process: handleAdd},
			{Key: "edit", Label: "menu_edit", Process: handleEdit},
			{Key: "remove", Label: "menu_remove", Process: handleRemove},
		},
		{
			{Key: "exit", Label: "menu_exit", End: true},
		},
	}

	// 初始化全局 operations 字典
	for _, row := range menuMap {
		for _, op := range row {
			operations[op.Key] = op
		}
	}

	// 注册隐藏的/快捷指令操作
	operations["s"] = Operation{Key: "s", Label: "menu_search", Process: handleSearch}
	operations["menu"] = Operation{Key: "menu", Label: "menu_menu", Process: handleMenu}
	operations["a"] = Operation{Key: "a", Label: "menu_all", Process: func(cfg *Config, args []string) error {
		cfg.ShowAll = true
		return cfg.saveConfig(false)
	}}
	operations["h"] = Operation{Key: "h", Label: "menu_hide", Process: func(cfg *Config, args []string) error {
		cfg.ShowAll = false
		return cfg.saveConfig(false)
	}}
}

func showMenu() {
	var columnsMaxWidths = make(map[int]int)

	for i := 0; i < len(menuMap); i++ {
		for j := 0; j < len(menuMap[i]); j++ {
			operation := menuMap[i][j]

			// 计算每列最大长度
			maxLen := int(utils.ZhLen(operationFormat(operation)))
			if _, exists := columnsMaxWidths[j]; !exists {
				columnsMaxWidths[j] = maxLen
			}
			if columnsMaxWidths[j] < maxLen {
				columnsMaxWidths[j] = maxLen
			}

			operations[operation.Key] = operation
		}
	}

	for i := 0; i < len(menuMap); i++ {
		var output = ""
		for j := 0; j < len(menuMap[i]); j++ {
			operation := menuMap[i][j]
			output += stringPadding(operationFormat(operation), columnsMaxWidths[j]) + "\t"
		}

		utils.Logln(strings.TrimSpace(output))
		output = ""
	}
}

func operationFormat(operation Operation) string {
	return utils.Colored("["+operation.Key+"]", utils.ColorGreen) + " " + i18n.T(operation.Label)
}

func stringPadding(str string, paddingLen int) string {
	if len(str) < paddingLen {
		return stringPadding(str+" ", paddingLen)
	} else {
		return str
	}
}
