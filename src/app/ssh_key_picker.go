package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	fzf "github.com/junegunn/fzf/src"
)

// pickSSHKey 弹出 fzf 选择器选择 SSH 密钥
func pickSSHKey() (string, error) {
	keys, err := listSSHKeys()
	if err != nil {
		return "", err
	}

	if len(keys) == 0 {
		return "", fmt.Errorf("~/.ssh 目录下未找到私钥文件")
	}

	// 准备数据通道
	inputChan := make(chan string)
	go func() {
		for _, k := range keys {
			inputChan <- k
		}
		close(inputChan)
	}()

	// 配置 fzf 参数
	outputChan := make(chan string, 1)
	fzfArgs := []string{
		"--reverse",
		"--height=30%",
		"--prompt=选择 SSH 密钥 (~/.ssh/)> ",
		"--header=模糊搜索私钥, Esc 手动输入",
		"--bind=esc:print(ESC)+abort",
		"--bind=ctrl-c:print(CTRL-C)+abort",
	}

	options, err := fzf.ParseOptions(true, fzfArgs)
	if err != nil {
		return "", err
	}

	options.Input = inputChan
	options.Output = outputChan

	// 运行 fzf
	code, err := fzf.Run(options)
	if code == 130 {
		out := ""
		select {
		case out = <-outputChan:
		case <-time.After(50 * time.Millisecond):
		}
		if out == "ESC" || out == "CTRL-C" {
			return "", nil // 用户主动退出，切回手动输入
		}
	}

	if err != nil {
		return "", err
	}

	// 获取结果
	select {
	case selected := <-outputChan:
		return selected, nil
	default:
		return "", nil
	}
}

// listSSHKeys 列出 ~/.ssh 下可能的私钥文件
func listSSHKeys() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	sshDir := filepath.Join(home, ".ssh")
	files, err := os.ReadDir(sshDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var keys []string
	excludeFiles := map[string]bool{
		"known_hosts":     true,
		"known_hosts.old": true,
		"config":          true,
		"authorized_keys": true,
		"environment":     true,
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()

		// 过滤规则
		if excludeFiles[name] {
			continue
		}
		if strings.HasSuffix(name, ".pub") {
			continue
		}
		if strings.HasPrefix(name, "id_") || strings.HasSuffix(name, ".pem") || strings.HasSuffix(name, ".key") {
			keys = append(keys, filepath.Join(sshDir, name))
		} else {
			// 对于不符合常见模式但可能是私钥的文件，也可以考虑加入，或者保持严谨
			// 这里我们保持包含常见模式
		}
	}

	return keys, nil
}
