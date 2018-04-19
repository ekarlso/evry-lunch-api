package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"

	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"github.com/ekarlso/evry-lunch-api/api/handlers"
	"github.com/ekarlso/evry-lunch-api/db"
)

func MakeCors() *cors.Cors {
	allowedHeaders := viper.GetStringSlice("cors_allowed_headers")
	allowedOrigins := viper.GetStringSlice("cors_allowed_origins")
	allowedMethods := viper.GetStringSlice("cors_allowed_methods")

	logrus.Debugf("CORS config: %s - %s - %s", allowedOrigins, allowedMethods, allowedHeaders)

	return cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: allowedMethods,
		AllowedHeaders: allowedHeaders,
	})
}

// SetupServerFlags - registers flags for server
func SetFlags() {
	pflag.String("path_base", "", "Base path")
	pflag.String("bind", ":5000", "Port to be running at.")
	pflag.String("db_connection", "user=postgres password=postgres dbname=postgres sslmode=disable host=localhost", "DB Connection.")
}

// func SeedData() {
// 	allergens := []models.Allergen{
// 		models.Allergen{Name: "Nuts", ID: uuid.NewV4()},
// 		models.Allergen{Name: "Gluten", ID: uuid.NewV4()},
// 	}

// 	for _, i := range allergens {
// 		db.Connection.Create(&i)
// 	}

// 	bologneseDish := models.Dish{
// 		Name: "Pasta Bolognese",
// 	}
// 	db.Connection.Create(&bologneseDish)

// 	bolognose

// 	tacoDish := models.Dish{
// 		Name: "Taco",
// 	}
// 	db.Connection.Create(&tacoDish)

// 	//err := db.Connection.Create(allergens)
// 	//fmt.Println(err)
// }

func initOpts() {
	//viper.SetEnvPrefix("bots")
	viper.AutomaticEnv()

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func main() {
	SetFlags()
	initOpts()

	logrus.SetLevel(logrus.DebugLevel)
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negronilogrus.NewMiddlewareFromLogger(log, "web"))

	cors := MakeCors()
	n.Use(cors)

	router := mux.NewRouter()
	menuRouter := handlers.MakeItemsRouter()
	router.PathPrefix("/menu").Handler(menuRouter)

	prefix := viper.GetString("path_base")
	if prefix != "" {
		n.UseHandler(http.StripPrefix(viper.GetString(prefix), router))
	} else {
		n.UseHandler(router)
	}

	err := db.Init()
	if err != nil {
		panic(err)
	}

	ex, err := os.Executable()
	exPath := filepath.Dir(ex)

	err = db.Migrate(exPath + "/migrations")
	if err != nil {
		fmt.Println(err.Error())
	}

	db.Connection.SetLogger(log)

	// SeedData()
	server := http.Server{
		Handler: n,
		Addr:    ":5000",
	}

	server.ListenAndServe()
}
