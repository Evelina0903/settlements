package models

type City struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"type:text;not null"`
	TypeID     uint
	DistrictID uint
	Type       Type
	District   District
	Population int     `gorm:"type:int;not null"`
	Childrens  float32 `gorm:"type:float;not null"`
	Latitude   float32 `gorm:"type:float;not null"`
	Longitude  float32 `gorm:"type:float;not null"`
}
