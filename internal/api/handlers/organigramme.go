package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Structures des entrées

type Module struct {
	Nom      string `json:"name"`
	Priorite int    `json:"priority"`
	Heures   int    `json:"hours"`
	Type     string `json:"type"` // cours, td, tp
}

type FicheVoeux struct {
	Prof    string   `json:"prof_name"`
	Modules []Module `json:"modules"`
}

type CaseEmploi struct {
	Groupe   string `json:"group"`
	TypeSlot string `json:"slot_type"`
	Affecte  string `json:"assigned"`
}

type EntreeOrga struct {
	EmploiAncien       []CaseEmploi `json:"time_slots"`
	FichesProfs        []FicheVoeux `json:"wish_lists"`
	ModulesDisponibles []Module     `json:"available_modules"`
}

// Structure de sortie

type Affectation struct {
	Prof   string `json:"professeur"`
	Groupe string `json:"groupe"`
	Module string `json:"module"`
	Type   string `json:"type"`
	Heures int    `json:"heures"`
}

func Orga(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var donnees EntreeOrga
	err := json.NewDecoder(req.Body).Decode(&donnees)
	if err != nil {
		http.Error(res, "Entrée invalide", http.StatusBadRequest)
		return
	}

	// Map pour suivre les heures de chaque professeur
	heuresProf := make(map[string]int)
	var resultats []Affectation

	for _, slot := range donnees.EmploiAncien {
		prof := slot.Affecte
		fiche := trouverFicheProf(prof, donnees.FichesProfs)
		if fiche == nil {
			continue
		}

		// Rechercher le module correspondant dans la fiche de vœux selon la priorité
		module := trouverModule(fiche.Modules, slot.TypeSlot, donnees.ModulesDisponibles, heuresProf[prof])
		if module != nil {
			heuresProf[prof] += module.Heures
			resultats = append(resultats, Affectation{
				Prof:   prof,
				Groupe: slot.Groupe,
				Module: module.Nom,
				Type:   slot.TypeSlot,
				Heures: module.Heures,
			})
			// Supprimer ce module des disponibles
			donnees.ModulesDisponibles = retirerModule(module.Nom, module.Type, donnees.ModulesDisponibles)
		}
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(resultats)
}

func trouverFicheProf(nom string, fiches []FicheVoeux) *FicheVoeux {
	for _, f := range fiches {
		if strings.EqualFold(f.Prof, nom) {
			return &f
		}
	}
	return nil
}

func trouverModule(modules []Module, typeSlot string, disponibles []Module, heuresActuelles int) *Module {
	// Essayer priorité 1 à 3
	for prio := 1; prio <= 3; prio++ {
		for _, m := range modules {
			if strings.EqualFold(m.Type, typeSlot) && m.Priorite == prio {
				if estDisponible(m, disponibles) && heuresActuelles+m.Heures <= 24 {
					return &m
				}
			}
		}
	}
	return nil
}

func estDisponible(module Module, liste []Module) bool {
	for _, m := range liste {
		if strings.EqualFold(m.Nom, module.Nom) && strings.EqualFold(m.Type, module.Type) {
			return true
		}
	}
	return false
}

func retirerModule(nom string, typeSlot string, liste []Module) []Module {
	var nouvelle []Module
	for _, m := range liste {
		if !(strings.EqualFold(m.Nom, nom) && strings.EqualFold(m.Type, typeSlot)) {
			nouvelle = append(nouvelle, m)
		}
	}
	return nouvelle
}
