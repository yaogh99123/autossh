package app

import (
	"autossh/src/i18n"
	"autossh/src/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// showSync 处理 rsync 同步需求
func showSync(configFile string) {
	cfg, err := loadConfig(configFile)
	if err != nil {
		utils.Errorln(err)
		return
	}

	// 查找同步源和目的
	args := os.Args
	var syncArgs []string
	var targetServer *Server

	for _, arg := range args {
		// 跳过程序名和子命令名
		if arg == "sync" || strings.HasSuffix(os.Args[0], arg) {
			continue
		}
		// 跳过配置参数
		if strings.HasPrefix(arg, "-c") || strings.HasPrefix(arg, "-config") {
			continue
		}

		if strings.Contains(arg, ":") {
			parts := strings.SplitN(arg, ":", 2)
			aliasOrIndex := parts[0]
			path := parts[1]

			if serverIndex, exists := cfg.serverIndex[aliasOrIndex]; exists {
				targetServer = serverIndex.server
				// 转换为 rsync 识别的格式: user@ip:path
				syncArgs = append(syncArgs, fmt.Sprintf("%s@%s:%s", targetServer.User, targetServer.Ip, path))
				continue
			}
		}
		syncArgs = append(syncArgs, arg)
	}

	if len(syncArgs) < 2 {
		fmt.Println(i18n.T("sync_usage"))
		fmt.Println(i18n.T("sync_example1"))
		fmt.Println(i18n.T("sync_example2"))
		return
	}

	if targetServer == nil {
		utils.Errorln("Error: No remote server alias or index found in arguments.")
		return
	}

	// 在传输前尝试获取源的大小并打印提示
	printSourceSize(targetServer, syncArgs[0])

	runRsync(targetServer, syncArgs)
}

func printSourceSize(srv *Server, src string) {
	var sizeStr string
	// 如果 src 包含 @，说明是远程源 (user@ip:path)
	if strings.Contains(src, "@") && strings.Contains(src, ":") {
		parts := strings.SplitN(src, ":", 2)
		remotePath := parts[1]

		// 构造 SSH 命令获取远程大小
		sshOptions := []string{fmt.Sprintf("-p %d", srv.Port), "-o StrictHostKeyChecking=no", "-o UserKnownHostsFile=/dev/null"}
		if srv.Method == "key" {
			keyPath, _ := utils.ParsePath(string(srv.Key))
			sshOptions = append(sshOptions, fmt.Sprintf("-i %s", keyPath))
		}

		// 执行 du -sh
		cmdStr := fmt.Sprintf("ssh %s %s@%s 'du -sh \"%s\" | cut -f1'", strings.Join(sshOptions, " "), srv.User, srv.Ip, remotePath)
		out, err := exec.Command("sh", "-c", cmdStr).Output()
		if err == nil {
			sizeStr = strings.TrimSpace(string(out))
		}
	} else {
		// 本地源
		absPath, _ := utils.ParsePath(src)
		info, err := os.Stat(absPath)
		if err == nil {
			if info.IsDir() {
				// 简单的本地目录大小统计 (调用系统 du)
				out, _ := exec.Command("du", "-sh", absPath).Output()
				fields := strings.Fields(string(out))
				if len(fields) > 0 {
					sizeStr = fields[0]
				}
			} else {
				sizeStr = utils.SizeFormat(float64(info.Size()))
			}
		}
	}

	if sizeStr != "" {
		fmt.Printf("Detected target (size: %s), using rsync for transfer...\n", sizeStr)
	}
}

func runRsync(srv *Server, args []string) {
	srv.Format()

	// 构建 SSH 参数
	sshOptions := []string{
		fmt.Sprintf("-p %d", srv.Port),
		"-o StrictHostKeyChecking=no",
		"-o UserKnownHostsFile=/dev/null",
	}

	if srv.Method == "key" {
		keyPath, _ := utils.ParsePath(string(srv.Key))
		sshOptions = append(sshOptions, fmt.Sprintf("-i %s", keyPath))
	}

	// 使用 --rsh 代替 -e，这样参数传递更稳定
	rshCmd := fmt.Sprintf("ssh %s", strings.Join(sshOptions, " "))
	rsyncArgs := []string{
		"-az",
		"--partial",
		"--info=progress2",
		"--rsh=" + rshCmd,
	}
	rsyncArgs = append(rsyncArgs, args...)

	// utils.Infoln(fmt.Sprintf("Executing: rsync %s", strings.Join(rsyncArgs, " ")))

	rsyncPath, err := exec.LookPath("rsync")
	if err != nil {
		utils.Errorln("Error: rsync not found in PATH. Please install rsync first.")
		return
	}

	cmd := exec.Command(rsyncPath, rsyncArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		utils.Errorln(fmt.Sprintf("\nRsync failed: %v", err))
	}
}
