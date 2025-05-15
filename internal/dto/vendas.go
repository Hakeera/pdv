package dto

type ItemVendaRequest struct {
	ProdutoID     uint    `json:"produto_id"`
	Quantidade    int     `json:"quantidade"`
	PrecoUnitario int64   `json:"preco_unitario"`
	Desconto      int64   `json:"desconto"`
}

type CriarVendaRequest struct {
	ClienteID     *uint              `json:"cliente_id"`     // nil = venda sem identificação
	DescontoTotal int64   `json:"desconto_total"` // opcional
	Itens         []ItemVendaRequest `json:"itens"`
}

