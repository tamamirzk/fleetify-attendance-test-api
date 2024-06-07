package controllers

import (
	"absensi/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateAttendanceReq struct {
	EmployeeID  uint   `json:"EmployeeID"`
	Description string `json:"description"`
}
type UpdateAttendanceReq struct {
	Description string `json:"description"`
}

func CreateAttendance(c *gin.Context) {
	var req CreateAttendanceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var attendance models.Attendance
	attendance.EmployeeID = req.EmployeeID
	attendance.ClockIn = time.Now()
	if err := models.DB.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create an entry in AttendanceHistory
	attendanceHistory := models.AttendanceHistory{
		EmployeeID:     attendance.EmployeeID,
		AttendanceID:   attendance.ID,
		DateAttendance: attendance.ClockIn,
		AttendanceType: 1,               // 1 means "in" or "masuk"
		Description:    req.Description, // Get description from request body
	}

	if err := models.DB.Create(&attendanceHistory).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attendance created successfully",
		"data": gin.H{
			"attendance":        attendance,
			"attendanceHistory": attendanceHistory,
		},
	})

}

func UpdateAttendance(c *gin.Context) {
	var req UpdateAttendanceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	var attendance models.Attendance
	if err := models.DB.First(&attendance, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	attendance.ClockOut = &now
	if err := models.DB.Save(&attendance).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create an entry in AttendanceHistory
	attendanceHistory := models.AttendanceHistory{
		EmployeeID:     attendance.EmployeeID,
		AttendanceID:   attendance.ID,
		DateAttendance: *attendance.ClockOut,
		AttendanceType: 2,               // 2 means "out" or "keluar"
		Description:    req.Description, // Get description from request body
	}

	if err := models.DB.Create(&attendanceHistory).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attendance updated successfully",
		"data": gin.H{
			"attendance":        attendance,
			"attendanceHistory": attendanceHistory,
		},
	})
}

func GetAttendanceLogs(c *gin.Context) {
	var attendances []models.Attendance
	query := models.DB.Preload("AttendanceHistories").Preload("Employee").Preload("Employee.Departement")

	// Filter by date
	if date := c.Query("date"); date != "" {
		query = query.Where("DATE(clock_in) = ?", date)
	}

	// Filter by department
	if departmentID := c.Query("department_id"); departmentID != "" {
		query = query.Joins("JOIN employees ON employees.id = attendances.employee_id").Where("employees.departement_id = ?", departmentID)
	}

	if err := query.Find(&attendances).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate punctuality
	type AttendanceLog struct {
		Attendance  models.Attendance
		IsOnTimeIn  bool
		IsOnTimeOut bool
	}

	var logs []AttendanceLog
	for _, attendance := range attendances {
		// Parse the MaxClockInTime and MaxClockOutTime from string to time.Time with default date
		maxClockInTime, err := time.Parse("2006-01-02 15:04", "1970-01-01 "+attendance.Employee.Departement.MaxClockInTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid MaxClockInTime format"})
			return
		}

		maxClockOutTime, err := time.Parse("2006-01-02 15:04", "1970-01-01 "+attendance.Employee.Departement.MaxClockOutTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid MaxClockOutTime format"})
			return
		}

		// Extract only time part from attendance.ClockIn
		clockInTime := time.Date(1970, 1, 1, attendance.ClockIn.Hour(), attendance.ClockIn.Minute(), attendance.ClockIn.Second(), 0, time.UTC)
		var clockOutTime time.Time
		if attendance.ClockOut != nil {
			clockOutTime = time.Date(1970, 1, 1, attendance.ClockOut.Hour(), attendance.ClockOut.Minute(), attendance.ClockOut.Second(), 0, time.UTC)
		}

		isOnTimeIn := clockInTime.Before(maxClockInTime) || clockInTime.Equal(maxClockInTime)
		isOnTimeOut := attendance.ClockOut == nil || clockOutTime.After(maxClockOutTime) || clockOutTime.Equal(maxClockOutTime)

		logs = append(logs, AttendanceLog{
			Attendance:  attendance,
			IsOnTimeIn:  isOnTimeIn,
			IsOnTimeOut: isOnTimeOut,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    logs,
	})
}
