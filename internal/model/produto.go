package model

import "gorm.io/gorm"

type Produto struct {
	gorm.Model
	CodigoBarras string `gorm:"uniqueIndex" json:"codigo_barras"`
	Nome         string `json:"nome"`
	Preco        int64  `json:"preco"`        
	Estoque      int    `json:"estoque"`
	SyncStatus   string `gorm:"default:pendente" json:"sync_status"`
}

