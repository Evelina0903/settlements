package models

type District struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:text;not null"`
	Citys []City
}
