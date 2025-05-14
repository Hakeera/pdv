package controller

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"pdv/dto"
	"pdv/internal/model"
)

type VendaController struct {
	DB *gorm.DB
}

func NewVendaController(db *gorm.DB) *VendaController {
	return &VendaController{DB: db}
}

func (vc *VendaController) CriarVenda(c echo.Context) error {
	var req dto.CriarVendaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"erro": "JSON inválido"})
	}

	// Validar itens
	if len(req.Itens) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"erro": "A venda deve conter ao menos um item"})
	}

	// Iniciar transação
	err := vc.DB.Transaction(func(tx *gorm.DB) error {
		var total float64
		var itens []model.ItemVenda

		// Validar estoque e calcular total
		for _, item := range req.Itens {
			var produto model.Produto
			if err := tx.First(&produto, item.ProdutoID).Error; err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Produto não encontrado")
			}
			if produto.Estoque < item.Quantidade {
				return echo.NewHTTPError(http.StatusBadRequest, "Estoque insuficiente para o produto: "+produto.Nome)
			}

			// Calcular valor do item
			valorItem := float64(item.Quantidade)*item.PrecoUnitario - item.Desconto
			total += valorItem

			itens = append(itens, model.ItemVenda{
				ProdutoID:     item.ProdutoID,
				Quantidade:    item.Quantidade,
				PrecoUnitario: item.PrecoUnitario,
				Desconto:      item.Desconto,
				Subtotal:      valorItem,
			})
		}

		// Aplicar desconto total se houver
		totalComDesconto := total - req.DescontoTotal
		if totalComDesconto < 0 {
			totalComDesconto = 0
		}

		// Criar venda
		venda := model.Venda{
			ClienteID:     req.ClienteID,
			DescontoTotal: req.DescontoTotal,
			Total:         totalComDesconto,
		}
		if err := tx.Create(&venda).Error; err != nil {
			return err
		}

		// Associar itens à venda e salvar
		for i := range itens {
			itens[i].VendaID = venda.ID
			if err := tx.Create(&itens[i]).Error; err != nil {
				return err
			}

			// Atualizar estoque
			if err := tx.Model(&model.Produto{}).
				Where("id = ?", itens[i].ProdutoID).
				Update("estoque", gorm.Expr("estoque - ?", itens[i].Quantidade)).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Println("Erro ao criar venda:", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"erro": "Erro ao processar venda"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"mensagem": "Venda registrada com sucesso"})
}

