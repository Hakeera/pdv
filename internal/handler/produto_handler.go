package handler

import (
	"net/http"
	"pdv/internal/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ProdutosPage(c echo.Context) error {
	return c.Render(http.StatusOK, "produtos.html", map[string]any{
		"Titulo": "Controle de Produtos",
	}) 
}

// Retorna apenas a lista de produtos em HTML (fragment para HTMX)
func GetProdutos(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var produtos []model.Produto
		if err := db.Find(&produtos).Error; err != nil {
			// Retorna HTML de erro
			return c.HTML(http.StatusInternalServerError, `
				<div class="error">
					<p>Erro ao carregar produtos do banco de dados</p>
				</div>
			`)
		}
		
		// Renderiza apenas o fragment da lista de produtos
		return c.Render(http.StatusOK, "produtos-list.html", map[string]any{
			"Produtos": produtos,
		})
	}
}

func CreateProduto(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var produto model.Produto
		if err := c.Bind(&produto); err != nil {
			return c.Render(http.StatusBadRequest, "partials/erro", echo.Map{
				"Erro": "Dados inválidos",
			})
		}

		if produto.Nome == "" || produto.CodigoBarras == "" {
			return c.Render(http.StatusBadRequest, "partials/erro", echo.Map{
				"Erro": "Campos obrigatórios ausentes",
			})
		}

		if err := db.Create(&produto).Error; err != nil {
			return c.Render(http.StatusInternalServerError, "partials/erro", echo.Map{
				"Erro": "Erro ao salvar produto: " + err.Error(),
			})
		}

		return c.Render(http.StatusOK, "partials/produto_sucesso", echo.Map{
			"Produto": produto,
		})
	}
}

func GetProdutoByCodigo(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		codigo := c.Param("codigo")
		var p model.Produto
		if err := db.Where("codigo_barras = ?", codigo).First(&p).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.JSON(http.StatusNotFound, echo.Map{"erro": "produto não encontrado"})
			}
			return c.JSON(http.StatusInternalServerError, echo.Map{"erro": "erro ao buscar produto"})
		}
		return c.JSON(http.StatusOK, p)
	}
}

func DeleteProduto(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var produto model.Produto

		if err := db.First(&produto, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"erro": "Produto não encontrado"})
		}

		if err := db.Delete(&produto).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"erro": "Erro ao excluir produto"})
		}

		return c.JSON(http.StatusOK, echo.Map{"mensagem": "Produto excluído com sucesso"})
	}
}

