package app

import (
	"autossh/src/i18n"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

var (
	Version string
	Build   string

	c       string
	v       bool
	h       bool
	upgrade bool
	cp      bool
	upload  bool
	up      bool
	down    bool
)

func init() {
	// 默认放在 ~/.config/autossh/config.yml
	home, _ := os.UserHomeDir()
	c = home + "/.config/autossh/config.yml"

	flag.StringVar(&c, "c", c, i18n.T("app_flag_c"))
	flag.StringVar(&c, "config", c, i18n.T("app_flag_c"))

	flag.BoolVar(&v, "v", v, i18n.T("app_flag_v"))
	flag.BoolVar(&v, "version", v, i18n.T("app_flag_v"))

	flag.BoolVar(&h, "h", h, i18n.T("app_flag_h"))
	flag.BoolVar(&h, "help", h, i18n.T("app_flag_h"))

	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) > 0 {
		arg := flag.Arg(0)
		switch arg {
		case "upgrade":
			upgrade = true
		case "cp":
			cp = true
		case "upload":
			upload = true
		case "up":
			up = true
		case "down":
			down = true
		default:
			defaultServer = arg
		}
	}
}

func Run() {
	// 监听系统信号，确保 Ctrl+C 始终有效
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		os.Exit(0)
	}()

	if v {
		showVersion()
	} else if h {
		showHelp()
	} else if upgrade {
		showUpgrade()
	} else if cp {
		showCp(c)
	} else if upload || up {
		showUpload(c)
	} else if down {
		showDownload(c)
	} else {
		showServers(c)
	}
}
