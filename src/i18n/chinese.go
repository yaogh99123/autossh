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

	"cp_err_missing_args":       "请输入完整参数",
	"cp_err_both_local":         "源和目标不能同时为本地地址",
	"cp_err_is_dir":             "是一个目录",
	"cp_err_format":             "%s 格式错误",
	"err_server_not_found":      "服务器 %s 不存在",
	"err_conn_fail":             "连接服务器失败: %v",

	"down_usage":                "用法: autossh down [-r] [-j 并发数] <服务器名/序号:远程路径> <本地路径>",
	"down_example1":             "示例: autossh down server1:/home/user/file.txt ./file.txt",
	"down_example2":             "示例: autossh down 01:/home/user/file.txt ~/Downloads/file.txt  (使用序号)",
	"down_example3":             "示例: autossh down -r -j 10 01:/home/user/dir/ ~/Downloads/  (10个并发)",
	"down_err_format":           "远程源格式错误，应为: 服务器名/序号:路径",
	"down_err_remote_not_found": "远程文件不存在: %s",
	"down_err_is_dir":           "远程路径是目录，请使用 -r 参数",
	"down_err_download":         "下载失败: %v",
	"down_success":              "下载完成!",
	"down_err_dir_fail":         "下载目录 %s 失败: %v\n",
	"down_err_file_fail":        "下载文件 %s 失败: %v\n",

	"up_usage":                  "用法: autossh upload/up [-r] [-j 并发数] <本地文件/目录> <服务器名/序号:远程路径>",
	"up_example1":               "示例: autossh upload file.txt server1:/home/user/",
	"up_example2":               "示例: autossh up file.txt 01:/home/user/  (使用序号)",
	"up_example3":               "示例: autossh upload -r -j 10 ./localdir 01:/home/user/  (10个并发)",
	"up_err_local_not_found":    "本地文件不存在: %s",
	"up_err_format":             "远程目标格式错误，应为: 服务器名/序号:路径",
	"up_err_local_stat":         "获取本地文件信息失败: %v",
	"up_err_is_dir":             "本地路径是目录，请使用 -r 参数",
	"up_err_upload":             "上传失败: %v",
	"up_success":                "上传完成!",
	"up_err_dir_fail":           "上传目录 %s 失败: %v\n",
	"up_err_file_fail":          "上传文件 %s 失败: %v\n",

	"search_no_servers":         "没有可用的服务器",
	"search_prompt":             "搜索服务器> ",
	"search_header":             "快捷菜单搜索 (fzf 模式, Esc 退出)",
	"search_err_init":           "fzf 初始化失败: %v",
	"search_selected":           "你选择了 %s",
	"search_err_run":            "fzf 运行失败 (code %d): %v",

	"upgrade_checking":          "正在检测最新版本",
	"upgrade_current_ver":       "当前版本：%s",
	"upgrade_up_to_date":        "感谢您的支持，当前已是最新版本。",
	"upgrade_new_ver":           "检测到新版本：%s",
	"upgrade_unsupported_os":    "暂不支持 %s 系统自动更新，请下载源码包手动编译。",
	"upgrade_err_download":      "下载失败：%v",
	"upgrade_err_unzip":         "解压缩失败：%v",
	"upgrade_err_install":       "安装失败",

	"scan_backup_fail":          "备份失败",
	"scan_invalid_input":        "输入有误，请重新输入",

	"remove_enter_index":        "请输入相应序号：",
	"remove_index_not_found":    "序号不存在",

	"edit_enter_index":          "请输入相应序号：",
	"edit_index_not_found":      "序号不存在",
	"add_other_group":           "[其他值]默认组",
	"add_enter_group":           "请输入要插入的组：",

	"app_flag_c":                "指定配置文件路径",
	"app_flag_v":                "版本信息",
	"app_flag_h":                "帮助信息",
}
