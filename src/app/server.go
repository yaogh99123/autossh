package app

import (
	"autossh/src/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/proxy"
)

type Server struct {
	Name     string                 `json:"name" yaml:"name"`
	Ip       string                 `json:"ip" yaml:"ip"`
	Port     int                    `json:"port" yaml:"port"`
	User     string                 `json:"user" yaml:"user"`
	Password string                 `json:"password" yaml:"password"`
	Method   string                 `json:"method" yaml:"method"`
	Key      string                 `json:"key" yaml:"key"`
	Options  map[string]interface{} `json:"options" yaml:"options"`
	Alias    string                 `json:"alias" yaml:"alias"`
	NoProxy  bool                   `json:"no_proxy" yaml:"no_proxy"`
	Log      *ServerLog             `json:"log,omitempty" yaml:"log,omitempty"`

	termWidth  int
	termHeight int
	groupName  string
	group      *Group
}

// 格式化，赋予默认值
func (server *Server) Format() {
	if server.Port == 0 {
		server.Port = 22
	}

	if server.Method == "" {
		server.Method = "password"
	}
}

// 合并选项
func (server *Server) MergeOptions(options map[string]interface{}, overwrite bool) {
	if server.Options == nil {
		server.Options = make(map[string]interface{})
	}

	for k, v := range options {
		if overwrite {
			server.Options[k] = v
		} else {
			if _, ok := server.Options[k]; !ok {
				server.Options[k] = v
			}
		}

	}
}

// 格式化输出，用于打印
func (server *Server) FormatPrint(flag string, ShowDetail bool, flagWidth, nameWidth int) string {
	alias := ""
	if server.Alias != "" {
		alias = "|" + server.Alias
	}

	fullFlag := "[" + flag + alias + "]"
	flagPart := utils.AppendRight(fullFlag, " ", flagWidth)
	coloredFlag := utils.Colored(flagPart, utils.ColorCyan)

	namePart := utils.AppendRight(server.Name, " ", nameWidth)

	method := server.Method
	if method == "" {
		method = "password"
	}
	methodPart := utils.Colored("["+method+"]", utils.ColorYellow)

	if ShowDetail {
		info := "[" + server.User + "@" + server.Ip + "]"
		infoPart := utils.AppendRight(info, " ", 35) // 补齐到 35 字符宽度
		return " " + coloredFlag + "  " + namePart + "  " + infoPart + "  " + methodPart
	} else {
		// 非详细模式下，我们也补齐 35 个空位，确保 Method 依然对齐
		emptyPart := utils.AppendRight("", " ", 35)
		return " " + coloredFlag + "  " + namePart + "  " + emptyPart + "  " + methodPart
	}
}

// 生成SSH Client
func (server *Server) GetSshClient() (*ssh.Client, error) {
	auth, err := parseAuthMethods(server)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: server.User,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// 默认端口为22
	if server.Port == 0 {
		server.Port = 22
	}

	addr := server.Ip + ":" + strconv.Itoa(server.Port)

	// 如果设置了不走代理，则强制直连
	if server.NoProxy {
		// 1. 使用 proxy.Direct 显式执行直连，由于是外部 Dial，它会绕过 Go 内部对环境变量的缓存逻辑
		dialer := proxy.Direct
		conn, err := dialer.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}

		// 2. 将连接升级为 SSH，手动指定 Handshake 超时
		config.Timeout = 15 * time.Second
		c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
		if err != nil {
			return nil, err
		}
		return ssh.NewClient(c, chans, reqs), nil
	}

	if server.group != nil && server.group.Proxy != nil {
		return server.proxySshClient(server.group.Proxy, addr, config)
	} else {
		return ssh.Dial("tcp", addr, config)
	}
}

func (server *Server) proxySshClient(p *Proxy, sshServerAddr string, sshConfig *ssh.ClientConfig) (client *ssh.Client, err error) {
	var dialer proxy.Dialer
	switch p.Type {
	case ProxyTypeSocks5:
		var auth proxy.Auth
		if p.User != "" {
			auth = proxy.Auth{
				User:     p.User,
				Password: p.Password,
			}
		}

		dialer, err = proxy.SOCKS5("tcp", p.Server+":"+strconv.Itoa(p.Port), &auth, proxy.Direct)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("unknown proxy type: %s", p.Type))
	}

	conn, err := dialer.Dial("tcp", sshServerAddr)
	if err != nil {
		return nil, err
	}

	c, chans, reqs, err := ssh.NewClientConn(conn, sshServerAddr, sshConfig)
	if err != nil {
		return nil, err
	}

	return ssh.NewClient(c, chans, reqs), nil
}

// 生成Sftp Client
func (server *Server) GetSftpClient() (*sftp.Client, error) {
	sshClient, err := server.GetSshClient()
	if err == nil {
		return sftp.NewClient(sshClient)
	} else {
		return nil, err
	}
}

