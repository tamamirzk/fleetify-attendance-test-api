package router

import (
	"absensi/controllers"
	"absensi/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	// Set the global DB connection
	models.DB = db

	// Employee routes
	r.GET("/employees", controllers.GetEmployees)
	r.GET("/employees/:id", controllers.GetEmployee)
	r.POST("/employees", controllers.CreateEmployee)
	r.PUT("/employees/:id", controllers.UpdateEmployee)
	r.DELETE("/employees/:id", controllers.DeleteEmployee)

	// Department routes
	r.GET("/departments", controllers.GetDepartements)
	r.GET("/departments/:id", controllers.GetDepartement)
	r.POST("/departments", controllers.CreateDepartement)
	r.PUT("/departments/:id", controllers.UpdateDepartement)
	r.DELETE("/departments/:id", controllers.DeleteDepartement)

	// Attendance routes
	r.POST("/attendances", controllers.CreateAttendance)
	r.PUT("/attendances/:id", controllers.UpdateAttendance)
	r.GET("/attendance_logs", controllers.GetAttendanceLogs)
}
