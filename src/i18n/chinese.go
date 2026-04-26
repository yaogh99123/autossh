package i18n

var zhCN = map[string]string{
	"version_info": "autossh %s 编译版本 %s。",
	"author_info":  "由 Timo 编写，项目地址：https://github.com/yaogh99123/autossh。",
	"help_usage": `一个ssh远程客户端，可一键登录远程服务器，主要用来弥补Mac/Linux Terminal ssh无法保存密码的不足。
Usage:
  autossh [options] [commands]

Options:
  -c, -config string    指定配置文件(default: ~/.config/autossh/config.yml)。
  -v, -version          显示版本信息。
  -h, -help             显示帮助信息。

Commands:
  cp [-r] source target    复制传输。
  upload/up [-r] [-j 并发数] local server:remote    简化的本地上传功能（支持多线程）。
  down [-r] [-j 并发数] server:remote local    从服务器下载到本地（支持多线程）。
  ${ServerNum}             使用编号登录指定服务器。
  ${ServerAlias}           使用别名登录指定服务器。
  upgrade                  检测并更新到最新版本。
`,
	"menu_add":     "添加",
	"menu_edit":    "编辑",
	"menu_remove":  "删除",
	"menu_exit":    "退出",
	"menu_search":  "搜索",
	"menu_menu":    "菜单",
	"menu_all":     "全部",
	"menu_hide":    "隐藏",
	
	"welcome_autossh":         " 欢迎使用 Auto SSH ",
	"more_servers_hidden":     "\n... (更多服务器已隐藏, 输入 'a' 显示全部)",
	"common_tips_prefix":      "常用提示: ",
	"common_tips_content":     "add.添加, edit.编辑, remove.删除, exit.退出",
	"shortcut_prefix":         "快捷指令: ",
	"shortcut_content":        "[s]搜索, [menu]菜单, [a]显示全部, [h]隐藏多余",
	"select_function":         "请选择功能 [序号, 别名, s]: ",
	
	"config_not_found_init":   "配置文件不存在，正在初始化默认配置...",
	"config_init_fail":        "初始化配置失败",
}
