package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"maicai.ddxq.com/etcdv3"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
)

//https://godoc.org/github.com/coreos/etcd/clientv3#example-KV--Put

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", ping)
	router.GET("/getKeyList", getKeyList)
	router.GET("/getKeyList2", getKeyList2)
	router.GET("/getKey", getKey)
	router.GET("/setKey", setKey)

}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func setKey(c *gin.Context) {
	endpoints := []string{"localhost:2379"}
	etcdv3.InitConfig(endpoints, "1")

	key := c.DefaultQuery("key", "key2")
	value := c.DefaultQuery("value", "value_set_key")
	val, err := etcdv3.PutKey(key, value)
	if err == nil {
		fmt.Printf("%s =>%s set success", key, val)
		fmt.Println()
	} else {
		fmt.Errorf("get %s found error:%s", key, err.Error())
		fmt.Println()
	}

}

func getKey(c *gin.Context) {
	endpoints := []string{"localhost:2379"}
	etcdv3.InitConfig(endpoints, "1")

	key := c.DefaultQuery("key", "key2")
	val, err := etcdv3.GetKey(key)
	if err == nil {
		fmt.Printf("%s =>%s", key, val)
		fmt.Println()
	} else {
		fmt.Errorf("get %s found error:%s", key, err.Error())
		fmt.Println()
	}

}

func getKeyList(c *gin.Context) {
	endpoints := []string{"localhost:2379"}
	etcdv3.InitConfig(endpoints, "1")

	key := c.DefaultQuery("key", "batch_key")

	keys, err := etcdv3.GetKeyList(key, clientv3.WithPrefix())
	if err != nil {
		fmt.Errorf("getKeyList error:%s", err.Error())
	}
	for _, kv := range keys.Kvs {
		fmt.Printf("%s => %s", kv.Key, kv.Value)
		fmt.Println()
	}
}

func getKeyList2(c *gin.Context) {
	endpoints := []string{"localhost:2379"}
	etcdv3.InitConfig(endpoints, "1")
	key := c.DefaultQuery("key", "batch")
	fmt.Println("GetKeyListWithPrefix")
	keys, err := etcdv3.GetKeyListWithPrefix(key)
	for k, v := range keys {
		fmt.Printf("%s => %s", k, v)
		fmt.Println()
	}
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/keys", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
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
