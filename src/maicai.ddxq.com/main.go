package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"maicai.ddxq.com/etcdv3"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
)

//https://godoc.org/github.com/coreos/etcd/clientv3#example-KV--Put

func main() {
	endpoints := []string{"localhost:2379"}
	etcdv3.InitConfig(endpoints, "1")
	/*
	key := "key2"
	val, err := etcdv3.GetKey(key)
	if err == nil {
		fmt.Printf("%s =>%s", key, val)
		fmt.Println()
	} else {
		fmt.Errorf("get %s found error:%s", key, err.Error())
	}*/
	/*for i:=0; i< 5;i++ {
		key := "batch_key:" + strconv.Itoa(i)
		val := key + ":value"
		etcdv3.PutKey(key,val)
	}*/
	keys,err := etcdv3.GetKeyList("batch_key", clientv3.WithPrefix())
	if err != nil {
		fmt.Errorf("getKeyList error:%s", err.Error())
	}
	for _,kv := range keys.Kvs {
		fmt.Printf("%s => %s", kv.Key, kv.Value)
		fmt.Println()
	}

	fmt.Println("GetKeyListWithPrefix")
	keys2,err := etcdv3.GetKeyListWithPrefix("batch")
	for k,v :=range keys2 {
		fmt.Printf("%s => %s", k, v)
		fmt.Println()
	}

}

func main1() {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Errorf("error:%s", err.Error())
	}
	defer cli.Close()

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	timeout, err := time.ParseDuration("10")
	if err == nil {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	key, value := "key1", "value1"
	resp, err := cli.Put(ctx, key, value)
	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v", err)
		default:
			log.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	} else {
		fmt.Printf("put resp:%s", resp)
		fmt.Println()
		resp, err := cli.Get(ctx, "key1")
		if err == nil {
			for _, kv := range resp.Kvs {
				fmt.Printf("get: key:%s => value:%s", kv.Key, kv.Value)
				fmt.Println()
			}
		} else {
			fmt.Printf("get error:%s", err.Error())
			fmt.Println()
		}
	}

}
