package handlers
import(
   "net/http"
   "Devenir_dev/cmd/database"
)


func Fiche_de_voeux(res http.ResponseWriter, req *http.Request){
	if req.Method == http.MethodGet {
        Rendertemplates(res, "Fiche",nil)
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
    session, _ := store.Get(req, "session-name")
	profID, ok := session.Values["user_id"].(int)  
    if !ok {
        http.Error(res, "Utilisateur non connecté", http.StatusUnauthorized)
        return
    }
    // Vérifier si le prof a déjà 3 choix enregistrés
    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM fiche WHERE prof_id = ?", profID).Scan(&count)
    if err != nil {
        http.Error(res, "Erreur lors de la vérification", http.StatusInternalServerError)
        return
    }
    if count >= 3 {
        http.Error(res, "Vous avez déjà soumis 3 choix.", http.StatusForbidden)
        return
    }

    // Récupérer les 3 choix depuis le formulaire
    choices := []Fiche{
        {Fillier: req.FormValue("fillier1"), Anner: req.FormValue("anner1"), Tp: formBool(req, "tp1"), Td: formBool(req, "td1"), Cour: formBool(req, "cour1"), Priority: 1},
        {Fillier: req.FormValue("fillier2"), Anner: req.FormValue("anner2"), Tp: formBool(req, "tp2"), Td: formBool(req, "td2"), Cour: formBool(req, "cour2"), Priority: 2},
        {Fillier: req.FormValue("fillier3"), Anner: req.FormValue("anner3"), Tp: formBool(req, "tp3"), Td: formBool(req, "td3"), Cour: formBool(req, "cour3"), Priority: 3},
    }

    stmt, err := db.Prepare("INSERT INTO fiche (prof_id, fillier, anner, tp, td, cour, priority) VALUES (?, ?, ?, ?, ?, ?, ?)")
    if err != nil {
        http.Error(res, "Erreur de préparation de la requête", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    for _, choice := range choices {
        if choice.Fillier != "" { // Vérifier si le choix est rempli
            _, err = stmt.Exec(profID, choice.Fillier, choice.Anner, choice.Tp, choice.Td, choice.Cour, choice.Priority)
            if err != nil {
                http.Error(res, "Erreur lors de l'insertion", http.StatusInternalServerError)
                return
            }
        }
    }


}