package db

import (
	"github.com/ekarlso/evry-lunch-api/models"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // Import pg driver
	"github.com/spf13/viper"
)

var Connection *gorm.DB

func Init() error {
	dsn := viper.GetString("db_connection")
	logrus.Debugln(dsn)

	db, err := gorm.Open("postgres", dsn)
	logrus.Debugln("Init DB")

	// Links to Azure limits https://docs.microsoft.com/en-us/azure/postgresql/concepts-limits
	db.DB().SetMaxOpenConns(25)
	if err != nil {
		return err
	}
	db.LogMode(true)
	Connection = db
	return nil
}

func Migrate(migrationsPath string) error {
	driver, err := postgres.WithInstance(Connection.DB(), &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}

func GetOrCreateAllergen(ref models.Allergen) {
	if Connection.Where("name = ?", ref.Name).First(&ref).RecordNotFound() {
		ref.ID = uuid.NewV4()
		Connection.Create(&ref)
	}
}

func GetOrCreateDish(ref models.Dish) {
	if Connection.Where("name = ?", ref.Name).First(&ref).RecordNotFound() {
		ref.ID = uuid.NewV4()
		Connection.Create(&ref)
	}
}
