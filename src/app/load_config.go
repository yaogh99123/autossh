package app

import (
	"autossh/src/utils"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// 加载配置
func loadConfig(configFile string) (cfg *Config, err error) {
	configFile, err = utils.ParsePath(configFile)
	if err != nil {
		return cfg, err
	}

	if exists, _ := utils.FileIsExists(configFile); !exists {
		utils.Infoln("配置文件不存在，正在初始化默认配置...")
		cfg, err = initConfig(configFile)
		if err != nil {
			return nil, errors.Wrap(err, "初始化配置失败")
		}
		return cfg, nil
	}

	b, _ := ioutil.ReadFile(configFile)
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return cfg, err
	}

	cfg.file = configFile
	cfg.createServerIndex()

	return cfg, nil
}

// 初始化默认配置
func initConfig(configFile string) (*Config, error) {
	cfg := &Config{
		ShowDetail: true,
		Options: map[string]interface{}{
			"ServerAliveInterval": 30,
			"TERM":                "xterm-256color",
		},
		Servers: []*Server{
			{
				Name:     "Demo-Password",
				Ip:       "127.0.0.1",
				Port:     22,
				User:     "root",
				Password: "password",
				Method:   "password",
				Alias:    "demo",
			},
		},
		Groups: []*Group{
			{
				GroupName: "Default Group",
				Prefix:    "g",
				Servers: []Server{
					{
						Name:     "Group-Demo",
						Ip:       "127.0.0.1",
						User:     "root",
						Password: "password",
					},
				},
			},
		},
	}
	cfg.file = configFile
	err := cfg.saveConfig(false)
	if err != nil {
		return nil, err
	}
	cfg.createServerIndex()
	return cfg, nil
}
