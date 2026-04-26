package i18n

import (
	"fmt"
	"os"
	"strings"
)

type Language string

const (
	LangEN Language = "en"
	LangZH Language = "zh"
)

var currentLang = LangEN



// dicts 聚合所有语言包
var dicts = map[Language]map[string]string{
	LangEN: enUS,
	LangZH: zhCN,
}

func init() {
	detectLanguage()
}

// detectLanguage 通过环境变量自动检测系统语言
func detectLanguage() {
	langEnv := os.Getenv("LANG")
	if langEnv == "" {
		langEnv = os.Getenv("LC_ALL")
	}

	langEnv = strings.ToLower(langEnv)
	// 判断是否包含中文特征
	if strings.Contains(langEnv, "zh_cn") || strings.Contains(langEnv, "zh-cn") ||
		strings.Contains(langEnv, "zh_tw") || strings.Contains(langEnv, "zh-tw") ||
		strings.Contains(langEnv, "zh_hk") {
		currentLang = LangZH
	} else {
		currentLang = LangEN // 默认 Fallback 到英文
	}
}

// SetLanguage 允许通过配置文件或命令行强制覆盖系统语言
func SetLanguage(lang string) {
	lang = strings.ToLower(lang)
	if strings.HasPrefix(lang, "zh") {
		currentLang = LangZH
	} else if strings.HasPrefix(lang, "en") {
		currentLang = LangEN
	}
}

// GetCurrentLang 获取当前使用的语言
func GetCurrentLang() Language {
	return currentLang
}

// T 翻译核心函数，支持占位符参数（类似 fmt.Sprintf）
func T(key string, args ...interface{}) string {
	dict, ok := dicts[currentLang]
	if !ok {
		// 防御性编程，如果当前语言字典不存在，回退到英文
		dict = dicts[LangEN]
	}

	format, ok := dict[key]
	if !ok {
		// 如果翻译键不存在，为了防止完全无输出，直接返回 key 本身作为 Fallback
		format = key
	}

	// 如果传入了参数，则执行格式化插值
	if len(args) > 0 {
		return fmt.Sprintf(format, args...)
	}
	return format
}
