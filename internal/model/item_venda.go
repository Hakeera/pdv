package model

import "gorm.io/gorm"

type ItemVenda struct {
	gorm.Model
	VendaID       uint    `json:"venda_id"`
	ProdutoID     uint    `json:"produto_id"`
	Produto       Produto `gorm:"foreignKey:ProdutoID"`
	Quantidade    int     `json:"quantidade"`
	PrecoUnitario int64   `json:"preco_unitario"` 
	Desconto      int64   `json:"desconto"`      
	Subtotal      int64   `json:"subtotal"`     
}

