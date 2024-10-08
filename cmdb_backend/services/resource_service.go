package services

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"cmdb/models" // 导入 models 包

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResourceService struct {
	DB *gorm.DB
}

func (s *ResourceService) GetServerResources(c *gin.Context) {
	var resources []models.ServerResource

	// 获取查询参数
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	query := s.DB

	// 如果提供了时间范围，则添加时间过滤条件
	if startDate != "" && endDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
			return
		}
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
			return
		}
		query = query.Where("date_time BETWEEN ? AND ?", start, end)
	}

	if err := query.Find(&resources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取 ClusterGroup 信息
	var clusterGroups []models.ClusterGroup
	if err := s.DB.Find(&clusterGroups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建 ClusterGroup 映射
	clusterGroupMap := make(map[string]string)
	for _, cg := range clusterGroups {
		clusterGroupMap[cg.ClusterName] = cg.DepartmentName
	}

	// 为每个资源添加 DepartmentName
	type ResourceWithDepartment struct {
		models.ServerResource
		DepartmentName string `json:"department_name"`
	}

	var enrichedResources []ResourceWithDepartment
	for _, resource := range resources {
		enrichedResource := ResourceWithDepartment{
			ServerResource: resource,
			DepartmentName: clusterGroupMap[resource.ClusterName],
		}
		enrichedResources = append(enrichedResources, enrichedResource)
	}

	c.JSON(http.StatusOK, enrichedResources)
}

func (s *ResourceService) GetClusterResourceUsage(c *gin.Context) {
	var resources []models.ServerResource
	if err := s.DB.Find(&resources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resources)
}

func (s *ResourceService) GetResourceAlerts(c *gin.Context) {
	var alerts []models.ServerResource
	if err := s.DB.Where("used_memory / total_memory > 0.8 OR used_disk / total_disk > 0.8 OR cpu_load > 80").Find(&alerts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alerts)
}

func (s *ResourceService) PredictDiskFullDate(c *gin.Context) {
	var resources []models.ServerResource
	if err := s.DB.Order("date_time desc").Limit(30).Find(&resources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	predictions := make(map[uint]string)
	for _, r := range resources {
		daysLeft := float64(r.TotalDisk-r.UsedDisk) / (float64(r.UsedDisk) / 30)
		predictions[r.PoolID] = r.DateTime.AddDate(0, 0, int(daysLeft)).Format("2006-01-02")
	}

	c.JSON(http.StatusOK, predictions)
}

func (s *ResourceService) InsertServerResource(c *gin.Context) {
	var resource models.ServerResource
	if err := c.ShouldBindJSON(&resource); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.DB.Create(&resource).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// 新增的 GetClusterGroups 方法
func (s *ResourceService) GetClusterGroups(c *gin.Context) {
	var clusterGroups []models.ClusterGroup
	if err := s.DB.Find(&clusterGroups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cluster groups"})
		return
	}
	c.JSON(http.StatusOK, clusterGroups)
}

func (s *ResourceService) generateMockData() error {
	// 生成 20 个集群组
	clusterGroups := []models.ClusterGroup{}
	for i := 1; i <= 20; i++ {
		clusterGroups = append(clusterGroups, models.ClusterGroup{
			GroupName:      fmt.Sprintf("Group%d", (i-1)/3+1),
			ClusterName:    fmt.Sprintf("Cluster%d", i),
			DepartmentName: []string{"IT", "Finance", "HR", "Marketing", "Sales", "Operations"}[rand.Intn(6)],
		})
	}
	if err := s.DB.Create(&clusterGroups).Error; err != nil {
		return err
	}

	// 生成至少 20 台服务器
	for i := 1; i <= 60; i++ {
		host := models.HostPool{
			HostName:     fmt.Sprintf("host-%d", i),
			HostIP:       fmt.Sprintf("192.168.%d.%d", (i-1)/255+1, (i-1)%255+1),
			VCPUs:        uint(rand.Intn(32) + 8),
			RAM:          uint((rand.Intn(8) + 1) * 8),
			DiskSize:     uint((rand.Intn(10) + 1) * 100),
			HostType:     []string{"0", "1"}[rand.Intn(2)],
			SerialNumber: fmt.Sprintf("SN%06d", i),
			RackNumber:   fmt.Sprintf("R%02d", (i-1)/6+1),
			RackHeight:   uint(rand.Intn(4) + 1),
		}

		if err := s.DB.Create(&host).Error; err != nil {
			return err
		}

		// 为每个主机生成 1-3 个应用
		numApps := rand.Intn(3) + 1
		for j := 0; j < numApps; j++ {
			clusterGroup := clusterGroups[(i-1)/3]
			app := models.HostApplication{
				PoolID:         host.ID,
				ServerType:     []string{"MySQL", "Redis", "MongoDB", "Nginx", "Kafka", "Elasticsearch"}[rand.Intn(6)],
				ServerVersion:  fmt.Sprintf("%d.%d.%d", rand.Intn(5)+1, rand.Intn(10), rand.Intn(20)),
				ServerProtocol: []string{"TCP", "HTTP", "HTTPS"}[rand.Intn(3)],
				ServerAddr:     fmt.Sprintf("%s:%d", host.HostIP, 3000+rand.Intn(3000)),
				ClusterName:    clusterGroup.ClusterName,
			}
			if err := s.DB.Create(&app).Error; err != nil {
				return err
			}
		}

		// 生成 server_resources 数据
		serverResource := models.ServerResource{
			PoolID:      host.ID,
			TotalMemory: float64(host.RAM * 1024),
			UsedMemory:  float64(host.RAM*1024) * (0.3 + rand.Float64()*0.6),
			TotalDisk:   float64(host.DiskSize * 1024),
			UsedDisk:    float64(host.DiskSize*1024) * (0.2 + rand.Float64()*0.7),
			CPUCores:    int(host.VCPUs),
			CPULoad:     20.0 + rand.Float64()*70.0,
			DateTime:    time.Now().Add(-time.Duration(rand.Intn(24)) * time.Hour),
		}
		if err := s.DB.Create(&serverResource).Error; err != nil {
			return err
		}
	}

	return nil
}
