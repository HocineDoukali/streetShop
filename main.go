package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func main() {
	// Charger le template
	tmpl, err := template.ParseFiles("streetShop.html")
	if err != nil {
		fmt.Println("Erreur lors du chargement du template :", err)
		os.Exit(1)
	}

	// Servir les fichiers statiques (CSS + images)
	http.Handle("/streetShop.css", http.FileServer(http.Dir(".")))
	http.Handle("/1.png", http.FileServer(http.Dir(".")))
	http.Handle("/19A.webp", http.FileServer(http.Dir(".")))
	http.Handle("/21A.webp", http.FileServer(http.Dir(".")))
	http.Handle("/22A.webp", http.FileServer(http.Dir(".")))
	http.Handle("/16A.webp", http.FileServer(http.Dir(".")))
	http.Handle("/18A.webp", http.FileServer(http.Dir(".")))
	http.Handle("/33B.webp", http.FileServer(http.Dir(".")))
	http.Handle("/34B.webp", http.FileServer(http.Dir(".")))

	// Route principale
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"Title":    "Produits - StreetShop",
			"SiteName": "StreetShop",
		}
		tmpl.Execute(w, data)
	})

	// Route principale
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"Title":    "Produits - StreetShop",
			"SiteName": "StreetShop",
		}
		tmpl.Execute(w, data)
	})

	// Route /add (pour plus tard)
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<h1>Page d’ajout de produit (à venir)</h1>")
	})

	// Lancer le serveur
	fmt.Println("🚀 Serveur StreetShop démarré sur http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erreur serveur :", err)
		os.Exit(1)
	}
}
