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
}
