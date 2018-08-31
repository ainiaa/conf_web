package config_manage

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"maicai.ddxq.com/config"
	"maicai.ddxq.com/etcdv3"
	"maicai.ddxq.com/util"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
)

//https://godoc.org/github.com/coreos/etcd/clientv3#example-KV--Put

type KeyInfo struct {
	Key   string `json:"key" form:"key"`
	Value string `json:"value" form:"value"`
}

func setKey(key string, value string) error {

	//endpoints := []string{"localhost:2379"}
	//etcdv3.InitConfig(endpoints, "1")
	etcdconfig := config.GetEtcdConfig()
	etcdv3.InitConfig(etcdconfig.Endpoints, etcdconfig.Timeout)

	val, err := etcdv3.PutKey(key, value)
	if err == nil {
		fmt.Printf("%s =>%s set success", key, val)
		fmt.Println()
	} else {
		fmt.Errorf("get %s found error:%s", key, err.Error())
		fmt.Println()
	}
	return err

}

func getKey(key string) KeyInfo {
	endpoints := []string{"localhost:2379"}
	etcdv3.InitConfig(endpoints, "1")

	val, err := etcdv3.GetKey(key)
	if err == nil {
		fmt.Printf("%s =>%s", key, val)
		fmt.Println()
	} else {
		fmt.Errorf("get %s found error:%s", key, err.Error())
		fmt.Println()
	}
	return KeyInfo{key, val}
}

func getKeyList(key string) []KeyInfo {
	endpoints := []string{"localhost:2379"}
	etcdv3.InitConfig(endpoints, "1")

	//keys, err := etcdv3.GetKeyList(key, clientv3.WithPrefix())
	keys, err := etcdv3.GetKeyList(key)
	if err != nil {
		fmt.Errorf("getKeyList error:%s", err.Error())
	}
	keyInfos := make([]KeyInfo, 0)
	for _, kv := range keys.Kvs {
		fmt.Printf("%s => %s", kv.Key, kv.Value)
		fmt.Println()
		keyInfo := KeyInfo{Key: util.ToString(kv.Key), Value: util.ToString(kv.Value)}
		keyInfos = append(keyInfos, keyInfo)
	}

	return keyInfos
}

func getKeyList2(key string) []KeyInfo {
	endpoints := []string{"localhost:2379"}
	etcdv3.InitConfig(endpoints, "1")
	fmt.Println("GetKeyListWithPrefix")
	keys, err := etcdv3.GetKeyListWithPrefix(key)
	if err != nil {
		fmt.Printf("etcdv3.GetKeyListWithPrefix found error, key:%s, error:%s", key, err.Error())
		fmt.Println()
	}
	keyInfos := make([]KeyInfo, 0)
	for k, v := range keys {
		fmt.Printf("%s => %s", k, v)
		fmt.Println()
		keyInfos = append(keyInfos, KeyInfo{k, v})
	}

	return keyInfos
}
