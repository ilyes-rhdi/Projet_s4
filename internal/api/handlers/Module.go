package handlers

import (
	"github.com/ilyes-rhdi/Projet_s4/internal/database"
	"github.com/ilyes-rhdi/Projet_s4/internal/api/models"
	"encoding/json"
	"net/http"
)

func GetAllModules(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	var modules []models.Module

	if err := db.Preload("Niveaux").Find(&modules).Error; err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(modules); err != nil {
		http.Error(w, "Erreur lors de l'envoi des données", http.StatusInternalServerError)
	}
}