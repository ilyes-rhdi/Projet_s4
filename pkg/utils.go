package utils

import (
	"Devenir_dev/internal/api/models"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Pagedata struct {
	Currentuser models.User
	Users       []models.User
}

// Rendertemplates charge et affiche les templates
func Rendertemplates(res http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("C:\\Users\\PC\\OneDrive\\Documents\\futur\\Devenir_dev\\templates\\" + tmpl + ".page.tmpl")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(res, data)
	if err != nil {
		http.Error(res, "Error executing template", http.StatusInternalServerError)
		fmt.Println("Error executing template:", err)
	}
}

// VerifyUser vérifie l'authentification de l'utilisateur avec GORM
func VerifyUser(db *gorm.DB, identifier, password string) (bool, models.Role, string) {
	var user models.User

	// Vérifie si l'identifiant est un email ou un nom d'utilisateur
	if strings.Contains(identifier, "@") {
		// Si c'est un email
		if err := db.Where("email = ?", identifier).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return false, "", "User not found."
			}
			log.Println("GORM Error:", err)
			return false, "", "Database error."
		}
	} else {
		// Si c'est un nom d'utilisateur
		if err := db.Where("nom = ?", identifier).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return false, "", "User not found."
			}
			log.Println("GORM Error:", err)
			return false, "", "Database error."
		}
	}

	// Vérifie le mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, "", "Incorrect password."
	}

	return true, user.Role, "User verified."
}

// ValidateInput vérifie la validité des champs utilisateur
func ValidateInput(user models.User) (bool, string) {
	// Vérification des champs vides
	if user.Nom == "" || user.Prenom == "" || user.Password == "" || user.Role == "" || user.Email == "" {
		return false, "All fields (nom, prenom, email, password, role) are required."
	}

	// Vérification de l'email avec une expression régulière
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		return false, "Invalid email format."
	}

	// Vérification de la longueur du mot de passe (ex: minimum 6 caractères)
	if len(user.Password) < 6 {
		return false, "Password must be at least 6 characters long."
	}

	return true, ""
}
func sanitizeRole(role models.Role) models.Role {
	switch role {
	case models.Admin, models.Professeur, models.Responsable:
		return models.Role(role) // Rôle valide
	default:
		// Retourne un rôle par défaut si le rôle est invalide
		return models.Professeur
	}
}

func SanitizeInput(user *models.User) {
	re := regexp.MustCompile("<.*?>")

	user.Nom = clean(user.Nom, re)
	user.Prenom = clean(user.Prenom, re)
	user.Password = clean(user.Password, re)
	user.Email = clean(user.Email, re)
	user.Role = sanitizeRole(user.Role)
}

// clean supprime les balises HTML et les espaces inutiles
func clean(s string, re *regexp.Regexp) string {
	return re.ReplaceAllString(strings.TrimSpace(s), "")
}

// FormBool vérifie si une case à cocher est activée dans un formulaire
func FormBool(r *http.Request, key string) bool {
	return r.FormValue(key) == "on"
}

// FindTeacher searches for a teacher by ID in the given list of teachers
func FindTeacher(teacherID uint, teachers []models.Teacher) *models.Teacher {
	for _, t := range teachers {
		if t.ID == teacherID {
			return &t
		}
	}
	return nil
}

// FindModuleForTeacher finds an appropriate module for a teacher based on their wishes and availability
func FindModuleForTeacher(teacherID uint, slotType string, wishes []models.Voeux, available []models.Module, currentHours int) *models.Module {
	// Try priorities 1 to 3
	for prio := 1; prio <= 3; prio++ {
		for _, wish := range wishes {
			if wish.TeacherID == teacherID && wish.Priority == prio {
				// Check if teacher wants this type of class
				if (slotType == "cours" && wish.Cours) ||
					(slotType == "td" && wish.Td) ||
					(slotType == "tp" && wish.Tp) {
					// Find the module in available modules
					for _, module := range available {
						if module.ID == wish.ModuleID {
							hours := GetHoursForType(&module, slotType)
							if hours > 0 && currentHours+hours <= 24 {
								return &module
							}
						}
					}
				}
			}
		}
	}
	return nil
}

// GetHoursForType returns the number of hours for a specific type of class in a module
func GetHoursForType(module *models.Module, slotType string) int {
	switch slotType {
	case "cours":
		return module.VolumeCours
	case "td":
		return module.VolumeTD
	case "tp":
		return module.VolumeTP
	default:
		return 0
	}
}

// RemoveModule removes a module from the list of available modules
func RemoveModule(moduleID uint, modules []models.Module) []models.Module {
	var result []models.Module
	for _, m := range modules {
		if m.ID != moduleID {
			result = append(result, m)
		}
	}
	return result
}

// SaveOrganigram saves the current organigram as a template or regular version
func SaveOrganigram(db *gorm.DB, user models.User, input struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Year        string                  `json:"year"`
	IsTemplate  bool                    `json:"is_template"`
	Data        models.OrganigramOutput `json:"data"`
}) (*models.SavedOrganigram, error) {
	// Convert data to JSON string
	dataJSON, err := json.Marshal(input.Data)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la conversion des données: %v", err)
	}

	savedOrg := models.SavedOrganigram{
		Name:        input.Name,
		Description: input.Description,
		Year:        input.Year,
		IsTemplate:  input.IsTemplate,
		CreatedBy:   user.ID,
		Data:        string(dataJSON),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(&savedOrg).Error; err != nil {
		return nil, fmt.Errorf("erreur lors de la sauvegarde: %v", err)
	}

	return &savedOrg, nil
}

// GetOrganigramTemplates retrieves all organigram templates
func GetOrganigramTemplates(db *gorm.DB) ([]models.SavedOrganigram, error) {
	var templates []models.SavedOrganigram
	if err := db.Where("is_template = ?", true).Find(&templates).Error; err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des templates: %v", err)
	}
	return templates, nil
}

// GetOrganigramByID retrieves a specific organigram by ID
func GetOrganigramByID(db *gorm.DB, id string) (*models.SavedOrganigram, *models.OrganigramOutput, error) {
	var organigram models.SavedOrganigram
	if err := db.First(&organigram, id).Error; err != nil {
		return nil, nil, fmt.Errorf("organigramme non trouvé: %v", err)
	}

	// Parse the JSON data back into OrganigramOutput
	var output models.OrganigramOutput
	if err := json.Unmarshal([]byte(organigram.Data), &output); err != nil {
		return nil, nil, fmt.Errorf("erreur lors de la conversion des données: %v", err)
	}

	return &organigram, &output, nil
}

// UpdateOrganigram updates an existing organigram
func UpdateOrganigram(db *gorm.DB, user models.User, input struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Data        models.OrganigramOutput `json:"data"`
}) error {
	// Convert data to JSON string
	dataJSON, err := json.Marshal(input.Data)
	if err != nil {
		return fmt.Errorf("erreur lors de la conversion des données: %v", err)
	}

	// Update the organigram
	result := db.Model(&models.SavedOrganigram{}).
		Where("id = ?", input.ID).
		Updates(map[string]interface{}{
			"name":        input.Name,
			"description": input.Description,
			"data":        string(dataJSON),
			"updated_at":  time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("erreur lors de la mise à jour: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("organigramme non trouvé")
	}

	return nil
}

// DeleteOrganigram deletes an organigram by ID
func DeleteOrganigram(db *gorm.DB, id string) error {
	result := db.Delete(&models.SavedOrganigram{}, id)
	if result.Error != nil {
		return fmt.Errorf("erreur lors de la suppression: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("organigramme non trouvé")
	}
	return nil
}
