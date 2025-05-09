package utils

import (
	"Devenir_dev/internal/api/models"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Pagedata struct {
	Currentuser models.User
	Users       []models.User
}

// Rendertemplates charge et affiche les templates
func Rendertemplates(res http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(fmt.Sprintf("./templates/%s.page.tmpl", tmpl))
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
	case models.StaffAdmin, models.Professeur, models.ChefDep:
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
func FindTeacher(teacherID int, teachers []models.Teacher) *models.Teacher {
	for _, t := range teachers {
		if teacherID >= 0 && t.ID == uint(teacherID) {
			return &t
		}
	}
	return nil
}

// FindModuleForTeacher finds an appropriate module for a teacher based on their wishes and availability
func FindModuleForTeacher(teacherID int, slotType string, wishes []models.Voeux, available []models.Module, currentHours int) *models.Module {
	// Try priorities 1 to 3
	for prio := 1; prio <= 3; prio++ {
		for _, wish := range wishes {
			if teacherID >= 0 && wish.TeacherID == uint(teacherID) && wish.Priority == prio {
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
func RemoveModule(moduleID int, modules []models.Module) []models.Module {
	var result []models.Module
	for _, m := range modules {
		if m.ID != uint(moduleID) {
			result = append(result, m)
		}
	}
	return result
}

// ValidateTeacherInput validates the input for a Teacher struct
func ValidateTeacherInput(teacher models.Teacher) (bool, string) {
	// Check for empty fields
	if teacher.Name == "" || teacher.SpecialtyID == 0 || teacher.MaxHours == 0 {
		return false, "All fields (name, specialty ID, max hours) are required."
	}

	// Ensure max hours is a positive number
	if teacher.MaxHours < 0 {
		return false, "Max hours must be a positive number."
	}

	return true, ""
}

// SanitizeTeacherInput sanitizes the input for a Teacher struct
func SanitizeTeacherInput(teacher *models.Teacher) {
	re := regexp.MustCompile("<.*?>")

	teacher.Name = clean(teacher.Name, re)
	// No need to sanitize numeric fields like SpecialtyID, MaxHours, etc.
}

// SanitizeModuleInput sanitizes the input for a Module struct
