package controllers

import (
	"absensi/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDepartements(c *gin.Context) {
	var departements []models.Departement
	if err := models.DB.Preload("Employees").Find(&departements).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    departements,
	})
}

func GetDepartement(c *gin.Context) {
	id := c.Param("id")
	var departement models.Departement
	if err := models.DB.Preload("Employees").First(&departement, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    departement,
	})
}

func CreateDepartement(c *gin.Context) {
	var departement models.Departement
	if err := c.ShouldBindJSON(&departement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.DB.Create(&departement).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Departement created",
		"data":    departement,
	})
}

func UpdateDepartement(c *gin.Context) {
	id := c.Param("id")
	var departement models.Departement
	if err := models.DB.First(&departement, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&departement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.DB.Save(&departement).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Departement updated",
		"data":    departement,
	})
}

func DeleteDepartement(c *gin.Context) {
	id := c.Param("id")
	if err := models.DB.Delete(&models.Departement{}, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Departement deleted",
		"data":    nil,
	})
}
