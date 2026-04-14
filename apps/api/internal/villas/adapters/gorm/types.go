package gorm

type Villa struct {
	ID     string `gorm:"primaryKey"`
	Numero int    `gorm:"unique"`
	Estado string `gorm:"column:estado"`
}
