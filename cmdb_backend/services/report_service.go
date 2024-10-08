package services

import (
	"cmdb/models"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ReportService struct {
	DB *gorm.DB
}

func (s *ReportService) GenerateClusterGroupReport(c *gin.Context) {
	f := excelize.NewFile()

	var groups []models.ClusterGroup
	if err := s.DB.Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, group := range groups {
		sheetName := group.GroupName
		index, err := f.NewSheet(sheetName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		f.SetCellValue(sheetName, "A1", "Cluster Name")
		f.SetCellValue(sheetName, "B1", "Server Count")
		f.SetCellValue(sheetName, "C1", "Avg CPU Usage")
		f.SetCellValue(sheetName, "D1", "Avg Memory Usage")
		f.SetCellValue(sheetName, "E1", "Avg Disk Usage")
		f.SetCellValue(sheetName, "F1", "Estimated Disk Full Date")

		// Add data (this is a simplified version, you'd need to join tables and do more complex queries in reality)
		f.SetCellValue(sheetName, "A2", "Cluster1")
		f.SetCellValue(sheetName, "B2", 10)
		f.SetCellValue(sheetName, "C2", "50%")
		f.SetCellValue(sheetName, "D2", "60%")
		f.SetCellValue(sheetName, "E2", "70%")
		f.SetCellValue(sheetName, "F2", "2023-12-31")

		f.SetActiveSheet(index)
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=cluster_group_report.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")

	f.Write(c.Writer)
}

func (s *ReportService) GenerateIDCReport(c *gin.Context) {
	var serverResources []models.ServerResource
	if err := s.DB.Find(&serverResources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	idcStats := make(map[string]*models.IDCUsage)

	for _, resource := range serverResources {
		idcName := getIDCNameFromIP(resource.IP)
		if _, exists := idcStats[idcName]; !exists {
			idcStats[idcName] = &models.IDCUsage{
				IDCName: idcName,
			}
		}

		usage := idcStats[idcName]
		usage.TotalInstances++
		usage.AvgCPUUsage += resource.CPULoad
		usage.AvgMemoryUsage += (resource.UsedMemory / resource.TotalMemory) * 100
		usage.AvgDiskUsage += (resource.UsedDisk / resource.TotalDisk) * 100
	}

	// Calculate averages
	for _, usage := range idcStats {
		if usage.TotalInstances > 0 {
			usage.AvgCPUUsage /= float64(usage.TotalInstances)
			usage.AvgMemoryUsage /= float64(usage.TotalInstances)
			usage.AvgDiskUsage /= float64(usage.TotalInstances)
		}
	}

	f := excelize.NewFile()
	sheetName := "IDC Report"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers
	headers := []string{"IDC Name", "Total Instances", "Avg CPU Usage (%)", "Avg Memory Usage (%)", "Avg Disk Usage (%)"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Fill data
	row := 2
	var idcNames []string
	for idcName := range idcStats {
		idcNames = append(idcNames, idcName)
	}
	sort.Strings(idcNames)

	for _, idcName := range idcNames {
		usage := idcStats[idcName]
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), usage.IDCName)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), usage.TotalInstances)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("%.2f", usage.AvgCPUUsage))
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("%.2f", usage.AvgMemoryUsage))
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), fmt.Sprintf("%.2f", usage.AvgDiskUsage))
		row++
	}

	f.SetActiveSheet(index)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=idc_report.xlsx")
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func getIDCNameFromIP(ip string) string {
	parts := strings.Split(ip, ".")
	if len(parts) >= 2 {
		switch parts[1] {
		case "1":
			return "IDC-1"
		case "2":
			return "IDC-2"
		case "3":
			return "IDC-3"
		case "4":
			return "IDC-4"
		case "5":
			return "IDC-5"
		case "6":
			return "IDC-6"
		default:
			return "Unknown-IDC"
		}
	}
	return "Unknown-IDC"
}
