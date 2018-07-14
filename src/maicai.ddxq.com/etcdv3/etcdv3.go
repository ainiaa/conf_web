package etcdv3

import (
	"log"
	"time"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/clientv3"
)

var globalEndpoints []string
var globalTimeout string

var globalKVC clientv3.KV

func InitConfig(endpoints []string, timeout string) {
	globalEndpoints = endpoints
	globalTimeout = timeout
}

func getKVC() clientv3.KV {
	if globalKVC == nil {
		cfg := clientv3.Config{
			Endpoints: globalEndpoints,
			// set timeout per request to fail fast when the target endpoint is unavailable
			DialTimeout: time.Second,
		}
		c, err := clientv3.New(cfg)
		//client.EnablecURLDebug()
		if err != nil {
			log.Fatal(err)
		}
		globalKVC = clientv3.NewKV(c)
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

	ctx, canel := getCtx(globalTimeout)
	defer canel()
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
	ctx, canel := getCtx(globalTimeout)
	defer canel()
	if len(opts) > 0 {
		return kvc.Put(ctx, key, value, opts...)
	}
	return kvc.Put(ctx, key, value)
}

func DeleteKey(key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	kvc := getKVC()
	ctx, canel := getCtx(globalTimeout)
	defer canel()
	if len(opts) > 0 {
		return kvc.Delete(ctx, key, opts...)
	}
	return kvc.Delete(ctx, key)
}
