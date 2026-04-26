package i18n

var enUS = map[string]string{
	"version_info": "autossh %s Build %s.",
	"author_info":  "Written by Timo, Project: https://github.com/yaogh99123/autossh.",
	"help_usage": `An SSH remote client with one-click login, mainly used to compensate for the inability of Mac/Linux Terminal to save SSH passwords.
Usage:
  autossh [options] [commands]

Options:
  -c, -config string    Specify configuration file (default: ~/.config/autossh/config.yml).
  -v, -version          Show version information.
  -h, -help             Show help information.

Commands:
  cp [-r] source target    Copy and transfer.
  upload/up [-r] [-j concurrency] local server:remote    Simplified local upload (multi-thread supported).
  down [-r] [-j concurrency] server:remote local    Download from server to local (multi-thread supported).
  ${ServerNum}             Login to specific server using number.
  ${ServerAlias}           Login to specific server using alias.
  upgrade                  Check and upgrade to the latest version.
`,
	"menu_add":     "Add",
	"menu_edit":    "Edit",
	"menu_remove":  "Remove",
	"menu_exit":    "Exit",
	"menu_search":  "Search",
	"menu_menu":    "Menu",
	"menu_all":     "All",
	"menu_hide":    "Hide",
	
	"welcome_autossh":         " Welcome to Auto SSH ",
	"more_servers_hidden":     "\n... (More servers hidden, enter 'a' to show all)",
	"common_tips_prefix":      "Common Tips: ",
	"common_tips_content":     "add.Add, edit.Edit, remove.Remove, exit.Exit",
	"shortcut_prefix":         "Shortcuts: ",
	"shortcut_content":        "[s]Search, [menu]Menu, [a]Show All, [h]Hide Extra",
	"select_function":         "Select function [Index, Alias, s]: ",
	
	"config_not_found_init":   "Config file not found, initializing default config...",
	"config_init_fail":        "Failed to initialize config",

	"cp_err_missing_args":       "Please provide complete parameters",
	"cp_err_both_local":         "Source and target cannot both be local addresses",
	"cp_err_is_dir":             "is a directory",
	"cp_err_format":             "%s format error",
	"err_server_not_found":      "Server %s not found",
	"err_conn_fail":             "Failed to connect to server: %v",

	"down_usage":                "Usage: autossh down [-r] [-j concurrency] <server_name/index:remote_path> <local_path>",
	"down_example1":             "Example: autossh down server1:/home/user/file.txt ./file.txt",
	"down_example2":             "Example: autossh down 01:/home/user/file.txt ~/Downloads/file.txt  (using index)",
	"down_example3":             "Example: autossh down -r -j 10 01:/home/user/dir/ ~/Downloads/  (10 concurrency)",
	"down_err_format":           "Remote source format error, should be: server_name/index:path",
	"down_err_remote_not_found": "Remote file not found: %s",
	"down_err_is_dir":           "Remote path is a directory, please use -r parameter",
	"down_err_download":         "Download failed: %v",
	"down_success":              "Download complete!",
	"down_err_dir_fail":         "Failed to download directory %s: %v\n",
	"down_err_file_fail":        "Failed to download file %s: %v\n",

	"up_usage":                  "Usage: autossh upload/up [-r] [-j concurrency] <local_file/directory> <server_name/index:remote_path>",
	"up_example1":               "Example: autossh upload file.txt server1:/home/user/",
	"up_example2":               "Example: autossh up file.txt 01:/home/user/  (using index)",
	"up_example3":               "Example: autossh upload -r -j 10 ./localdir 01:/home/user/  (10 concurrency)",
	"up_err_local_not_found":    "Local file not found: %s",
	"up_err_format":             "Remote target format error, should be: server_name/index:path",
	"up_err_local_stat":         "Failed to get local file info: %v",
	"up_err_is_dir":             "Local path is a directory, please use -r parameter",
	"up_err_upload":             "Upload failed: %v",
	"up_success":                "Upload complete!",
	"up_err_dir_fail":           "Failed to upload directory %s: %v\n",
	"up_err_file_fail":          "Failed to upload file %s: %v\n",
}
