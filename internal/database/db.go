package database

import (
	"Devenir_dev/internal/api/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initialise la connexion à la base de données en utilisant GORM
func InitDB() {
	var err error
	// Chaîne de connexion MySQL (modifie selon ton utilisateur et mot de passe)
	dsn := "root:ilyesgamer2005@@tcp(localhost:3306)/gestion_universite?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erreur de connexion à la base de données:", err)
	}
	err = DB.AutoMigrate(
		&models.User{},
		&models.Teacher{},
		&models.Module{},
		&models.Niveau{},
		&models.Voeux{},
		&models.Speciality{},
		&models.SavedOrganigram{},
		// ajoute d'autres structs ici
	)

	if err != nil {
		log.Fatal("Erreur de migration:", err)
	}
	// Vérifie que la connexion fonctionne
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Erreur lors de l'accès à la connexion SQL:", err)
	}

	// Teste la connexion
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Impossible de se connecter à la base de données:", err)
	}

	fmt.Println("Connexion à la base de données MySQL réussie.")
}

// GetDB retourne la connexion à la base de données
func GetDB() *gorm.DB {
	return DB
}

// CloseDB ferme la connexion à la base de données
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Erreur lors de la fermeture de la base de données:", err)
	}
	sqlDB.Close()
}
