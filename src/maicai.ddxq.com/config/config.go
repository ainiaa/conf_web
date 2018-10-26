package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Etcdconfig etcd配置项
type Etcdconfig struct {
	Endpoints []string `yaml:"endpoints,flow"`
	Timeout   string   `yaml:"timeout"`
	Lease     int64    `yaml:"lease"`
}

// MyConf 自定义配置项
type MyConf struct {
	Etcdconfig      Etcdconfig `yaml:"etcdconfig"`
	ProjectStoreKey string     `yaml:"project_store_key"`
}

var config MyConf
var initConfig = false

// GetConfig 获取config内容
func GetConfig() MyConf {

	data, _ := ioutil.ReadFile("config/config.yml")
	fmt.Printf("config:%s", string(data))
	err := yaml.Unmarshal([]byte(data), &config)
	fmt.Printf("Unmarshal data:%v", config)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}
	fmt.Println("初始数据", config)
	fmt.Printf("config:%v", config.Etcdconfig)
	fmt.Println()
	fmt.Printf("timeout:%s", config.Etcdconfig.Timeout)
	initConfig = true
	return config
}

// GetEtcdConfig 获得etcd配置项
func GetEtcdConfig() Etcdconfig {
	if initConfig == false {
		config = GetConfig()
	}
	return config.Etcdconfig
}
