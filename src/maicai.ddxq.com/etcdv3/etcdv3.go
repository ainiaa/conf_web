package etcdv3

import (
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"maicai.ddxq.com/config"
)

var globalEndpoints = []string{"localhost:2379"}
var globalTimeout = "2"

var globalKVC clientv3.KV
var globalCli *clientv3.Client
var initedConfig = false

func InitGlobalConfig() {
	if initedConfig == false {
		c := config.GetEtcdConfig()
		globalEndpoints = c.Endpoints
		globalTimeout = c.Timeout
		initedConfig = true
	}
}

func InitConfig(endpoints []string, timeout string) {
	globalEndpoints = endpoints
	globalTimeout = timeout
}

func getCli() *clientv3.Client {
	if globalCli == nil {
		cfg := clientv3.Config{
			Endpoints:   globalEndpoints,
			DialTimeout: time.Second,
		}
		cli, err := clientv3.New(cfg)
		if err != nil {
			log.Fatal(err)
		}
		globalCli = cli
	}
	return globalCli
}

func getKVC() clientv3.KV {
	if globalKVC == nil {
		globalCli = getCli()
		globalKVC = clientv3.NewKV(globalCli)
	}
	return globalKVC
}

func getCtx(duration string) (context.Context, context.CancelFunc) {

	timeout, err := time.ParseDuration(duration)
	var ctx context.Context
	var cancel context.CancelFunc
	if err == nil {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	return ctx, cancel
}

func GetKey(key string) (string, error) {
	gresp, err := GetKeyList(key)
	if err == nil {
		if len(gresp.Kvs) > 0 {
			kv := gresp.Kvs[0]
			return string(kv.Value), err
		}
	}
	return "", err
}

func GetKeyRev(key string, rev int64) (string, error) {
	gresp, err := GetKeyList(key, clientv3.WithRev(rev))
	if err == nil {
		if len(gresp.Kvs) > 0 {
			kv := gresp.Kvs[0]
			return string(kv.Value), err
		}
	}
	return "", err
}

func GetKeyList(key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	kvc := getKVC()

	ctx, cancel := getCtx(globalTimeout)
	defer cancel()
	if len(opts) > 0 {
		return kvc.Get(ctx, key, opts...)
	}

	return kvc.Get(ctx, key)
}

func GetKeyListWithPrefix(key string) (map[string]string, error) {
	gresp, err := GetKeyList(key, clientv3.WithPrefix())
	m := make(map[string]string)
	if err == nil {
		for _, kv := range gresp.Kvs {
			m[string(kv.Key)] = string(kv.Value)
		}
	}
	return m, nil
}

func PutKey(key, value string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	kvc := getKVC()
	ctx, cancel := getCtx(globalTimeout)
	defer cancel()
	if len(opts) > 0 {
		return kvc.Put(ctx, key, value, opts...)
	}
	return kvc.Put(ctx, key, value)
}

func DeleteKey(key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	kvc := getKVC()
	ctx, cancel := getCtx(globalTimeout)
	defer cancel()
	if len(opts) > 0 {
		return kvc.Delete(ctx, key, opts...)
	}
	return kvc.Delete(ctx, key)
}
