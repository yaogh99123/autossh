package app

import (
	"autossh/src/i18n"
	"autossh/src/utils"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
	utils.Blueln("Calculating target size... ")

	var sizeStr string
	// 如果 src 包含 @，说明是远程源 (user@ip:path)
	if strings.Contains(src, "@") && strings.Contains(src, ":") {
		parts := strings.SplitN(src, ":", 2)
		remotePath := parts[1]

		sshArgs := []string{
			fmt.Sprintf("-p%d", srv.Port),
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			"-o", "ConnectTimeout=10",
		}
		if srv.Method == "key" {
			keyPath, _ := utils.ParsePath(string(srv.Key))
			sshArgs = append(sshArgs, "-i", keyPath)
		}
		sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", srv.User, srv.Ip))

		// 优化远程探测：增加 2>/dev/null 并放宽超时到 60 秒
		remoteCmd := fmt.Sprintf("[ -f \"%s\" ] && stat -c%%s \"%s\" || du -sh \"%s\" 2>/dev/null | cut -f1", remotePath, remotePath, remotePath)
		sshArgs = append(sshArgs, remoteCmd)

		// 使用带 context 的命令，设置总超时 60 秒
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		out, err := exec.CommandContext(ctx, "ssh", sshArgs...).Output()
		if err == nil {
			raw := strings.TrimSpace(string(out))
			// 如果返回的是纯数字（stat 的结果），转换一下格式
			if val, err := strconv.ParseFloat(raw, 64); err == nil {
				sizeStr = utils.SizeFormat(val)
			} else {
				sizeStr = raw
			}
		}
	} else {
		// 本地源处理
		absPath, _ := utils.ParsePath(src)
		info, err := os.Stat(absPath)
		if err == nil {
			if info.IsDir() {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				out, _ := exec.CommandContext(ctx, "du", "-sh", absPath).Output()
				fields := strings.Fields(string(out))
				if len(fields) > 0 {
					sizeStr = fields[0]
				}
			} else {
				sizeStr = utils.SizeFormat(float64(info.Size()))
			}
		}
	}

	// 清除 "Calculating..." 行并打印结果
	fmt.Print("\r\033[K") // 清除当前行
	if sizeStr != "" {
		utils.Blueln(fmt.Sprintf("Detected target (size: %s), using rsync for transfer...", sizeStr))
	} else {
		utils.Blueln("Size calculation It's too big, starting transfer directly...")
	}
}

func runRsync(srv *Server, args []string) {
	srv.Format()

	// 构建 SSH 参数
	sshOptions := []string{
		fmt.Sprintf("-p%d", srv.Port),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
	}

	if srv.Method == "key" {
		keyPath, _ := utils.ParsePath(string(srv.Key))
		sshOptions = append(sshOptions, "-i", keyPath)
	}

	sshCmd := "ssh " + strings.Join(sshOptions, " ")
	rsyncArgs := []string{
		"-az",
		"--partial",
		"--info=progress2",
		"-e", sshCmd,
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
