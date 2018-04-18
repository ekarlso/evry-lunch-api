package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ekarlso/evry-lynsjer/db"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initFunc() {
	viper.AutomaticEnv()
	logrus.SetLevel(logrus.DebugLevel)

	err := db.Init()

	if err != nil {
		panic(err)
	}

}

func clearTable() {
	db.Connection.Exec("DELETE FROM menu_entries")
	db.Connection.Exec("DELETE FROM dish_allergens")
	db.Connection.Exec("DELETE FROM dishes")
	db.Connection.Exec("DELETE FROM allergens")
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(handler http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}

func checkJSONHeader(t *testing.T, rr *httptest.ResponseRecorder) {
	if rr.HeaderMap.Get("Content-Type") != "application/json; charset=UTF-8" {
		t.Fatal()
	}
}

func checkIsEmpty(t *testing.T, rr *httptest.ResponseRecorder) {
	if body := rr.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreateVisit(t *testing.T) {
	initFunc()
	clearTable()
	// Create a visit in DB

	router := MakeItemsRouter()

	menu := &Menu{
		Day: time.Time{},
		Dishes: []Dish{
			Dish{
				Name:      "Pasta Bolognese",
				Allergens: []string{"tomatoes"},
			},
		},
	}
	refSerialized, _ := json.Marshal(menu)
	req, _ := http.NewRequest("POST", "/menu", bytes.NewReader(refSerialized))
	rr := executeRequest(router, req)

	checkResponseCode(t, http.StatusCreated, rr.Code)
	// checkJSONHeader(t, rr)

	// var decoded *Menu
	// json.Unmarshal(rr.Body.Bytes(), &decoded)
}
