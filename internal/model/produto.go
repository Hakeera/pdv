package model

import "gorm.io/gorm"

type Produto struct {
	gorm.Model
	CodigoBarras string `gorm:"uniqueIndex" json:"codigo_barras" form:"codigo_barras"`
	Nome         string `json:"nome" form:"nome"`
	Preco        int64  `json:"preco" form:"preco"`        
	Estoque      int    `json:"estoque" form:"estoque"`
	SyncStatus   string `gorm:"default:pendente" json:"sync_status" form:"sync_status"`
}
