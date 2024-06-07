package controllers

import (
	"absensi/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetEmployees(c *gin.Context) {
	var employees []models.Employee
	if err := models.DB.Preload("Departement").Find(&employees).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    employees,
	})
}

func GetEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee
	if err := models.DB.Preload("Departement").Preload("Attendances").First(&employee, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah ada attendance yang dimuat
	if len(employee.Attendances) > 0 {
		// Mengambil attendance terbaru berdasarkan ClockIn
		latestAttendance := findLatestAttendance(employee.Attendances)
		if latestAttendance != nil {
			// Mengganti slice Attendances dengan attendance terbaru saja
			employee.Attendances = []models.Attendance{*latestAttendance}
			fmt.Print("employee.Attendances", employee.Attendances[0].ID)

			var attendanceHistories []models.AttendanceHistory
			if err := models.DB.Where("attendance_id = ?", employee.Attendances[0].ID).Find(&attendanceHistories).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Mengganti slice AttendanceHistories di Attendance terbaru dengan hasil query
			employee.Attendances[0].AttendanceHistories = attendanceHistories
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    employee,
	})
}

func CreateEmployee(c *gin.Context) {
	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.DB.Create(&employee).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Employee created",
		"data":    employee,
	})
}

func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee
	if err := models.DB.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.DB.Save(&employee).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Employee updated",
		"data":    employee,
	})
}

func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	// Hapus attendance_histories yang berelasi dengan employee
	if err := models.DB.Where("employee_id = ?", id).Delete(&models.AttendanceHistory{}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hapus attendances yang berelasi dengan employee
	if err := models.DB.Where("employee_id = ?", id).Delete(&models.Attendance{}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hapus employee
	if err := models.DB.Delete(&models.Employee{}, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Employee deleted",
		"data":    nil,
	})
}

// Fungsi untuk mencari attendance terbaru berdasarkan ClockIn
func findLatestAttendance(attendances []models.Attendance) *models.Attendance {
	var latestAttendance *models.Attendance
	for _, attendance := range attendances {
		if latestAttendance == nil || attendance.ClockIn.After(latestAttendance.ClockIn) {
			latestAttendance = &attendance
		}
	}
	return latestAttendance
}
