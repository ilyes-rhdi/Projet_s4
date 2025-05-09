package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"Devenir_dev/internal/api/models"
    "Devenir_dev/pkg"
 	"Devenir_dev/internal/database"
)

func AddTeacher(res http.ResponseWriter, req *http.Request) {
    if req.Method == http.MethodGet {
        // Render the add teacher page (e.g., HTML page)
        utils.Rendertemplates(res, "AddTeacher", nil)
        return
    }
    if req.Method != http.MethodPost {
        http.Error(res, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    db := database.GetDB()
    err := req.ParseForm()
    if err != nil {
        http.Error(res, "Error parsing form data", http.StatusBadRequest)
        return
    }

    // Create a Teacher struct from form data
    teacher := models.Teacher{
        Name:        req.FormValue("name"),
        SpecialtyID: uint(atoiOrDefault(req.FormValue("specialty_id"), 0)), // Convert int to uint
        MaxHours:    atoiOrDefault(req.FormValue("max_hours"), 0),
        CurrentHours: atoiOrDefault(req.FormValue("current_hours"), 0),
    }

    // Validate and sanitize input
    utils.SanitizeTeacherInput(&teacher) 

    // Prepare SQL statement
    // Insert the teacher into the database using GORM
    if err := db.Create(&teacher).Error; err != nil {
        fmt.Println("Database insert error:", err)
        http.Error(res, "Failed to insert teacher into database", http.StatusInternalServerError)
        return
    }

    // Send success response
    http.Redirect(res, req, "/teachers", http.StatusFound)
}

// Helper function to convert string to int with a default value
func atoiOrDefault(value string, defaultValue int) int {
    result, err := strconv.Atoi(value)
    if err != nil {
        return defaultValue
    }
    return result
}
func EditProfile(w http.ResponseWriter, r *http.Request) {
    // Logic to edit the profile
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Profile edited successfully"))
}