package etcdv3

import (
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"maicai.ddxq.com/config"
)

var globalEndpoints = []string{"localhost:2379"}
var globalTimeout = "2"
var globalLease int64 = 5

var globalKVC clientv3.KV
var globalCli *clientv3.Client
var initedConfig = false

//LeaseID LeaseID
type LeaseID = clientv3.LeaseID

// InitGlobalConfig 初始化公共配置项
func InitGlobalConfig() {
	if initedConfig == false {
		c := config.GetEtcdConfig()
		globalEndpoints = c.Endpoints
		globalTimeout = c.Timeout
		globalLease = c.Lease
		initedConfig = true
	}
}

// InitConfig 初始化配置项
func InitConfig(endpoints []string, timeout string) {
	globalEndpoints = endpoints
	globalTimeout = timeout
}

func getCli() (*clientv3.Client, error) {
	var err error
	var cli *clientv3.Client
	if globalCli == nil {
		cfg := clientv3.Config{
			Endpoints:   globalEndpoints,
			DialTimeout: time.Second,
		}
		cli, err = clientv3.New(cfg)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		globalCli = cli
	}
	return globalCli, err
}

func getKVC() (clientv3.KV, error) {
	var err error
	if globalKVC == nil {
		globalCli, err = getCli()
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		globalKVC = clientv3.NewKV(globalCli)
	}
	return globalKVC, err
}

func getCtx(duration string) (context.Context, context.CancelFunc) {

	timeout, err := time.ParseDuration(duration)
	if err != nil {
		log.Fatal(err)
	}
	var ctx context.Context
	var cancel context.CancelFunc
	if err == nil {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	return ctx, cancel
}

// GetKey 获取key的相关信息
func GetKey(key string) (string, error) {
	gresp, err := GetKeyList(key)
	if err != nil {
		return "", err
	}
	if len(gresp.Kvs) > 0 {
		kv := gresp.Kvs[0]
		return string(kv.Value), err
	}
	return "", err
}

// GetKeyRev 获取key的相关信息（包含版本信息）
func GetKeyRev(key string, rev int64) (string, error) {
	gresp, err := GetKeyList(key, clientv3.WithRev(rev))
	if err != nil {
		return "", err
	}
	if len(gresp.Kvs) > 0 {
		kv := gresp.Kvs[0]
		return string(kv.Value), err
	}
	return "", err
}

// GetKeyList 获取key列表
func GetKeyList(key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	kvc, err := getKVC()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	ctx, cancel := getCtx(globalTimeout)
	defer cancel()
	if len(opts) > 0 {
		return kvc.Get(ctx, key, opts...)
	}

	return kvc.Get(ctx, key)
}

// GetKeyListWithPrefix 根据前缀获取key列表
func GetKeyListWithPrefix(key string) (map[string]string, error) {
	gresp, err := GetKeyList(key, clientv3.WithPrefix())
	m := make(map[string]string)
	if err != nil {
		//reutrn m,err
	}
	for _, kv := range gresp.Kvs {
		m[string(kv.Key)] = string(kv.Value)
	}
	return m, nil
}

// PutKeyWithLease 新增key （带lease）
func PutKeyWithLease(key, val string, ttl int64) (*clientv3.LeaseGrantResponse, *clientv3.PutResponse, error) {
	cli, err := getCli()
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	//defer cli.Close()
	ctx, cancel := getCtx(globalTimeout)
	defer cancel()
	resp, err := cli.Grant(ctx, ttl)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	fmt.Printf("resp.ID:%d ", resp.ID)
	fmt.Println()
	r, err := cli.Put(ctx, key, val, clientv3.WithLease(resp.ID))
	return resp, r, err
}

// LeaseRevoke 删除租约
func LeaseRevoke(respID LeaseID) error {
	cli, err := getCli()
	if err != nil {
		log.Fatal(err)
		return err
	}
	//defer cli.Close()
	ctx, cancel := getCtx(globalTimeout)
	defer cancel()

	_, err = cli.Revoke(ctx, respID)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}

// LeaseKeepAlive  取消租约
func LeaseKeepAlive(respID LeaseID) error {
	cli, err := getCli()
	if err != nil {
		log.Fatal(err)
		return err
	}
	//defer cli.Close()
	ctx, cancel := getCtx(globalTimeout)
	defer cancel()
	_, err = cli.KeepAlive(ctx, respID)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}

// LeaseKeepAliveOnce 续约
func LeaseKeepAliveOnce(respID LeaseID) error {
	cli, err := getCli()
	if err != nil {
		log.Fatal(err)
		return err
	}
	//defer cli.Close()
	ctx, cancel := getCtx(globalTimeout)
	defer cancel()
	_, err = cli.KeepAliveOnce(ctx, respID)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}

// TimeToLive 获得TTL相关信息
func TimeToLive(respID LeaseID) (*clientv3.LeaseTimeToLiveResponse, error) {
	cli, err := getCli()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	//defer cli.Close()
	ctx, cancel := getCtx(globalTimeout)
	defer cancel()
	lresp, err := cli.TimeToLive(ctx, respID)
	if err != nil {
		log.Fatal(err)
		return lresp, err
	}
	return lresp, err
}

// PutKey 新增key
func PutKey(key, value string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	kvc, err := getKVC()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	ctx, cancel := getCtx(globalTimeout)
	defer cancel()
	if len(opts) > 0 {
		return kvc.Put(ctx, key, value, opts...)
	}
	return kvc.Put(ctx, key, value)
}

// DeleteKey 删除key
func DeleteKey(key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	kvc, err := getKVC()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	ctx, cancel := getCtx(globalTimeout)
	defer cancel()
	if len(opts) > 0 {
		return kvc.Delete(ctx, key, opts...)
	}
	return kvc.Delete(ctx, key)
}
