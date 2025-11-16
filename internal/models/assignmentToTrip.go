package models

import "gorm.io/gorm"

type AssignmentToTrip struct {
	gorm.Model
	ID           uint `gorm:"primaryKey"`
	MoneySpent   int
	Employee     Employee     `gorm:"foreignKey:EmployeeID;references:ID"`
	BusinessTrip BusinessTrip `gorm:"foreignKey:BusinessTripID;references:ID"`
}
