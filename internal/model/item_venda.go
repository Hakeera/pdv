package model

import "gorm.io/gorm"

type ItemVenda struct {
	gorm.Model
	VendaID       uint    `json:"venda_id"`
	ProdutoID     uint    `json:"produto_id"`
	Produto       Produto `gorm:"foreignKey:ProdutoID"`
	Quantidade    int     `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"` 
	Desconto      float64 `json:"desconto"`      
	Subtotal      float64 `json:"subtotal"`     
}

