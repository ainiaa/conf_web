package main

import (
	"fmt"

	"maicai.ddxq.com/api"
)

func main() {

	fmt.Println("start")
	api.Setup()
	//config := config.GetEtcdConfig()
	//fmt.Printf("etcdConfig:%v", config)
}
