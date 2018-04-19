package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Allergen - an allergen..
type Allergen struct {
	ID          uuid.UUID `json:"id" sql:"type:varchar(36)"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name" sql:"type:varchar(40)"`
	Description string    `json:"description" sql:"type:text"`
}

// Dish - a predefined meal to use when composing a menu
type Dish struct {
	ID          uuid.UUID `json:"id" sql:"type:varchar(36)"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`

	Allergens []Allergen `gorm:"many2many:dish_allergens"`
}

// DishAllergen - Which allergens does a dish contain ?
type DishAllergen struct {
	ID         uuid.UUID `json:"id" sql:"type:varchar(36)"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	AllergenID uuid.UUID `json:"allergen_id" sql:"type:varchar(36)"`
	DishID     uuid.UUID `json:"dish_id" sql:"type:varchar(36)"`
}

// MenuEntry - An entry in the menu for a given day etc...
type MenuEntry struct {
	ID        uuid.UUID `json:"id" sql:"type:varchar(36)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Which day is this for
	Day    time.Time `json:"date" sql:"type:date"`
	DishID uuid.UUID `json:"dish_id" sql:"type:varchar(36)"`
}
