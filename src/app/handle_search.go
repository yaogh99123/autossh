package app

import (
	"autossh/src/utils"
	"fmt"
	"strings"
	"time"

	fzf "github.com/junegunn/fzf/src"
)

// handleSearch 处理搜索并连接
func handleSearch(cfg *Config, args []string) error {
	var allServers []*Server

	// 1. 收集所有服务器
	for _, s := range cfg.Servers {
		allServers = append(allServers, s)
	}

	for i := range cfg.Groups {
		g := cfg.Groups[i]
		for j := range g.Servers {
			s := &g.Servers[j]
			allServers = append(allServers, s)
		}
	}

	if len(allServers) == 0 {
		return fmt.Errorf("没有可用的服务器")
	}

	// 2. 准备数据通道
	inputChan := make(chan string)
	go func() {
		for _, s := range allServers {
			groupInfo := ""
			if s.groupName != "" {
				groupInfo = "[" + s.groupName + "] "
			}
			// 格式: 名称 \t 分组信息 (用户@IP)
			inputChan <- fmt.Sprintf("%s\t%s(%s@%s)", s.Name, groupInfo, s.User, s.Ip)
		}
		close(inputChan)
	}()

	// 3. 配置 fzf 参数 (不清屏模式)
	outputChan := make(chan string, 1)
	fzfArgs := []string{
		"--reverse",
		"--height=40%",
		"--prompt=搜索服务器> ",
		"--header=快捷菜单搜索 (fzf 模式, Esc 退出)",
		"--bind=esc:print(ESC)+abort",
		"--bind=ctrl-c:print(CTRL-C)+abort",
		"--delimiter=\t",
	}

	options, err := fzf.ParseOptions(true, fzfArgs)
	if err != nil {
		return fmt.Errorf("fzf 初始化失败: %v", err)
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
		return fmt.Errorf("fzf 运行失败 (code %d): %v", code, err)
	}

	// 5. 处理结果
	select {
	case selected := <-outputChan:
		if selected != "" {
			parts := strings.Split(selected, "\t")
			if len(parts) > 0 {
				name := strings.TrimSpace(parts[0])
				// 在 allServers 中找回对象
				var target *Server
				for _, s := range allServers {
					if s.Name == name {
						target = s
						break
					}
				}
				if target != nil {
					utils.Infoln("\n你选择了", target.Name)
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
