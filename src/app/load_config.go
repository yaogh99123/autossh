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
		return cfg, errors.New("Can't read configFile file:" + configFile)
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
