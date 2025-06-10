package handler

import (
	"fmt"
	"net/http"
	"pdv/internal/model"
	"strings"

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
		
		// Bind dos dados do formulário
		if err := c.Bind(&produto); err != nil {
			return c.HTML(http.StatusBadRequest, `
				<div class="error-message">
					<h3>❌ Erro ao cadastrar produto</h3>
					<p>Dados inválidos: `+err.Error()+`</p>
				</div>
			`)
		}

		// Validação dos campos obrigatórios
		if produto.Nome == "" || produto.CodigoBarras == "" {
			return c.HTML(http.StatusBadRequest, `
				<div class="error-message">
					<h3>❌ Erro ao cadastrar produto</h3>
					<p>Nome e Código de Barras são obrigatórios</p>
				</div>
			`)
		}

		// Validações adicionais
		if produto.Preco < 0 {
			return c.HTML(http.StatusBadRequest, `
				<div class="error-message">
					<h3>❌ Erro ao cadastrar produto</h3>
					<p>Preço não pode ser negativo</p>
				</div>
			`)
		}

		if produto.Estoque < 0 {
			return c.HTML(http.StatusBadRequest, `
				<div class="error-message">
					<h3>❌ Erro ao cadastrar produto</h3>
					<p>Estoque não pode ser negativo</p>
				</div>
			`)
		}

		// Criar produto no banco de dados
		if err := db.Create(&produto).Error; err != nil {
			// Verificar se é erro de duplicação (código de barras único)
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || 
			   strings.Contains(err.Error(), "duplicate key") {
				return c.HTML(http.StatusBadRequest, `
					<div class="error-message">
						<h3>❌ Erro ao cadastrar produto</h3>
						<p>Produto com este código de barras já existe</p>
					</div>
				`)
			}
			
			return c.HTML(http.StatusInternalServerError, `
				<div class="error-message">
					<h3>❌ Erro ao cadastrar produto</h3>
					<p>Erro ao salvar produto: `+err.Error()+`</p>
				</div>
			`)
		}

		// Formatear preço para exibição (centavos para reais)
		precoFormatado := fmt.Sprintf("%.2f", float64(produto.Preco)/100.0)

		// Sucesso - retornar HTML de sucesso
		successHTML := fmt.Sprintf(`
			<div class="success-message">
				<h3>✅ Produto cadastrado com sucesso!</h3>
				<div class="produto-info">
					<h3>%s</h3>
					<p><strong>Código:</strong> %s</p>
					<p><strong>Preço:</strong> <span class="price">R$ %s</span></p>
					<p><strong>Estoque:</strong> %d unidades</p>
					<p><strong>ID:</strong> %d</p>
				</div>
				<p><small>Modal será fechado automaticamente...</small></p>
			</div>
		`, produto.Nome, produto.CodigoBarras, precoFormatado, produto.Estoque, produto.ID)

		return c.HTML(http.StatusOK, successHTML)
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

