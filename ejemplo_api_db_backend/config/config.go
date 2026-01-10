package config

import (
	"log"

	"github.com/spf13/viper"
)

func InitConfig(fileName string) *viper.Viper {
	// instanciamos viper
	config := viper.New()

	// indicamos el nombre del archivo de configuración
	config.SetConfigName(fileName)

	// indicamos las rutas donde buscar el archivo de configuración
	config.AddConfigPath(".")
	config.AddConfigPath("$HOME")

	// leemos el archivo de configuración
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal("Error while parsing configuration file", err)
	}

	return config
}
