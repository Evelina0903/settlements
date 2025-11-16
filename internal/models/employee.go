package models

type Employee struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:text;not null"`
	BusinessTrips []BusinessTrip `gorm:"many2many:assignment_to_trips;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
