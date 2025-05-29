package handler

import (
	"net/http"
	"pdv/internal/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

	type ItemVendaInput struct {
	ProdutoID  uint `json:"produto_id"`
	Quantidade int  `json:"quantidade"`
}

type VendaInput struct {
	ClienteID *uint            `json:"cliente_id"`
	Itens     []ItemVendaInput `json:"itens"`
}

func CreateVenda(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input VendaInput
		if err := c.Bind(&input); err != nil || len(input.Itens) == 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{"erro": "dados inválidos ou itens ausentes"})
		}
		
		var total int64 
		var itens []model.ItemVenda

		// Transação
		err := db.Transaction(func(tx *gorm.DB) error {
			for _, item := range input.Itens {
				var prod model.Produto
				if err := tx.First(&prod, item.ProdutoID).Error; err != nil {
					return err
				}

				if prod.Estoque < item.Quantidade {
					return echo.NewHTTPError(http.StatusBadRequest, "Estoque insuficiente para o produto "+prod.Nome)
				}

				// Atualiza estoque
				prod.Estoque -= item.Quantidade
				if err := tx.Save(&prod).Error; err != nil {
					return err
				}

				subtotal := (int64(item.Quantidade) * prod.Preco)
				total += subtotal

				itens = append(itens, model.ItemVenda{
					ProdutoID:  prod.ID,
					Quantidade: item.Quantidade,
					PrecoUnitario:  prod.Preco,
				})
			}

			venda := model.Venda{
				Total:     total,
				Itens:     itens,
			}

			return tx.Create(&venda).Error
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"erro": err.Error()})
		}

		return c.JSON(http.StatusCreated, echo.Map{"mensagem": "venda criada com sucesso", "total": total})
	}
}

func GetVendas(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var vendas []model.Venda
		if err := db.Preload("Itens.Produto").Find(&vendas).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"erro": "erro ao buscar vendas"})
		}
		return c.JSON(http.StatusOK, vendas)
	}
}

func DeleteVenda(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var venda model.Venda

		if err := db.Preload("Itens").First(&venda, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"erro": "Venda não encontrada"})
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			// Restaurar estoque
			for _, item := range venda.Itens {
				if err := tx.Model(&model.Produto{}).
					Where("id = ?", item.ProdutoID).
					Update("estoque", gorm.Expr("estoque + ?", item.Quantidade)).Error; err != nil {
					return err
				}
			}

			// Deletar venda (soft delete)
			if err := tx.Delete(&venda).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"erro": "Erro ao excluir venda"})
		}

		return c.JSON(http.StatusOK, echo.Map{"mensagem": "Venda excluída com sucesso"})
	}
}

