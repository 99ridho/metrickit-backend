package db

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

// InitailizeDatabase is used to initialize database connection, and return the connection
func Initialize() *sqlx.DB {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	databaseName := viper.GetString("database.name")
	driverName := viper.GetString("database.driverName")

	dsn := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?parseTime=true",
		username, password, host, port, databaseName)

	db, err := sqlx.Open(driverName, dsn)

	if err != nil {
		log.Fatalln(err)
	}

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(err)
	}

	return db
}
