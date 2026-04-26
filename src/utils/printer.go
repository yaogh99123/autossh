package utils

import "fmt"

// 颜色定义 (对齐 dcli 风格)
const (
	ColorRed    = "\033[0;31m"
	ColorGreen  = "\033[0;32m"
	ColorYellow = "\033[1;33m"
	ColorBlue   = "\033[0;34m"
	ColorCyan   = "\033[0;36m"
	ColorBold   = "\033[1m"
	ColorNC     = "\033[0m" // No Color
)

// 打印一行信息
func Logln(a ...interface{}) {
	fmt.Println(a...)
}

// 打印（不换行）
func Log(a ...interface{}) {
	fmt.Print(a...)
}

// 打印一行信息 (对齐 dcli Green 风格)
func Infoln(a ...interface{}) {
	fmt.Print(ColorGreen)
	Logln(a...)
	fmt.Print(ColorNC)
}

// 打印信息（不换行）
func Info(a ...interface{}) {
	fmt.Print(ColorGreen)
	fmt.Print(a...)
	fmt.Print(ColorNC)
}

// 打印一行错误 (对齐 dcli Red 风格)
func Errorln(a ...interface{}) {
	fmt.Print(ColorRed)
	Logln(a...)
	fmt.Print(ColorNC)
}

// 打印黄色信息
func Yellowln(a ...interface{}) {
	fmt.Print(ColorYellow)
	Logln(a...)
	fmt.Print(ColorNC)
}

// 打印青色信息
func Cyanln(a ...interface{}) {
	fmt.Print(ColorCyan)
	Logln(a...)
	fmt.Print(ColorNC)
}

// 打印蓝色信息
func Blueln(a ...interface{}) {
	fmt.Print(ColorBlue)
	Logln(a...)
	fmt.Print(ColorNC)
}

// 格式化输出颜色字符串
func Colored(msg string, colorCode string) string {
	return colorCode + msg + ColorNC
}

// 二维数组对齐
//func Align(arr [][]string) [][]string {
//	for column := 0; column < 2; column++ {
//		columnWidth := getColumnWidth(arr, column)
//
//		for index := range arr {
//			arr[index][column] = AppendRight(arr[index][column], " ", columnWidth)
//		}
//	}
//
//	return arr
//}
//
//func getColumnWidth(arr [][]string, column int) int {
//	maxWidth := 0
//	for _, row := range arr {
//		width := int(ZhLen(row[column]))
//		if maxWidth < width {
//			maxWidth = width
//		}
//	}
//
//	return maxWidth
//}
