package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"sync"
)

// Structure produit
type Product struct {
	ID          int
	Name        string
	Price       int
	Image       string
	Description string
	Promo       bool
	PromoPrice  int
}

// Article du panier
type CartItem struct {
	Product  Product
	Quantity int
}

// Données du panier pour template
type CartData struct {
	Items []CartItem
	Total int
}

var (
	products = []Product{
		{ID: 1, Name: "Pull Palace Capuche", Price: 148, Image: "/static/img/products/16A.webp", Promo: false, PromoPrice: 0, Description: "Tissu doux pour l'hiver."},
		{ID: 2, Name: "Pull Palace Marine", Price: 138, Image: "/static/img/products/21A.webp", Promo: true, PromoPrice: 108, Description: "Confort et style streetwear."},
		{ID: 3, Name: "Pull Palace Crew Noir", Price: 128, Image: "/static/img/products/22A.webp", Promo: false, PromoPrice: 0, Description: "Pull classique noir pour tous les jours."},
	}
	cart      = make(map[string][]CartItem)
	cartMutex = &sync.Mutex{}
)

func main() {
	// Fonction de multiplication pour template
	funcMap := template.FuncMap{
		"mul": func(a, b int) int {
			return a * b
		},
	}

	// Charger templates avec fonction mul
	StreetShopTmpl := template.Must(template.New("StreetShop.html").Funcs(funcMap).ParseFiles("StreetShop.html"))
	ArticleTmpl := template.Must(template.New("article.html").Funcs(funcMap).ParseFiles("article.html"))
	AddTmpl := template.Must(template.New("add.html").Funcs(funcMap).ParseFiles("add.html"))
	CartTmpl := template.Must(template.New("cart.html").Funcs(funcMap).ParseFiles("cart.html"))

	// Gestion session
	getSessionID := func(r *http.Request) string {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			return ""
		}
		return cookie.Value
	}

	setSessionID := func(w http.ResponseWriter) string {
		sessionID := fmt.Sprintf("session_%d", len(cart))
		http.SetCookie(w, &http.Cookie{Name: "session_id", Value: sessionID, Path: "/"})
		return sessionID
	}

	// CSS et images
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Page principale
	http.HandleFunc("/StreetShop", func(w http.ResponseWriter, r *http.Request) {
		if err := StreetShopTmpl.Execute(w, products); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Page détail produit
	http.HandleFunc("/article", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)
		for _, p := range products {
			if p.ID == id {
				ArticleTmpl.Execute(w, p)
				return
			}
		}
		http.NotFound(w, r)
	})

	// Ajouter un produit (template)
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			AddTmpl.Execute(w, nil)
			return
		}

		// POST → ajout
		name := r.FormValue("name")
		price, _ := strconv.Atoi(r.FormValue("price"))
		promo := r.FormValue("promo") == "on"
		promoPrice, _ := strconv.Atoi(r.FormValue("promoPrice"))
		image := r.FormValue("image")
		desc := r.FormValue("description")

		newID := len(products) + 1
		p := Product{
			ID:          newID,
			Name:        name,
			Price:       price,
			Promo:       promo,
			PromoPrice:  promoPrice,
			Image:       image,
			Description: desc,
		}
		products = append(products, p)
		http.Redirect(w, r, "/StreetShop", http.StatusSeeOther)
	})

	// Ajouter au panier
	http.HandleFunc("/add-to-cart", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/StreetShop", http.StatusSeeOther)
			return
		}

		id, _ := strconv.Atoi(r.FormValue("product_id"))
		qty, _ := strconv.Atoi(r.FormValue("quantity"))
		if qty <= 0 {
			qty = 1
		}

		var prod Product
		for _, p := range products {
			if p.ID == id {
				prod = p
				break
			}
		}

		sessionID := getSessionID(r)
		if sessionID == "" {
			sessionID = setSessionID(w)
		}

		cartMutex.Lock()
		found := false
		for i, item := range cart[sessionID] {
			if item.Product.ID == id {
				cart[sessionID][i].Quantity += qty
				found = true
				break
			}
		}
		if !found {
			cart[sessionID] = append(cart[sessionID], CartItem{Product: prod, Quantity: qty})
		}
		cartMutex.Unlock()

		http.Redirect(w, r, "/cart", http.StatusSeeOther)
	})

	// Supprimer du panier
	http.HandleFunc("/remove-from-cart", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/cart", http.StatusSeeOther)
			return
		}

		id, _ := strconv.Atoi(r.FormValue("product_id"))
		sessionID := getSessionID(r)
		if sessionID == "" {
			http.Redirect(w, r, "/StreetShop", http.StatusSeeOther)
			return
		}

		cartMutex.Lock()
		items := cart[sessionID]
		for i, item := range items {
			if item.Product.ID == id {
				cart[sessionID] = append(items[:i], items[i+1:]...)
				break
			}
		}
		cartMutex.Unlock()
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
	})

	// Page panier
	http.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		sessionID := getSessionID(r)
		cartMutex.Lock()
		items := cart[sessionID]
		cartMutex.Unlock()

		total := 0
		for _, item := range items {
			price := item.Product.Price
			if item.Product.Promo && item.Product.PromoPrice > 0 {
				price = item.Product.PromoPrice
			}
			total += price * item.Quantity
		}

		data := CartData{Items: items, Total: total}
		CartTmpl.Execute(w, data)
	})

	fmt.Println("🚀 Serveur démarré sur http://localhost:8080/StreetShop")
	http.ListenAndServe(":8080", nil)
}
