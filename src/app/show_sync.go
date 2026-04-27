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

	runRsync(targetServer, syncArgs)
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
		"--info=progress2",
		"--rsh=" + rshCmd,
	}
	rsyncArgs = append(rsyncArgs, args...)

	utils.Infoln(fmt.Sprintf("Executing: rsync %s", strings.Join(rsyncArgs, " ")))

	cmd := exec.Command("rsync", rsyncArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		utils.Errorln(fmt.Sprintf("\nRsync failed: %v", err))
	}
}
