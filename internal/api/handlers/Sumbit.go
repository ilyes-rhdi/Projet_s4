package handlers

import (
	"Devenir_dev/internal/database"
    "Devenir_dev/internal/api/models"
	"fmt"
	"net/http"
    "Devenir_dev/pkg"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)


func Submit(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
        // Render the login page (e.g., HTML page)
        utils.Rendertemplates(res, "Submit",nil)
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
    // Create a User struct from form data
    user := models.User{
        Username:     req.FormValue("username"),
        Email:    req.FormValue("email"),
        PasswordHash:      req.FormValue("password"),
        Role:    req.FormValue("role"),
        FullName: req.FormValue("name"),
    }

    // Validate and sanitize input
    utils.ValidateInput(user)
    utils.SanitizeInput(&user)
    password,_:=bcrypt.GenerateFromPassword([]byte(user.PasswordHash),14)

    // Prepare SQL statement
    stmt, err := db.Prepare("INSERT INTO users(name, email, password, isAdmin ,Speciality ,Year_entrance,Grade) VALUES(?, ?, ?, ?, ?, ?, ?)")
    if err != nil {
        fmt.Println("Database prepare error:", err)  // Log the actual error
        http.Error(res, "Database error", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    _, err = stmt.Exec(user.Username, user.Email, password ,user.Role,user.FullName)
    if err != nil {
        fmt.Println("Database exec error:", err)  
        http.Error(res, "Failed to insert user into database", http.StatusInternalServerError)
        return
    }

    // Send success response
        http.Redirect(res, req, "/Home", http.StatusFound)
    
    
}
