package model

import (
	"gorm.io/gorm"
)

type Venda struct {
	gorm.Model
	ClienteID *uint    `json:"cliente_id"` 
	Cliente   *Cliente `json:"cliente,omitempty"`
	Total     float64  `json:"total"`
	Itens     []ItemVenda
	DescontoTotal float64 `json:"descontotal"`
  
}

