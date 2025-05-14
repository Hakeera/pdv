package model

import "gorm.io/gorm"

type Cliente struct {
	gorm.Model
	Nome       string
	Documento  string
	Endereco   string
	SyncStatus string `gorm:"default:pendente"`
}

