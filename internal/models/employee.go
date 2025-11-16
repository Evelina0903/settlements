package models

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:text;not null"`
	BusinessTrips []BusinessTrip `gorm:"many2many:assignment_to_trip;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
