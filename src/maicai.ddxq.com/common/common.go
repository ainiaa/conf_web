package common

// KeyInfo kv结构
type KeyInfo struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// MenuNodeAttributes 菜单节点属性结构
// Url 菜单节点的url
// Icon 菜单节点的icon
type MenuNodeAttributes struct {
	Url  string `json:"url"`
	Icon string `json:"icon"`
}

// MenuNode 菜单节点结构
// ID 菜单id属性
// Text 菜单标签
// State 菜单状态
// Attributes 菜单属性
// MenuNodes 菜单节点
type MenuNode struct {
	ID         int                `json:"id"`
	Text       string             `json:"text"`
	State      string             `json:"state"`
	Attributes MenuNodeAttributes `json:"attributes"`
	MenuNodes  []*MenuNode        `json:"children,omitempty"`
}

// DataGrid DataGrid属性
// Total DataGrid的总量
// DataGridNodeList 节点属性
type DataGrid struct {
	Total            int            `json:"total"`
	DataGridNodeList []DataGridNode `json:"rows"`
}

// DataGridNode 节点结构
// ProductID 产品id
// ProductName 产品名称
// UnitCost 产品单价
// Status 产品状态
// ListPrice 产品列表价格
// Attr1 属性1
// Itemid 属性id
type DataGridNode struct {
	ProductID   string  `json:"productid,omitempty"`
	ProductName string  `json:"productname,omitempty"`
	UnitCost    float32 `json:"unitcost,omitempty"`
	Status      string  `json:"status,omitempty"`
	ListPrice   float32 `json:"listprice,omitempty"`
	Attr1       string  `json:"attr1,omitempty"`
	ItemID      string  `json:"itemid,omitempty"`
}
