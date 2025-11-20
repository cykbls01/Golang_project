package Model

import "time"

type Service struct {
	ClusterID    string    `gorm:"column:cluster_id;type:varchar(120);default;null" json:"clusterID"`
	ClusterName  string    `gorm:"column:cluster_name;type:varchar(120);default;null" json:"clusterName"`
	Namespace    string    `gorm:"column:namespace;type:varchar(120);default;null" json:"namespace"`
	ServiceID    string    `gorm:"column:service_id;type:varchar(120);primaryKey" json:"serviceID"`
	IP           string    `gorm:"column:ip;type:varchar(5000);default;null" json:"ip"`
	Port         string    `gorm:"column:port;type:varchar(5000);default;null" json:"port"`
	ServiceName  string    `gorm:"column:service_name;type:varchar(120);default;null" json:"serviceName"`
	ServiceType  string    `gorm:"column:service_type;type:varchar(120);default;null" json:"serviceType"`
	NodePort     string    `gorm:"column:node_port;type:varchar(5000);default;null" json:"nodePort"`
	TargetPort   string    `gorm:"column:target_port;type:varchar(5000);default;null" json:"targetPort"`
	WorkloadName string    `gorm:"column:workload_name;type:varchar(120);default:null" json:"workloadName"`
	WorkloadType string    `gorm:"column:workload_type;type:varchar(120);default:null" json:"workloadType"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime;null" json:"create_time"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName 指定数据库表名（默认是结构体名小写复数，这里显式指定与Django一致）
func (s *Service) TableName() string {
	return "service"
}

// String 实现fmt.Stringer接口，对应Django的__str__方法
func (s *Service) String() string {
	return s.ServiceName
}
