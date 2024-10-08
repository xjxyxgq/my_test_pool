package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"cmdb/models"
	"cmdb/services"

	"context"

	"github.com/chromedp/chromedp"
)

var db *gorm.DB

func main() {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 初始化数据库连接和生成模拟数据
	initDB()

	r := gin.Default()

	// 添加CORS中间件
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // 前端运行的地址
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	r.Use(cors.New(config))

	// 初始化服务
	resourceService := &services.ResourceService{DB: db}
	reportService := &services.ReportService{DB: db}
	emailService := &services.EmailService{
		SMTPHost:     "smtp.163.com",
		SMTPPort:     465,
		SMTPUser:     "xjxyxgq@163.com",
		SMTPPassword: "JSSHABDKJMXUHARS",
	}

	// 设置路由
	r.GET("/api/cmdb/v1/get_hosts_pool_detail", getHostsPoolDetail)
	r.GET("/api/cmdb/v1/get_host_detail/:id", getHostDetail)
	r.POST("/api/cmdb/v1/collect_applications", collectApplications)
	r.GET("/api/cmdb/v1/get_cluster_usage", getClusterUsage)

	// 添加新的接口
	r.GET("/api/cmdb/v1/cluster-resource-usage", resourceService.GetClusterResourceUsage)
	r.GET("/api/cmdb/v1/resource-alerts", resourceService.GetResourceAlerts)
	r.GET("/api/cmdb/v1/disk-full-prediction", resourceService.PredictDiskFullDate)
	r.GET("/api/cmdb/v1/cluster-group-report", reportService.GenerateClusterGroupReport)
	r.GET("/api/cmdb/v1/idc-report", reportService.GenerateIDCReport)
	r.GET("/api/cmdb/v1/server-resources", resourceService.GetServerResources)
	r.POST("/api/cmdb/v1/insert-server-resource", resourceService.InsertServerResource)
	r.POST("/api/cmdb/v1/send-email", emailService.SendEmail)
	r.GET("/api/cmdb/v1/cluster-groups", resourceService.GetClusterGroups)
	r.POST("/api/cmdb/v1/generate-and-send-report", generateAndSendReport)
	r.POST("/api/cmdb/v1/trigger-report", triggerReport)

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getHostsPoolDetail(c *gin.Context) {
	var hosts []models.HostPool
	if err := db.Preload("HostApplications").Find(&hosts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hosts"})
		return
	}

	// 获取所有的 ClusterGroup
	var clusterGroups []models.ClusterGroup
	if err := db.Find(&clusterGroups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cluster groups"})
		return
	}

	// 创建 ClusterGroup 映射
	clusterGroupMap := make(map[string]string)
	for _, cg := range clusterGroups {
		clusterGroupMap[cg.ClusterName] = cg.DepartmentName
	}

	// 为每个 HostApplication 添加 DepartmentName
	for i := range hosts {
		for j := range hosts[i].HostApplications {
			hosts[i].HostApplications[j].DepartmentName = clusterGroupMap[hosts[i].HostApplications[j].ClusterName]
		}
	}

	c.JSON(http.StatusOK, hosts)
}

func collectApplications(c *gin.Context) {
	// 这里我们将模拟数据生成
	generateMockData()
	c.JSON(http.StatusOK, gin.H{"message": "Mock data generated successfully"})
}

func initDB() {
	var err error
	dsn := "root:nov24feb11@tcp(127.0.0.1:3311)/cmdb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 删除现有的表
	db.Migrator().DropTable(&models.HostPool{}, &models.HostApplication{}, &models.ServerResource{}, &models.ClusterGroup{})

	// 重新创建表
	if err := db.AutoMigrate(&models.HostPool{}, &models.HostApplication{}, &models.ServerResource{}, &models.ClusterGroup{}); err != nil {
		log.Fatal("Failed to auto migrate:", err)
	}

	// 生成模拟数据
	generateMockData()
	log.Println("Mock data has been generated.")
}

func generateMockData() {
	// 1. 生成 cluster_group 数据
	clusterGroups := []models.ClusterGroup{}
	for i := 1; i <= 7; i++ {
		for j := 1; j <= 3; j++ {
			clusterGroups = append(clusterGroups, models.ClusterGroup{
				GroupName:      fmt.Sprintf("Group%d", i),
				ClusterName:    fmt.Sprintf("Cluster%d-%d", i, j),
				DepartmentName: []string{"IT", "Finance", "HR", "Marketing", "Sales", "Operations"}[rand.Intn(6)],
			})
		}
	}
	if err := db.Create(&clusterGroups).Error; err != nil {
		log.Fatal("Failed to create cluster groups:", err)
	}

	// 2. 生成 hosts_pool 数据
	hosts := []models.HostPool{}
	for i := 1; i <= 60; i++ {
		host := models.HostPool{
			HostName:     fmt.Sprintf("host-%d", i),
			HostIP:       fmt.Sprintf("192.%d.%d.%d", (i-1)%5+1, (i-1)/255+1, (i-1)%255+1),
			VCPUs:        uint(rand.Intn(32) + 8),
			RAM:          uint((rand.Intn(8) + 1) * 8),
			DiskSize:     uint((rand.Intn(10) + 1) * 100),
			HostType:     []string{"0", "1"}[rand.Intn(2)],
			SerialNumber: fmt.Sprintf("SN%06d", i),
			RackNumber:   fmt.Sprintf("R%02d", (i-1)/6+1),
			RackHeight:   uint(rand.Intn(4) + 1),
		}
		hosts = append(hosts, host)
	}
	if err := db.Create(&hosts).Error; err != nil {
		log.Fatal("Failed to create hosts:", err)
	}

	// 3. 生成 hosts_applications 数据
	for _, host := range hosts {
		numApps := rand.Intn(3) + 1
		for j := 0; j < numApps; j++ {
			clusterGroup := clusterGroups[rand.Intn(len(clusterGroups))]
			app := models.HostApplication{
				PoolID:         host.ID,
				ServerType:     []string{"MySQL", "Redis", "MongoDB", "Nginx", "Kafka", "Elasticsearch"}[rand.Intn(6)],
				ServerVersion:  fmt.Sprintf("%d.%d.%d", rand.Intn(5)+1, rand.Intn(10), rand.Intn(20)),
				ServerProtocol: []string{"TCP", "HTTP", "HTTPS"}[rand.Intn(3)],
				ServerAddr:     fmt.Sprintf("%s:%d", host.HostIP, 3000+rand.Intn(3000)),
				ClusterName:    clusterGroup.ClusterName,
				ServerPort:     3000 + rand.Intn(3000),
				ServerRole:     []string{"master", "slave"}[rand.Intn(2)],
				ServerStatus:   []string{"running", "stopped", "maintenance"}[rand.Intn(3)],
			}
			if err := db.Create(&app).Error; err != nil {
				log.Fatal("Failed to create application:", err)
			}
		}
	}

	// 4. 生成 server_resources 数据
	for _, host := range hosts {
		// 随机选择一个 ClusterGroup
		clusterGroup := clusterGroups[rand.Intn(len(clusterGroups))]

		serverResource := models.ServerResource{
			PoolID:       host.ID,
			ClusterName:  clusterGroup.ClusterName,
			GroupName:    clusterGroup.GroupName,
			IP:           host.HostIP,
			Port:         uint(3000 + rand.Intn(3000)),
			InstanceRole: []string{"master", "slave"}[rand.Intn(2)],
			TotalMemory:  float64(host.RAM * 1024),
			UsedMemory:   float64(host.RAM*1024) * (0.3 + rand.Float64()*0.6),
			TotalDisk:    float64(host.DiskSize * 1024),
			UsedDisk:     float64(host.DiskSize*1024) * (0.2 + rand.Float64()*0.7),
			CPUCores:     int(host.VCPUs),
			CPULoad:      20.0 + rand.Float64()*70.0,
			DateTime:     time.Now().Add(-time.Duration(rand.Intn(24)) * time.Hour),
		}
		if err := db.Create(&serverResource).Error; err != nil {
			log.Fatal("Failed to create server resource:", err)
		}
	}
}

func getHostDetail(c *gin.Context) {
	id := c.Param("id")
	var host models.HostPool
	if err := db.Preload("HostApplications").First(&host, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Host not found"})
		return
	}
	c.JSON(http.StatusOK, host)
}

func getClusterUsage(c *gin.Context) {
	var serverResources []models.ServerResource
	if err := db.Find(&serverResources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch server resources"})
		return
	}

	idcUsageMap := make(map[string]*models.IDCUsage)

	for _, resource := range serverResources {
		var host models.HostPool
		if err := db.First(&host, resource.PoolID).Error; err != nil {
			log.Printf("Failed to find host for PoolID %d: %v", resource.PoolID, err)
			continue
		}

		idcName := getIDCNameFromIP(host.HostIP)

		if _, exists := idcUsageMap[idcName]; !exists {
			idcUsageMap[idcName] = &models.IDCUsage{IDCName: idcName}
		}

		idcUsage := idcUsageMap[idcName]
		idcUsage.TotalInstances++
		idcUsage.AvgCPUUsage += resource.CPULoad
		idcUsage.AvgMemoryUsage += (resource.UsedMemory / resource.TotalMemory) * 100
		idcUsage.AvgDiskUsage += (resource.UsedDisk / resource.TotalDisk) * 100
	}

	var idcUsages []models.IDCUsage
	for _, usage := range idcUsageMap {
		if usage.TotalInstances > 0 {
			usage.AvgCPUUsage /= float64(usage.TotalInstances)
			usage.AvgMemoryUsage /= float64(usage.TotalInstances)
			usage.AvgDiskUsage /= float64(usage.TotalInstances)
		}
		idcUsages = append(idcUsages, *usage)
	}

	c.JSON(http.StatusOK, idcUsages)
}

func getIDCNameFromIP(ip string) string {
	// 这里实现根据 IP 地址判断 IDC 的逻辑
	// 例如，可以根据 IP 地址的前两位来判断
	parts := strings.Split(ip, ".")
	if len(parts) >= 2 {
		switch parts[1] {
		case "1":
			return "P1"
		case "2":
			return "P2"
		case "3":
			return "P3"
		case "4":
			return "P4"
		case "5":
			return "P5"
		case "6":
			return "P6"
		default:
			return "Unknown-IDC"
		}
	}
	return "Unknown-IDC"
}

func generateAndSendReport(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
		HTML  string `json:"html"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用 chromedp 生成截图
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate("data:text/html,"+request.HTML),
		chromedp.CaptureScreenshot(&buf),
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate screenshot"})
		return
	}

	// 创建邮件内容
	emailContent := fmt.Sprintf(`
        <html>
            <body>
                <h1>服务器资源使用情况报告</h1>
                <img src="cid:screenshot" alt="Server Resources" style="max-width: 100%%;" />
            </body>
        </html>
    `)

	// 发送邮件
	err := sendEmail(request.Email, "服务器资源使用情况报告", emailContent, buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report generated and sent successfully"})
}

func sendEmail(to, subject, body string, attachment []byte) error {
	// 这里使用您的邮件服务配置
	emailService := &services.EmailService{
		SMTPHost:     "smtp.163.com",
		SMTPPort:     465,
		SMTPUser:     "xjxyxgq@163.com",
		SMTPPassword: "JSSHABDKJMXUHARS",
	}

	return emailService.SendEmailWithAttachment(to, subject, body, attachment)
}

func triggerReport(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	htmlContent := generateHTMLContent()

	err := generateAndSendReportLogic(request.Email, htmlContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate and send report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report triggered successfully"})
}

func generateHTMLContent() string {
	var serverResources []models.ServerResource
	db.Find(&serverResources)

	var clusterGroups []models.ClusterGroup
	db.Find(&clusterGroups)

	var html strings.Builder
	html.WriteString(`
    <html>
    <head>
        <style>
            body { font-family: Arial, sans-serif; }
            table { border-collapse: collapse; width: 100%; }
            th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
            th { background-color: #f2f2f2; }
            .alert { padding: 10px; margin-bottom: 10px; border-radius: 4px; }
            .alert-warning { background-color: #fff3cd; color: #856404; }
            .alert-danger { background-color: #f8d7da; color: #721c24; }
        </style>
    </head>
    <body>
        <h1>服务器资源使用情况报告</h1>
    `)

	// 资源警报
	html.WriteString("<h2>资源警报</h2>")
	for _, resource := range serverResources {
		memoryUsage := (resource.UsedMemory / resource.TotalMemory) * 100
		diskUsage := (resource.UsedDisk / resource.TotalDisk) * 100
		cpuUsage := resource.CPULoad

		if memoryUsage > 80 || diskUsage > 80 || cpuUsage > 80 {
			html.WriteString(fmt.Sprintf(`
                <div class="alert alert-danger">
                    %s (%s %s) | 内存: %.2f%% | 磁盘: %.2f%% | CPU: %.2f%%
                </div>
            `, resource.IP, resource.GroupName, resource.ClusterName, memoryUsage, diskUsage, cpuUsage))
		} else if memoryUsage < 10 || diskUsage < 10 || cpuUsage < 10 {
			html.WriteString(fmt.Sprintf(`
                <div class="alert alert-warning">
                    %s (%s %s) | 内存: %.2f%% | 磁盘: %.2f%% | CPU: %.2f%%
                </div>
            `, resource.IP, resource.GroupName, resource.ClusterName, memoryUsage, diskUsage, cpuUsage))
		}
	}

	// 集群资源使用情况
	html.WriteString("<h2>集群资源使用情况</h2>")
	for _, group := range clusterGroups {
		html.WriteString(fmt.Sprintf("<h3>%s-%s</h3>", group.GroupName, group.ClusterName))
		html.WriteString("<table>")
		html.WriteString("<tr><th>资源</th><th>使用率</th></tr>")

		var memoryUsage, diskUsage, cpuUsage float64
		var count int
		for _, resource := range serverResources {
			if resource.GroupName == group.GroupName {
				memoryUsage += (resource.UsedMemory / resource.TotalMemory) * 100
				diskUsage += (resource.UsedDisk / resource.TotalDisk) * 100
				cpuUsage += resource.CPULoad
				count++
			}
		}
		if count > 0 {
			html.WriteString(fmt.Sprintf("<tr><td>内存</td><td>%.2f%%</td></tr>", memoryUsage/float64(count)))
			html.WriteString(fmt.Sprintf("<tr><td>磁盘</td><td>%.2f%%</td></tr>", diskUsage/float64(count)))
			html.WriteString(fmt.Sprintf("<tr><td>CPU</td><td>%.2f%%</td></tr>", cpuUsage/float64(count)))
		}
		html.WriteString("</table>")
	}

	// 服务器资源详情
	html.WriteString("<h2>服务器资源详情</h2>")
	html.WriteString(`
        <table>
            <tr>
                <th>Instance ID</th>
                <th>IP</th>
                <th>Cluster Name</th>
                <th>Group Name</th>
                <th>CPU Usage</th>
                <th>Memory Usage</th>
                <th>Disk Usage</th>
            </tr>
    `)
	for _, resource := range serverResources {
		html.WriteString(fmt.Sprintf(`
            <tr>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%.2f%%</td>
                <td>%.2f%%</td>
                <td>%.2f%%</td>
            </tr>
        `,
			resource.IP,
			resource.ClusterName,
			resource.GroupName,
			resource.CPULoad,
			(resource.UsedMemory/resource.TotalMemory)*100,
			(resource.UsedDisk/resource.TotalDisk)*100))
	}
	html.WriteString("</table>")

	html.WriteString("</body></html>")
	return html.String()
}

func generateAndSendReportLogic(email, htmlContent string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate("data:text/html,"+htmlContent),
		chromedp.CaptureScreenshot(&buf),
	); err != nil {
		return err
	}

	emailBody := fmt.Sprintf(`
        <html>
            <body>
                <h1>服务器资源使用情况报告</h1>
                <img src="cid:screenshot" alt="Server Resources" style="max-width: 100%%;" />
            </body>
        </html>
    `)

	return sendEmail(email, "服务器资源使用情况报告", emailBody, buf)
}
