package main

import (
	"bytes"
	"log"
	"net/http"
	"text/template"
	"encoding/json"
	"os"
)

var _temp = template.Must(template.ParseGlob("./templates/*.html"))

type Project struct {
    Id              string   `json:"id"`
    Titre           string   `json:"titre"`
    Image           string   `json:"image"`
    Technologies    string   `json:"technologies"`
    Introduction    string   `json:"introduction"`
    Presentation    string   `json:"presentation"`
    Fonctionnalites []string `json:"fonctionnalites"`
    Depot           string   `json:"depot"`
    Video           string   `json:"video"`
}

var listProjects []Project

func main() {

	jErr := loadProjects()
	if jErr != nil {
		log.Fatalf("Erreur chargement JSON : %v", jErr)
	}

	router := http.NewServeMux()

	router.HandleFunc("/{$}", HomeHandler)
	router.HandleFunc("/projects", ProjectsHandler)
	router.HandleFunc("/projects/{id}", ProjectsByIdHandler)
	router.HandleFunc("/", NotfoundHandler)

	fileServer := http.FileServer(http.Dir("./assets"))
	router.Handle("/static/", http.StripPrefix("/static/", fileServer))

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatalf("Erreur (Serveur) : Impossible de lancer le serveur : %v", err)
		return
	}

}

func loadProjects() error {
    file, err := os.Open("data.json")
    if err != nil {
        return err
    }
    defer file.Close()

    err = json.NewDecoder(file).Decode(&listProjects)
    return err
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {

	buffer := bytes.Buffer{}

	errRender := _temp.ExecuteTemplate(&buffer, name, data)
	if errRender != nil {
		log.Printf("Erreur (Template) : Impossible de rendre le template %s : %v", name, errRender)
		http.Error(w, "Erreur chargement de la page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buffer.WriteTo(w)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request : %s %s", r.Method, r.URL.Path)

	renderTemplate(w, "index", nil)
}

func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request : %s %s", r.Method, r.URL.Path)

	renderTemplate(w, "projects", listProjects)
}

func ProjectsByIdHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Request : %s %s", r.Method, r.URL.Path)

    id := r.PathValue("id")

    for _, product := range listProjects {
        if product.Id == id {
            renderTemplate(w, "project-details", product)
            return
        }
    }

    http.Error(w, "Projet introuvable...", http.StatusNotFound)
}

func NotfoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request : %s %s", r.Method, r.URL.Path)

	renderTemplate(w, "error", nil)
}

