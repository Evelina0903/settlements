package models

import "gorm.io/gorm"

type BusinessTrip struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Destination string `gorm:"type:text;not null"`
	StartAt   string `gorm:"type:date;not null"`
	EndAt     string `gorm:"type:date;not null"`
	Employees []Employee `gorm:"many2many:assignment_to_trip;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
