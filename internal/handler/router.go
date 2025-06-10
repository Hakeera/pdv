package handler

import (

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	// Rota raiz
	e.GET("/", Home)

	// Produtos
	e.GET("/produtos/page", ProdutosPage) 
	e.GET("/produtos", GetProdutos(db))
	e.POST("/produtos", CreateProduto(db))
	e.GET("/produtos/:codigo", GetProdutoByCodigo(db))
	e.PUT("/produtos/:id", UpdateProduto(db))
	e.DELETE("/produtos/:id", DeleteProduto(db))

	// Vendas
	e.GET("/vendas", GetVendas(db))
	e.GET("/vendas/:id", GetVendaByID(db))
	e.POST("/vendas", CreateVenda(db)) 
	e.DELETE("/vendas/:id", DeleteVenda(db))

}

