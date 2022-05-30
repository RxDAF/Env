package cfg

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Configure struct {
	MasterServer    string   `yaml:"masterServer"`
	Roles           []string `yaml:"roles"`           // 自己要运行的程序列表
	DownloadBufPath string   `yaml:"downloadBufPath"` // 可以为.
	ServicePath     string   `yaml:"servicePath"`
	Address         string   `yaml:"address"` //服务器的ip地址
}

func NewCFG(file string) (*Configure, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	cfg := new(Configure)
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
