package server

import (
	"database/sql"
	"log"

	"github.com/spf13/viper"
)

func InitDatabase(config *viper.Viper) *sql.DB {
	// cadena de conexión a la base de datos
	connectionString := config.GetString("database.connection_string")
	// configuramos las conexones a mantener abiertas, máximas y el tiempo máximo de vida de una conexión
	maxIdleConnections := config.GetInt("database.max_idle_connections")
	maxOpenConnections := config.GetInt("database.max_open_connections")
	connectionMaxLifetime := config.GetDuration("database.connection_max_lifetime")
	// obtenemos el nombre del driver de base de datos
	driverName := config.GetString("database.driver_name")

	if connectionString == "" {
		log.Fatalf("Database connectin string is missing")
	}

	// creamos la conexión a la base de datos
	dbHandler, err := sql.Open(driverName, connectionString)
	if err != nil {
		log.Fatalf("Error while initializing database: %v", err)
	}

	dbHandler.SetMaxIdleConns(maxIdleConnections)
	dbHandler.SetMaxOpenConns(maxOpenConnections)
	dbHandler.SetConnMaxLifetime(connectionMaxLifetime)

	// nos conectamos a la base de datos
	err = dbHandler.Ping()
	if err != nil {
		dbHandler.Close()
		log.Fatalf("Error while validating database: %v", err)
	}

	return dbHandler
}
