package handler

import (
	"net/http"
	"pdv/internal/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetProdutos(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var produtos []model.Produto
		if err := db.Find(&produtos).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"erro": "erro ao buscar produtos"})
		}
		return c.JSON(http.StatusOK, produtos)
	}
}

func CreateProduto(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var p model.Produto
		if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"erro": "dados inválidos"})
		}

		if p.CodigoBarras == "" || p.Nome == "" || p.Preco <= 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{"erro": "campos obrigatórios ausentes"})
		}

		if err := db.Create(&p).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"erro": "erro ao criar produto"})
		}

		return c.JSON(http.StatusCreated, p)
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

