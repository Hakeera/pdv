package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"pdv/config"
	"pdv/internal/handler"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("Erro ao renderizar template %s: %v", name, err)
	}
	return err
}

func main() { db, err := config.InitDB()
	if err != nil {
		log.Fatal("erro ao conectar ao banco:", err)
	}

	e := echo.New()
	
	t := &Template{
	    templates: template.Must(template.ParseGlob("internal/view/**/*.html")),
	}
	e.Renderer = t

	// Inicializa rotas
	handler.RegisterRoutes(e, db)

	log.Println("Servidor rodando em http://localhost:1323")
	for _, tmpl := range t.templates.Templates() {
		fmt.Println("Template carregado:", tmpl.Name())
	}

	e.Logger.Fatal(e.Start(":1323"))
}

