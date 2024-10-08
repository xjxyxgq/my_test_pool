package models

import (
	"time"

	"gorm.io/gorm"
)

type HostPool struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	HostName         string            `gorm:"size:50;not null" json:"host_name"`
	HostIP           string            `gorm:"size:50;not null" json:"host_ip"`
	HostType         string            `gorm:"size:10" json:"host_type"`
	H3cID            string            `gorm:"size:50" json:"h3c_id"`
	H3cStatus        string            `gorm:"size:20" json:"h3c_status"`
	DiskSize         uint              `gorm:"default:null" json:"disk_size"`
	RAM              uint              `gorm:"default:null" json:"ram"`
	VCPUs            uint              `gorm:"default:null" json:"vcpus"`
	IfH3cSync        string            `gorm:"size:10" json:"if_h3c_sync"`
	H3cImgID         string            `gorm:"size:50" json:"h3c_img_id"`
	H3cHmName        string            `gorm:"size:1000" json:"h3c_hm_name"`
	IsDelete         string            `gorm:"size:10" json:"is_delete"`
	LeafNumber       string            `gorm:"size:50" json:"leaf_number"`
	RackNumber       string            `gorm:"size:10" json:"rack_number"`
	RackHeight       uint              `gorm:"default:null" json:"rack_height"`
	RackStartNumber  uint              `gorm:"default:null" json:"rack_start_number"`
	FromFactor       uint              `gorm:"default:null" json:"from_factor"`
	SerialNumber     string            `gorm:"size:50" json:"serial_number"`
	IsDeleted        bool              `gorm:"not null;default:false" json:"is_deleted"`
	IsStatic         bool              `gorm:"not null;default:false" json:"is_static"`
	CreateTime       time.Time         `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime       time.Time         `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"update_time"`
	HostApplications []HostApplication `gorm:"foreignKey:PoolID" json:"host_applications"`
}

func (HostPool) TableName() string {
	return "hosts_pool"
}

type HostApplication struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	PoolID         uint           `gorm:"not null" json:"pool_id"`
	HostPool       HostPool       `gorm:"foreignKey:PoolID" json:"-"`
	ServerType     string         `gorm:"size:30" json:"server_type"`
	ServerVersion  string         `gorm:"size:30" json:"server_version"`
	ServerSubtitle string         `gorm:"size:30" json:"server_subtitle"`
	ClusterName    string         `gorm:"size:64" json:"cluster_name"`
	ServerProtocol string         `gorm:"size:64" json:"server_protocol"`
	ServerAddr     string         `gorm:"size:100" json:"server_addr"`
	ServerPort     int            `gorm:"not null" json:"server_port"`
	ServerRole     string         `gorm:"size:100" json:"server_role"`
	ServerStatus   string         `gorm:"size:100" json:"server_status"`
	DepartmentName string         `gorm:"size:100" json:"department_name"`
	CreateTime     time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime     time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"update_time"`
}

func (HostApplication) TableName() string {
	return "hosts_applications"
}

// ServerResource 表示服务器资源的结构
type ServerResource struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	PoolID       uint           `gorm:"not null" json:"pool_id"`
	ClusterName  string         `json:"cluster_name"`
	GroupName    string         `json:"group_name"` // 新增字段
	IP           string         `json:"ip"`
	Port         uint           `json:"port"`
	InstanceRole string         `json:"instance_role"`
	TotalMemory  float64        `json:"total_memory"`
	UsedMemory   float64        `json:"used_memory"`
	TotalDisk    float64        `json:"total_disk"`
	UsedDisk     float64        `json:"used_disk"`
	CPUCores     int            `json:"cpu_cores"`
	CPULoad      float64        `json:"cpu_load"`
	DateTime     time.Time      `json:"date_time"`
}

// ClusterGroup 表示集群组的结构
type ClusterGroup struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	GroupName      string         `json:"group_name" gorm:"not null"`
	ClusterName    string         `json:"cluster_name" gorm:"not null"`
	DepartmentName string         `json:"department_name" gorm:"not null"` // 新增字段
}

// IDCUsage 表示 IDC 使用情况的结构
type IDCUsage struct {
	IDCName        string  `json:"idc_name"`
	TotalInstances int     `json:"total_instances"`
	AvgCPUUsage    float64 `json:"avg_cpu_usage"`
	AvgMemoryUsage float64 `json:"avg_memory_usage"`
	AvgDiskUsage   float64 `json:"avg_disk_usage"`
}
