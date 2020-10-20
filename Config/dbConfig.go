package Config

import (
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

var CLI *mongo.Client

type DBConfig struct {
	Host     string
	Port     int
	Username string
	DBName   string
	Password string
}

func BuildConfig() *DBConfig {
	if len(os.Getenv("DB_HOST")) == 0 {
		os.Setenv("DB_HOST", "localhost:27017")
		os.Setenv("DB_USERNAME", "root")
		os.Setenv("DB_DBNAME", "auth-mongo")
		os.Setenv("DB_PASSWORD", "*******")
	}
	dbconfig := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Username: os.Getenv("DB_USERNAME"),
		DBName:   os.Getenv("DB_DBNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	return &dbconfig
}
func DbURL(BuildConfig *DBConfig) string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s",
		BuildConfig.Username,
		BuildConfig.Password,
		BuildConfig.Host,
	)
}
