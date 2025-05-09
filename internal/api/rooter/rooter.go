package rooter

import (
    "Devenir_dev/internal/api/handlers"
    "github.com/gorilla/mux"
    "Devenir_dev/internal/api/middleware"
    "Devenir_dev/internal/api/models"
    "net/http"

)

func NewRouter() *mux.Router {
    router := mux.NewRouter()

    // Définition des routes
    router.HandleFunc("/login", handlers.Login).Methods("GET", "POST")
    router.HandleFunc("/submit", handlers.Submit).Methods("POST")
    router.HandleFunc("/home", handlers.Home).Methods("GET")
    router.HandleFunc("/home/fiche-de-voeux", handlers.Fiche_de_voeux).Methods("GET", "POST")
    router.HandleFunc("/teachers", handlers.AddTeacher).Methods("GET")
    router.HandleFunc("/addTeacher", handlers.AddTeacher).Methods("POST")
    router.HandleFunc("/updateTeacher", handlers.UpdateTeacher).Methods("PUT")
    router.HandleFunc("/deleteTeacher", handlers.DeleteTeacher).Methods("DELETE")
    // routes for chef departement
     router.Handle("/organigram/edit", middleware.RequireRole(models.ChefDep, "edit_organigram")(http.HandlerFunc(handlers.EditOrganigram))).Methods("POST")

    // Routes for StaffAdmin
    router.Handle("/comments", middleware.RequireRole(models.StaffAdmin, "leave_comments")(http.HandlerFunc(handlers.LeaveComments))).Methods("POST")

    // Routes for Professeur
    router.Handle("/profile", middleware.RequireRole(models.Professeur, "edit_profile")(http.HandlerFunc(handlers.EditProfile))).Methods("PUT")
    return router
}