package models

type Role struct {
	ID          uint16       `gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"index:,unique"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions"`
}