// 执行远程连接
func (server *Server) Connect(globalLog ServerLog) error {
	client, err := server.GetSshClient()
	if err != nil {
		if utils.ErrorAssert(err, "ssh: unable to authenticate") {
			return errors.New("连接失败，请检查密码/密钥是否有误")
		}

		return errors.New("ssh dial fail:" + err.Error())
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return errors.New("create session fail:" + err.Error())
	}

	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return errors.New("创建文件描述符出错:" + err.Error())
	}
	defer terminal.Restore(fd, oldState)

	stopKeepAliveLoop := server.startKeepAliveLoop(session)
	defer close(stopKeepAliveLoop)

	err = server.stdIO(session, globalLog)
	if err != nil {
		return err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	server.termWidth, server.termHeight, _ = terminal.GetSize(fd)
	termType := ""
	if val, ok := server.Options["TERM"]; ok && val != nil {
		termType = val.(string)
	}
	if termType == "" {
		termType = os.Getenv("TERM")
	}
	if termType == "" || termType == "xterm-ghostty" {
		termType = "xterm-256color"
	}
	if err := session.RequestPty(termType, server.termHeight, server.termWidth, modes); err != nil {
		return errors.New("创建终端出错:" + err.Error())
	}

	server.listenWindowChange(session, fd)

	err = session.Shell()
	if err != nil {
		return errors.New("执行Shell出错:" + err.Error())
	}

	_ = session.Wait()
	//if err != nil {
	//	return errors.New("执行Wait出错:" + err.Error())
	//}

	return nil
}

// 重定向标准输入输出
func (server *Server) stdIO(session *ssh.Session, globalLog ServerLog) error {
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	logConfig := globalLog
	if server.Log != nil {
		logConfig = *server.Log
	}

	if logConfig.Enable {
		ch, err := session.StdoutPipe()
		if err != nil {
			return err
		}

		go func() {
			flag := os.O_RDWR | os.O_CREATE
			switch logConfig.Mode {
			case LogModeAppend:
				flag = flag | os.O_APPEND
			case LogModeCover:
			}

			f, err := os.OpenFile(server.formatLogFilename(logConfig.Filename), flag, 0644)
			if err != nil {
				utils.Logger.Error("Open file fail ", err)
				return
			}

			for {
				buff := [4096]byte{}
				n, _ := ch.Read(buff[:])
				if n > 0 {
					if _, err := f.Write(buff[:n]); err != nil {
						utils.Logger.Error("Write file buffer fail ", err)
					}

					if _, err := os.Stdout.Write(buff[:n]); err != nil {
						utils.Logger.Error("Write stdout buffer fail ", err)
					}
				}
			}
		}()
	} else {
		session.Stdout = os.Stdout
	}

	return nil
}

// 格式化日志文件名
func (server *Server) formatLogFilename(filename string) string {
	kvs := []map[string]string{
		{"%g": server.groupName},
		{"%n": server.Name},
		{"%dt": time.Now().Format("20060102150405")},
		{"%d": time.Now().Format("20060102")},
		{"%u": server.User},
		{"%a": server.Alias},
	}

	for _, kv := range kvs {
		for k, v := range kv {
			filename = strings.ReplaceAll(filename, k, v)
		}
	}

	return filename
}

// 解析鉴权方式
func parseAuthMethods(server *Server) ([]ssh.AuthMethod, error) {
	var sshs []ssh.AuthMethod

	switch strings.ToLower(server.Method) {
	case "password":
		sshs = append(sshs, ssh.Password(server.Password))
		break

	case "key":
		method, err := pemKey(server)
		if err != nil {
			return nil, err
		}
		sshs = append(sshs, method)
		break

		// 默认以password方式
	default:
		sshs = append(sshs, ssh.Password(server.Password))
	}

	return sshs, nil
}

// 解析密钥
func pemKey(server *Server) (ssh.AuthMethod, error) {
	if server.Key == "" {
		server.Key = "~/.ssh/id_rsa"
	}
	server.Key, _ = utils.ParsePath(server.Key)

	pemBytes, err := ioutil.ReadFile(server.Key)
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer
	if server.Password == "" {
		signer, err = ssh.ParsePrivateKey(pemBytes)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(server.Password))
	}

	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}

// 发送心跳包
func (server *Server) startKeepAliveLoop(session *ssh.Session) chan struct{} {
	terminate := make(chan struct{})
	go func() {
		for {
			select {
			case <-terminate:
				return
			default:
				if val, ok := server.Options["ServerAliveInterval"]; ok && val != nil {
					_, err := session.SendRequest("keepalive@bbr", true, nil)
					if err != nil {
						utils.Logger.Category("server").Error("keepAliveLoop fail", err)
					}

					var interval int
					switch v := server.Options["ServerAliveInterval"].(type) {
					case int:
						interval = v
					case float64:
						interval = int(v)
					case int64:
						interval = int(v)
					default:
						interval = 30
					}
					t := time.Duration(interval)
					time.Sleep(time.Second * t)
				} else {
					return
				}
			}
		}
	}()
	return terminate
}

// 监听终端窗口变化
func (server *Server) listenWindowChange(session *ssh.Session, fd int) {
	go func() {
		sigwinchCh := make(chan os.Signal, 1)
		defer close(sigwinchCh)

		signal.Notify(sigwinchCh, syscall.SIGWINCH)
		termWidth, termHeight, err := terminal.GetSize(fd)
		if err != nil {
			utils.Logger.Error(err)
		}

		for {
			select {
			// 阻塞读取
			case sigwinch := <-sigwinchCh:
				if sigwinch == nil {
					return
				}
				currTermWidth, currTermHeight, err := terminal.GetSize(fd)

				// 判断一下窗口尺寸是否有改变
				if currTermHeight == termHeight && currTermWidth == termWidth {
					continue
				}

				// 更新远端大小
				session.WindowChange(currTermHeight, currTermWidth)
				if err != nil {
					utils.Logger.Error(err)
					continue
				}

				termWidth, termHeight = currTermWidth, currTermHeight
			}
		}
	}()
}
