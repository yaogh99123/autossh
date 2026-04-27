package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

// ReadLine 统一提供带有退格、历史记录支持的终端读取，避免终端乱码
func ReadLine(prompt string) (string, error) {
	rl, err := readline.New(prompt)
	if err != nil {
		return "", err
	}
	defer rl.Close()

	line, err := rl.Readline()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func ReadPassword(prompt string) (string, error) {
	fmt.Println(prompt)
	rl, err := readline.New("")
	if err != nil {
		return "", err
	}
	defer rl.Close()

	cfg := rl.GenPasswordConfig()
	cfg.MaskRune = '*'
	b, err := rl.ReadPasswordWithConfig(cfg)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

// Scanln 保持对老代码的兼容，但底层改用原生的 readline
func Scanln(a *string) {
	line, err := ReadLine("")
	if err != nil {
		if err.Error() == "Interrupt" || err.Error() == "EOF" {
			fmt.Println()
			os.Exit(0)
		}
		fmt.Println("读取输入失败:", err)
		return
	}
	*a = line
}
