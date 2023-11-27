package database

import (
	"fmt"
	"log"
	"time"

	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

func seedProfileAdmin(db *gorm.DB) {
	query := `
        INSERT INTO profiles (id, created_at, updated_at, deleted_at, first_name, last_name, address, city, state, postal_code, phone_number, user_id)
        SELECT ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
        WHERE NOT EXISTS (
            SELECT 1 FROM profiles WHERE id = ?
        )
    `

	result := db.Exec(query,
		2,
		time.Now(),
		time.Now(),
		nil,
		"William",
		"Schwaab",
		"Rua Marechal Deodoro Da Fonseca 840",
		"Foz Do Iguacu",
		"Paraná",
		"22222-222",
		"11 1111 1111",
		5,
		2,
	)

	if result.Error != nil {
		fmt.Println("Error executing query:", result.Error)
		return
	}
}

func seedUserAdmin(db *gorm.DB) {
	query := `
     INSERT INTO users (id, created_at, updated_at, deleted_at,  email, password, is_active, is_staff, first_login, last_login)
        SELECT ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
        WHERE NOT EXISTS (
            SELECT 1 FROM users WHERE id = ?
        )
    `
	result := db.Exec(query,
		5,
		time.Now(),
		time.Now(),
		nil,
		"william@schwaab.me",
		"$2a$10$0.7NK60o45Qw1T4yw9oJSub9WUxsGkdyMf7HLWnGEsSrBkoJthY9W",
		1,
		1,
		time.Time{},
		time.Time{},
		5,
	)

	if result.Error != nil {
		fmt.Println("Error executing query:", result.Error)
		return
	}
}

func seedStoreConfig(db *gorm.DB) {
	var count int64

	if err := db.Model(&models.StoreConfig{}).Count(&count).Error; err != nil {
		log.Fatalf("Erro ao contar registros de StoreConfig: %v", err)
	}

	if count == 0 {
		storeConfig := models.StoreConfig{
			PixKeyType: models.PixTypeCNPJ,
		}

		if err := db.Create(&storeConfig).Error; err != nil {
			log.Fatalf("Falha ao criar StoreConfig: %v", err)
		}
	}
}

func seedDatabase(db *gorm.DB) {
	seedStoreConfig(db)
	seedUserAdmin(db)
	seedProfileAdmin(db)
}
