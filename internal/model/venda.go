package model

import (
	"gorm.io/gorm"
)

type Venda struct {
	gorm.Model
	Total     int64	    `json:"total"`
	Itens     []ItemVenda
	DescontoTotal int64 `json:"descontototal"`
  
}

