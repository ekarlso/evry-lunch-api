package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/ekarlso/evry-lynsjer/db"
	"github.com/ekarlso/evry-lynsjer/models"
	"github.com/gorilla/mux"
)

type Dish struct {
	Name      string   `json:"name"`
	Allergens []string `json:"allergens"`
}

type Menu struct {
	Day    time.Time `json:"day"`
	Dishes []Dish    `json:"dishes"`
}

func GetMenu(w http.ResponseWriter, r *http.Request) {
	// db.Connection.Order("created_at desc").Find(&items)
	var menuItems []Menu

	var entries []models.MenuEntry
	db.Connection.Find(&entries)

	for _, entry := range entries {
		e := Menu{
			Day: entry.Day,
		}

		dishes := []Dish{}
		var dishRefs []models.Dish
		db.Connection.Where(models.Dish{ID: entry.DishID}).Find(&dishRefs)
		for _, dishRef := range dishRefs {
			var allergenRefs []models.Allergen
			db.Connection.Model(&dishRef).Related(&allergenRefs, "Allergens")

			allergens := []string{}
			for _, a := range allergenRefs {
				allergens = append(allergens, a.Name)
			}

			dishes = append(dishes, Dish{Name: dishRef.Name, Allergens: allergens})
		}

		e.Dishes = dishes

		menuItems = append(menuItems, e)
	}

	WriteJson(w, menuItems)
}

func CreateMenu(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	var menu Menu
	err := json.NewDecoder(r.Body).Decode(&menu)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	for _, dish := range menu.Dishes {
		dishRef := models.Dish{Name: dish.Name}
		if db.Connection.Where(models.Dish{Name: dish.Name}).First(&dishRef).RecordNotFound() {
			dishRef.ID = uuid.NewV4()
			db.Connection.Create(&dishRef)
		}

		var dishAllergens = []models.Allergen{}
		for _, allergenName := range dish.Allergens {
			// Get or create all allergens the current dish has...
			allergenRef := models.Allergen{Name: allergenName}
			if db.Connection.Where(allergenRef).First(&allergenRef).RecordNotFound() {
				allergenRef.ID = uuid.NewV4()
				db.Connection.Create(&allergenRef)
			}

			dishAllergens = append(dishAllergens, allergenRef)
		}

		for _, allergen := range dishAllergens {
			dishAllergen := models.DishAllergen{DishID: dishRef.ID, AllergenID: allergen.ID}

			if db.Connection.Where(dishAllergen).First(&dishAllergen).RecordNotFound() {
				dishAllergen.ID = uuid.NewV4()
				db.Connection.Create(&dishAllergen)
			}
		}

		menuEntry := models.MenuEntry{Day: menu.Day, DishID: dishRef.ID}
		if db.Connection.Where(models.MenuEntry{Day: menu.Day, DishID: dishRef.ID}).First(&menuEntry).RecordNotFound() {
			menuEntry.ID = uuid.NewV4()
			db.Connection.Create(&menuEntry)
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func MakeItemsRouter() *mux.Router {
	router := mux.NewRouter()
	router.Methods("GET").Path("/menu").HandlerFunc(GetMenu)
	router.Methods("POST").Path("/menu").HandlerFunc(CreateMenu)
	return router
}
