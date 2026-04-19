package gorm

type Unidad struct {
	ID     string `gorm:"primaryKey"`
	Codigo string `gorm:"unique"`
	Estado string `gorm:"column:estado"`
}
