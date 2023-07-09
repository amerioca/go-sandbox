package models

type Permission struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `json:"name" gorm:"index:,unique"`
}
