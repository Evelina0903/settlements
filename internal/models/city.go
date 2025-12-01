package models

type City struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"type:text;not null"`
	TypeID     uint
	DistrictID uint
	Type       Type
	District   District
	Population int     `gorm:"type:int;not null"`
	Childrens  int     `gorm:"type:int;not null"`
	Latitude   float64 `gorm:"not null"`
	Longitude  float64 `gorm:"not null"`
}
