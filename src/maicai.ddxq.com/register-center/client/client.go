package register

import (
	"fmt"

	etcd "maicai.ddxq.com/etcdv3"
)

func register() error {
	etcd.PutKey()
}

func keepAlive() error {

}
