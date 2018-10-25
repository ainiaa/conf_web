package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Etcdconfig struct {
	Endpoints []string `yaml:"endpoints,flow"`
	Timeout   string   `yaml:"timeout"`
	Lease     int64    `yaml:"lease"`
}
type MyConf struct {
	Etcdconfig      Etcdconfig `yaml:"etcdconfig"`
	ProjectStoreKey string     `yaml:"project_store_key"`
}

var config MyConf
var initConfig = false

func GetConfig() MyConf {
	data, _ := ioutil.ReadFile("config/config.yml")
	fmt.Printf("config:%s", string(data))
	err := yaml.Unmarshal([]byte(data), &config)
	fmt.Printf("Unmarshal data:%v", config)
	if err != nil {
		fmt.Errorf("Unmarshal: %v", err)
	}
	fmt.Println("初始数据", config)
	fmt.Printf("config:%v", config.Etcdconfig)
	fmt.Println()
	fmt.Printf("timeout:%d", config.Etcdconfig.Timeout)
	initConfig = true
	return config
}

func GetEtcdConfig() Etcdconfig {
	if initConfig == false {
		config = GetConfig()
	}
	return config.Etcdconfig
}
