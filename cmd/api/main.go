package main

import (
	"net/http"
	"github.com/ilyes-rhdi/Projet_s4/internal/api/rooter"
	"github.com/ilyes-rhdi/Projet_s4/internal/api/middleware"
	"github.com/ilyes-rhdi/Projet_s4/internal/database"
	"os"
	"log"
)

func main() {
	database.InitDB()
	app := rooter.NewRouter()
	app.Use(middleware.JwtMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Port par défaut localement
	}

	err := http.ListenAndServe("0.0.0.0:"+port, app)
	if err != nil {
		log.Fatal(err)
	}
}
  



 