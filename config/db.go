package config

import (
	"log"
	"pdv/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("local.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Ativar integridade referencial no SQLite
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	_, err = sqlDB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Println("Erro ao ativar foreign_keys:", err)
	}

	// Migrar tabelas
	err = db.AutoMigrate(
		&model.Produto{},
		&model.Cliente{},
		&model.Venda{},
		&model.ItemVenda{},
	)
	if err != nil {
		log.Println("Erro ao migrar modelos:", err)
	}

	return db, nil
}

