package manage

import (
	"fmt"

	"maicai.ddxq.com/etcdv3"
	"maicai.ddxq.com/util"
)

//https://godoc.org/github.com/coreos/etcd/clientv3#example-KV--Put

type KeyInfo struct {
	Key   string `json:"key" form:"key"`
	Value string `json:"value" form:"value"`
}

func SetKey(key string, value string) error {

	etcdv3.InitGlobalConfig()
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

func SetKeyWithLease(key, value string, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	etcdv3.InitGlobalConfig()
	lresp, presp, err := etcdv3.PutKeyWithLease(key, value, ttl)
	if err == nil {
		fmt.Printf("%s =>%s set success", key, value)
		fmt.Println()
	} else {
		fmt.Errorf("get %s found error:%s", key, err.Error())
		fmt.Println()
	}
	return lresp, err
}

func LeaseRevoke(respID etcdv3.RespID) error {
	return etcdv3.LeaseRevoke(respID)
}

func LeaseKeepAlive(respID etcdv3.RespID) error {
	return etcdv3.LeaseKeepAlive(respID)
}

func LeaseKeepAliveOnce(respID etcdv3.RespID) error {
	return etcdv3.LeaseKeepAliveOnce(respID)
}

func GetKey(key string) KeyInfo {
	etcdv3.InitGlobalConfig()
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

func GetKeyList(key string) []KeyInfo {
	etcdv3.InitGlobalConfig()
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

func GetKeyList2(key string) []KeyInfo {
	etcdv3.InitGlobalConfig()
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
