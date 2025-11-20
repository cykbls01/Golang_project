package Model

import (
	"time"
)

// Pod 对应Django的Pod模型，与数据库表映射
type Pod struct {
	ClusterID      string    `gorm:"column:cluster_id;type:varchar(120);default:null" json:"clusterID"`            // 取自传递参数
	ClusterName    string    `gorm:"column:cluster_name;type:varchar(120);default:null" json:"clusterName"`        // 取自传递参数
	NodeID         string    `gorm:"column:node_id;type:varchar(120);default:null" json:"nodeID"`                  // 节点ID
	NodeIP         string    `gorm:"column:node_ip;type:varchar(120);default:null" json:"nodeIP"`                  // 节点IP
	WorkloadName   string    `gorm:"column:workload_name;type:varchar(120);default:null" json:"workloadName"`      // 工作负载名称
	Kind           string    `gorm:"column:kind;type:varchar(120);default:null" json:"kind"`                       // 无状态、有状态、守护集等，取自workload表
	WorkloadID     string    `gorm:"column:workload_id;type:varchar(120);default:null" json:"workloadID"`          // 属于哪个service（注：字段注释与名称可能不一致，保持原逻辑）
	PodID          string    `gorm:"column:pod_id;type:varchar(120);primaryKey" json:"podID"`                      // 主键
	PodName        string    `gorm:"column:pod_name;type:varchar(120);default:null" json:"podName"`                // Pod名称
	PodIP          string    `gorm:"column:pod_ip;type:varchar(120);default:null" json:"podIP"`                    // PodIP
	Namespace      string    `gorm:"column:namespace;type:varchar(120);default:null" json:"namespace"`             // 命名空间
	ContainersName string    `gorm:"column:containers_name;type:varchar(1200);default:null" json:"containersName"` // 容器名称（多个用分隔符）
	ImagesVersion  string    `gorm:"column:images_version;type:varchar(1200);default:null" json:"imagesVersion"`   // 镜像版本（多个用分隔符）
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime;default:null" json:"create_time"`            // 创建时间（自动设置）
	UpdateTime     time.Time `gorm:"column:update_time;autoUpdateTime;default:null" json:"update_time"`            // 更新时间（自动更新）
	SelectTime     time.Time `gorm:"column:select_time;autoCreateTime" json:"select_time"`                         // 选择时间（自动设置）
	CPU            string    `gorm:"column:cpu;type:varchar(120);default:null" json:"cpu"`                         // CPU配置
	Memory         string    `gorm:"column:memory;type:varchar(120);default:null" json:"memory"`                   // 内存配置
	Status         string    `gorm:"column:status;type:varchar(120);default:null" json:"status"`                   // 是否在运行（状态字段）
	Region         string    `gorm:"column:region;type:varchar(120);default:null" json:"region"`                   // 取自传递参数
}

// TableName 指定数据库表名（默认会转为复数pods，如需指定原表名可修改）
func (p *Pod) TableName() string {
	return "pod" // 若Django表名是pods，可改为"pods"，根据实际数据库表名调整
}

// String 实现Stringer接口，对应Django的__str__方法
func (p *Pod) String() string {
	return p.PodName
}
