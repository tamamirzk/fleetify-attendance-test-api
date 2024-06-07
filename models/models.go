package models

import (
	"time"

	"gorm.io/gorm"
)

var DB *gorm.DB

type Departement struct {
	ID              uint       `gorm:"primaryKey"`
	DepartementName string     `gorm:"type:varchar(255)"`
	MaxClockInTime  string     `gorm:"type:varchar(8)"` // Format HH:MM:SS
	MaxClockOutTime string     `gorm:"type:varchar(8)"` // Format HH:MM:SS
	Employees       []Employee `gorm:"foreignKey:DepartementID"`
}

type Employee struct {
	ID uint `gorm:"primaryKey"`
	// EmployeeID    string       `gorm:"unique;type:varchar(50)"`
	DepartementID uint         `gorm:"index"`
	Departement   Departement  `gorm:"foreignKey:DepartementID"`
	Name          string       `gorm:"type:varchar(255)"`
	Address       string       `gorm:"type:text"`
	Attendances   []Attendance `gorm:"foreignKey:EmployeeID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Attendance struct {
	ID                  uint                `gorm:"primaryKey"`
	EmployeeID          uint                `gorm:"index"`
	Employee            Employee            `gorm:"foreignKey:EmployeeID"`
	ClockIn             time.Time           `gorm:"type:datetime"`
	ClockOut            *time.Time          `gorm:"type:datetime"` // Pointer to allow null
	AttendanceHistories []AttendanceHistory `gorm:"foreignKey:AttendanceID"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type AttendanceHistory struct {
	ID             uint       `gorm:"primaryKey"`
	EmployeeID     uint       `gorm:"index"`
	Employee       Employee   `gorm:"foreignKey:EmployeeID"`
	AttendanceID   uint       `gorm:"index"`
	Attendance     Attendance `gorm:"foreignKey:AttendanceID"`
	DateAttendance time.Time  `gorm:"type:datetime"`
	AttendanceType int8       `gorm:"type:tinyint"`
	Description    string     `gorm:"type:text"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
