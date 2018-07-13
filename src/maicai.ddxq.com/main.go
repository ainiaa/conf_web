package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
)

func main() {

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
