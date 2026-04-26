package app

import (
	"autossh/src/i18n"
	"autossh/src/utils"
	"fmt"
	"strings"
	"time"

	fzf "github.com/junegunn/fzf/src"
)

// handleSearch 处理搜索并连接
func handleSearch(cfg *Config, args []string) error {
	// 1. 收集所有服务器及其对应的显示序号
	type serverInfo struct {
		server *Server
		index  string
	}
	var allInfos []serverInfo

	// 顶级服务器序号: 1, 2, 3...
	for i, s := range cfg.Servers {
		allInfos = append(allInfos, serverInfo{server: s, index: fmt.Sprintf("%d", i+1)})
	}

	// 分组服务器序号: a1, a2, b1, g1...
	for i := range cfg.Groups {
		g := cfg.Groups[i]
		for j := range g.Servers {
			s := &g.Servers[j]
			allInfos = append(allInfos, serverInfo{server: s, index: fmt.Sprintf("%s%d", g.Prefix, j+1)})
		}
	}

	if len(allInfos) == 0 {
		return fmt.Errorf(i18n.T("search_no_servers"))
	}

	// 2. 准备数据通道
	inputChan := make(chan string)
	go func() {
		for _, info := range allInfos {
			s := info.server
			groupInfo := ""
			if s.groupName != "" {
				groupInfo = "[" + s.groupName + "] "
			}
			// 格式: 序号 \t 名称 \t 分组信息 (用户@IP)
			// 注意: 这样序号也会被 fzf 检索到
			inputChan <- fmt.Sprintf("%s\t%s\t%s(%s@%s)", info.index, s.Name, groupInfo, s.User, s.Ip)
		}
		close(inputChan)
	}()

	// 3. 配置 fzf 参数 (不清屏模式)
	outputChan := make(chan string, 1)
	fzfArgs := []string{
		"--reverse",
		"--height=40%",
		"--prompt=" + i18n.T("search_prompt"),
		"--header=" + i18n.T("search_header"),
		"--bind=esc:print(ESC)+abort",
		"--bind=ctrl-c:print(CTRL-C)+abort",
		"--delimiter=\t",
		"--with-nth=1,2,3", // 显示序号、名称和详细信息
	}

	options, err := fzf.ParseOptions(true, fzfArgs)
	if err != nil {
		return fmt.Errorf(i18n.T("search_err_init", err))
	}

	options.Input = inputChan
	options.Output = outputChan

	// 4. 运行 fzf
	code, err := fzf.Run(options)
	if code == 130 {
		// 处理退出
		out := ""
		select {
		case out = <-outputChan:
		case <-time.After(50 * time.Millisecond):
		}
		if out == "CTRL-C" {
			return nil
		}
		return nil // ESC 退出
	}

	if err != nil {
		return fmt.Errorf(i18n.T("search_err_run", code, err))
	}

	// 5. 处理结果
	select {
	case selected := <-outputChan:
		if selected != "" {
			parts := strings.Split(selected, "\t")
			if len(parts) >= 2 {
				name := strings.TrimSpace(parts[1]) // 第 2 列是名称
				index := strings.TrimSpace(parts[0]) // 第 1 列是序号
				
				// 在 allInfos 中找回对象
				var target *Server
				for _, info := range allInfos {
					if info.server.Name == name && info.index == index {
						target = info.server
						break
					}
				}
				if target != nil {
					utils.Infoln("\n"+i18n.T("search_selected", target.Name))
					return target.Connect()
				}
			}
		}
	default:
	}

	return nil
}

// handleMenu 处理菜单显示
func handleMenu(cfg *Config, args []string) error {
	return nil
}
