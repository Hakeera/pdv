package model

import "gorm.io/gorm"

type Produto struct {
	gorm.Model
	CodigoBarras string  `gorm:"uniqueIndex"`
	Nome         string
	Preco        float64
	Estoque      int
	SyncStatus   string `gorm:"default:pendente"`
}

