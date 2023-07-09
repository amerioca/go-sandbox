package models

import "gorm.io/gorm"

type Product struct {
	ID          uint    `gorm:"primaryKey"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	gorm.Model
}

func (product *Product) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&Product{}).Count(&total)

	return total
}

func (product *Product) Take(db *gorm.DB, limit int, offset int) interface{} {
	var products []Product

	db.Offset(offset).Limit(limit).Find(&products)

	return products
}
