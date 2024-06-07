package main

import (
	"absensi/models"
	"absensi/router"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dbUser := "root"
	dbPassword := ""
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbName := "golang-api"

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Assign db to the models.DB variable for global access
	models.DB = db

	// Auto-migrate the schema in the correct order
	err = db.AutoMigrate(
		&models.Departement{},
		&models.Employee{},
		&models.Attendance{},
		&models.AttendanceHistory{},
	)
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	r := gin.Default()

	// Middleware CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Handle OPTIONS method
	r.Use(func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	router.SetupRouter(r, db)
	r.Run(":8080")
}
