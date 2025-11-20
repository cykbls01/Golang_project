package Model

type Node struct {
	NodeID      string `gorm:"column:node_id;type:varchar(120);primaryKey" json:"nodeID"`
	NodeIP      string `gorm:"column:node_ip;type:varchar(120);default:null" json:"nodeIP"`
	NodeName    string `gorm:"column:node_name;type:varchar(120);default:null" json:"nodeName"`
	ClusterID   string `gorm:"column:cluster_id;type:varchar(120);default;null" json:"clusterID"`
	ClusterName string `gorm:"column:cluster_name;type:varchar(120);default;null" json:"clusterName"`
}

func (n *Node) TableName() string {
	return "node"
}

func (n *Node) String() string {
	return n.NodeName
}
