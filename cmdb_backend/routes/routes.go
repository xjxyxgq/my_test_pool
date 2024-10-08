package routes

import (
	"cmdb/controllers"
	"cmdb/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/cmdb/v1/get_hosts_pool_detail", func(c *gin.Context) {
		controllers.GetHostsPoolDetail(c, db)
	})

	r.GET("/api/cmdb/v1/get_host_detail/:id", func(c *gin.Context) {
		controllers.GetHostDetail(c, db)
	})
}

// 更新任何使用了 cluster_id 或 group_name 的路由处理函数
// 例如：
func getServerResources(c *gin.Context, resourceService *services.ResourceService) {
	resourceService.GetServerResources(c)
}
