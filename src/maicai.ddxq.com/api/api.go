package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	//"github.com/casbin/casbin"
	//"github.com/gin-contrib/authz"

	"github.com/gin-gonic/gin"
	cm "maicai.ddxq.com/manage"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
)

//https://godoc.org/github.com/coreos/etcd/clientv3#example-KV--Put

type KeyInfo struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MenuNodeAttributes struct {
	Url  string `json:"url"`
	Icon string `json:"icon"`
}

type MenuNode struct {
	Id         int                `json:"id"`
	Text       string             `json:"text"`
	State      string             `json:"state"`
	Attributes MenuNodeAttributes `json:"attributes"`
	MenuNodes  []*MenuNode        `json:"children,omitempty"`
}

type DataGrid struct {
	Total            int            `json:"total"`
	DataGridNodeList []DataGridNode `json:"rows"`
}
type DataGridNode struct {
	ProductId   string  `json:"productid,omitempty"`
	ProductName string  `json:"productname,omitempty"`
	UnitCost    float32 `json:"unitcost,omitempty"`
	Status      string  `json:"status,omitempty"`
	ListPrice   float32 `json:"listprice,omitempty"`
	Attr1       string  `json:"attr1,omitempty"`
	Itemid      string  `json:"itemid,omitempty"`
}

func setupRouter() *gin.Engine {
	//e := casbin.NewEnforcer("config/authz_model.conf", "config/authz_policy.csv")
	//router := gin.New()
	//router.Use(authz.NewAuthorizer(e))
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/css", "templates/css")
	router.Static("/easyui", "templates/easyui")
	router.Static("/js", "templates/js")
	router.Static("/images", "templates/images")
	router.Static("/temp", "templates/temp")

	router.GET("/ping", ping)
	router.GET("/getKeyList", getKeyListHandler)
	router.GET("/getKeyList2", getKeyList2Handler)
	router.GET("/getKey", getKeyHandler)
	router.GET("/setKey", setKeyHandler)

	router.GET("/index", indexHandler)
	router.POST("/getMenu", getMenuHandler)
	router.GET("/getMenu", getMenuHandler)
	router.POST("/getDataGrid", getDataGridHandler)
	router.GET("/getDataGrid", getDataGridHandler)
	return router
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Users",
	})
}

func getMenuHandler(c *gin.Context) {
	data, err := ioutil.ReadFile("./templates/temp/menu.json")
	//fmt.Printf("json:%s", data)
	if err != nil {
		fmt.Errorf("read error:%s", err.Error())
	}

	var menuList []MenuNode
	err = json.Unmarshal([]byte(data), &menuList)
	if err != nil {
		fmt.Errorf("json.Unmarshal error:%s", err.Error())
	}
	//fmt.Printf("menuList:%+v\n", menuList)
	c.JSON(http.StatusOK, menuList)
}

func getDataGridHandler(c *gin.Context) {
	data, err := ioutil.ReadFile("./templates/temp/datagrid.json")
	fmt.Printf("json:%s", data)
	if err != nil {
		fmt.Errorf("read error:%s", err.Error())
	}

	dataGrid := DataGrid{}
	err = json.Unmarshal([]byte(data), &dataGrid)
	if err != nil {
		fmt.Errorf("json.Unmarshal error:%s", err.Error())
	}
	//fmt.Printf("dataGrid:%+v\n", dataGrid)
	c.JSON(http.StatusOK, dataGrid)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func setKeyHandler(c *gin.Context) {
	key := c.DefaultQuery("key", "key2")
	value := c.DefaultQuery("value", "value_set_key")
	err := cm.SetKey(key, value)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "set success",
		})
	} else {

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Errorf("set failure, error:%s", err.Error()),
		})
	}
}

func getKeyHandler(c *gin.Context) {
	key := c.DefaultQuery("key", "key2")
	keyInfo := cm.GetKey(key)
	c.JSON(http.StatusOK, gin.H{
		"message": "getKeyHandler",
		"data":    keyInfo,
	})
}

func getKeyListHandler(c *gin.Context) {
	key := c.DefaultQuery("key", "batch_key")
	keyInfos := cm.GetKeyList(key)
	c.JSON(http.StatusOK, gin.H{
		"message": "getKeyListHandler",
		"data":    keyInfos,
	})
}

func getKeyList2Handler(c *gin.Context) {
	key := c.DefaultQuery("key", "batch")
	keyInfos := cm.GetKeyList2(key)
	c.JSON(http.StatusOK, gin.H{
		"message": "getKeyListHandler",
		"data":    keyInfos,
	})
}

func Setup() {
	router := setupRouter()
	router.Run()
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
