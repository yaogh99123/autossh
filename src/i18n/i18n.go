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
	if strings.Contains(langEnv, "zh_cn") || strings.Contains(langEnv, "zh-cn") ||
		strings.Contains(langEnv, "zh_tw") || strings.Contains(langEnv, "zh-tw") ||
		strings.Contains(langEnv, "zh_hk") {
		currentLang = LangZH
	} else {
		currentLang = LangEN
	}
}

// SetLanguage
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
		dict = dicts[LangEN]
	}

	format, ok := dict[key]
	if !ok {
		format = key
	}

	if len(args) > 0 {
		return fmt.Sprintf(format, args...)
	}
	return format
}
