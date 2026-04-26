package app

import (
	"flag"
	"os"
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

	flag.StringVar(&c, "c", c, "指定配置文件路径")
	flag.StringVar(&c, "config", c, "指定配置文件路径")

	flag.BoolVar(&v, "v", v, "版本信息")
	flag.BoolVar(&v, "version", v, "版本信息")

	flag.BoolVar(&h, "h", h, "帮助信息")
	flag.BoolVar(&h, "help", h, "帮助信息")

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
