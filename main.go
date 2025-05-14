package main

import (
	"log"
	"pdv/config"
	"pdv/internal/handler"

	"github.com/labstack/echo/v4"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("erro ao conectar ao banco:", err)
	}

	e := echo.New()

	// Inicializa rotas
	handler.RegisterRoutes(e, db)

	log.Println("Servidor rodando em http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}

