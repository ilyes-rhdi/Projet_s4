package main

import (
	"Devenir_dev/cmd/database"
	"Devenir_dev/cmd/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const port = ":3000"

func main() {
	database.InitDB()
	app := mux.NewRouter()
	/*	app.Use(middelware)
		app.HandleFunc("/login", handlers.Login)
		app.HandleFunc("/Submit", handlers.Submit)
		app.HandleFunc("/Home", handlers.Home)
		app.HandleFunc("/Home/profs", handlers.List)
		app.HandleFunc("/deleteUser", handlers.DeleteUserHandler)
		app.HandleFunc("/Home/Fiche de voeux",handlers.Fiche_de_voeux)
		fmt.Println("(http://localhost:3000/login) le serveur est lancer sur ce lien ")*/
	app.HandleFunc("/organigramme", handlers.Orga).Methods("POST")
	app.HandleFunc("/organigramme", handlers.Orga).Methods("GET")

	err := http.ListenAndServe(port, app)
	if err != nil {
		log.Fatal(err)
	}
}
