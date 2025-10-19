package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func main() {
	// Charger le template
	tmpl, err := template.ParseFiles("StreetShop.html")
	if err != nil {
		fmt.Println("Erreur lors du chargement du template :", err)
		os.Exit(1)
	}

	// Servir le CSS
	http.Handle("/StreetShop.css", http.FileServer(http.Dir(".")))

	// Servir les images dans static/img
	fsImages := http.FileServer(http.Dir("./static/img"))
	http.Handle("/static/img/", http.StripPrefix("/static/img/", fsImages))

	// Route principale
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"Title": "Produits - StreetShop",
		}
		tmpl.Execute(w, data)
	})

	// Lancer le serveur
	fmt.Println("🚀 Serveur StreetShop démarré sur http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erreur serveur :", err)
		os.Exit(1)
	}
}
