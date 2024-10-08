package controllers

import (
	"cmdb/models"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetHostsPoolDetail(c *gin.Context, db *gorm.DB) {
	var hosts []models.HostPool
	var applications []models.HostApplication

	db.Find(&hosts)
	db.Find(&applications)

	c.JSON(http.StatusOK, gin.H{
		"hosts":        hosts,
		"applications": applications,
	})
}

func CollectApplications(c *gin.Context, db *gorm.DB) {
	rand.Seed(time.Now().UnixNano())

	// Clear existing data
	db.Exec("DELETE FROM hosts_applications")
	db.Exec("DELETE FROM hosts_pool")

	// Simulate 30 hosts
	for i := 0; i < 30; i++ {
		host := models.HostPool{
			HostName:        "host" + string(rune(i+'0')),
			HostIP:          "192." + string(rune(i%5+1)) + "." + string(rune(i/256+'0')) + "." + string(rune(i%256+'0')),
			HostType:        "0",
			H3cID:           "h3c" + string(rune(i+'0')),
			H3cStatus:       "active",
			DiskSize:        uint(rand.Intn(1000)),
			RAM:             uint(rand.Intn(256)),
			VCPUs:           uint(rand.Intn(64)),
			IfH3cSync:       "yes",
			H3cImgID:        "img" + string(rune(i+'0')),
			H3cHmName:       "hm" + string(rune(i+'0')),
			IsDelete:        "no",
			LeafNumber:      "leaf" + string(rune(i+'0')),
			RackNumber:      "rack" + string(rune(i+'0')),
			RackHeight:      uint(rand.Intn(10)),
			RackStartNumber: uint(rand.Intn(100)),
			FromFactor:      uint(rand.Intn(5)),
			SerialNumber:    "serial" + string(rune(i+'0')),
			IsDeleted:       false,
			IsStatic:        false,
		}
		result := db.Create(&host)
		if result.Error != nil {
			log.Printf("Error creating host: %v", result.Error)
		}
	}

	// Simulate 50 applications
	for i := 0; i < 50; i++ {
		application := models.HostApplication{
			PoolID:         uint(rand.Intn(30) + 1),
			ServerType:     "type" + string(rune(i+'0')),
			ServerVersion:  "v1.0",
			ServerSubtitle: "subtitle" + string(rune(i+'0')),
			ClusterName:    "cluster" + string(rune(i+'0')),
			ServerProtocol: "http",
			ServerAddr:     "192.168." + string(rune(i/256+'0')) + "." + string(rune(i%256+'0')),
			DepartmentName: "department" + string(rune(i+'0')),
		}
		result := db.Create(&application)
		if result.Error != nil {
			log.Printf("Error creating application: %v", result.Error)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data collected successfully",
	})
}

func GetHostDetail(c *gin.Context, db *gorm.DB) {
	hostID := c.Param("id")
	var host models.HostPool
	var applications []models.HostApplication

	if err := db.First(&host, hostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Host not found"})
		return
	}

	if err := db.Where("pool_id = ?", hostID).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching applications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"host":         host,
		"applications": applications,
	})
}
