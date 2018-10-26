package manage

import (
	"fmt"

	"github.com/coreos/etcd/clientv3"
	"maicai.ddxq.com/common"
	"maicai.ddxq.com/etcdv3"
	"maicai.ddxq.com/util"
)

//https://godoc.org/github.com/coreos/etcd/clientv3#example-KV--Put

// SetKey 设置key
// key key名称
// value key值
func SetKey(key string, value string) error {

	etcdv3.InitGlobalConfig()
	val, err := etcdv3.PutKey(key, value)
	if err == nil {
		fmt.Printf("%s =>%s set success", key, val)
		fmt.Println()
	} else {
		fmt.Printf("get %s found error:%s", key, err.Error())
		fmt.Println()
	}
	return err

}

// SetKeyWithLease 设置key（带有lease）
// key key名称
// value key值
// ttl 租约期
func SetKeyWithLease(key, value string, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	etcdv3.InitGlobalConfig()
	lresp, _, err := etcdv3.PutKeyWithLease(key, value, ttl)
	if err == nil {
		fmt.Printf("%s =>%s set success", key, value)
		fmt.Println()
	} else {
		fmt.Printf("get %s found error:%s", key, err.Error())
		fmt.Println()
	}
	fmt.Printf("lrsp.ID:%d \r\n", lresp.ID)
	fmt.Printf("lrsp:%+v", lresp)
	fmt.Println()
	return lresp, err
}

// LeaseRevoke 删除租约
func LeaseRevoke(leaseID etcdv3.LeaseID) error {
	return etcdv3.LeaseRevoke(leaseID)
}

// LeaseKeepAlive 取消租约为永久
func LeaseKeepAlive(leaseID etcdv3.LeaseID) error {
	return etcdv3.LeaseKeepAlive(leaseID)
}

// LeaseKeepAliveOnce 续约
func LeaseKeepAliveOnce(leaseID etcdv3.LeaseID) error {
	return etcdv3.LeaseKeepAliveOnce(leaseID)
}

// TimetoLive 获得租约信息
func TimetoLive(leaseID etcdv3.LeaseID) (*clientv3.LeaseTimeToLiveResponse, error) {
	return etcdv3.TimeToLive(leaseID)
}

// GetKey 获取key的相关信息
func GetKey(key string) common.KeyInfo {
	etcdv3.InitGlobalConfig()
	val, err := etcdv3.GetKey(key)
	if err == nil {
		fmt.Printf("%s =>%s", key, val)
		fmt.Println()
	} else {
		fmt.Printf("get %s found error:%s", key, err.Error())
		fmt.Println()
	}
	return common.KeyInfo{Key: key, Value: val}
}

// GetKeyList 获取key的列表信息
func GetKeyList(key string) []common.KeyInfo {
	etcdv3.InitGlobalConfig()
	//keys, err := etcdv3.GetKeyList(key, clientv3.WithPrefix())
	keys, err := etcdv3.GetKeyList(key)
	if err != nil {
		fmt.Printf("getKeyList error:%s", err.Error())
	}
	keyInfos := make([]common.KeyInfo, 0)
	for _, kv := range keys.Kvs {
		fmt.Printf("%s => %s", kv.Key, kv.Value)
		fmt.Println()
		keyInfo := common.KeyInfo{Key: util.ToString(kv.Key), Value: util.ToString(kv.Value)}
		keyInfos = append(keyInfos, keyInfo)
	}

	return keyInfos
}

// GetKeyList2 获取key的列表信息
func GetKeyList2(key string) []common.KeyInfo {
	etcdv3.InitGlobalConfig()
	fmt.Println("GetKeyListWithPrefix")
	keys, err := etcdv3.GetKeyListWithPrefix(key)
	if err != nil {
		fmt.Printf("etcdv3.GetKeyListWithPrefix found error, key:%s, error:%s", key, err.Error())
		fmt.Println()
	}
	keyInfos := make([]common.KeyInfo, 0)
	for k, v := range keys {
		fmt.Printf("%s => %s", k, v)
		fmt.Println()
		keyInfos = append(keyInfos, common.KeyInfo{Key: k, Value: v})
	}

	return keyInfos
}
