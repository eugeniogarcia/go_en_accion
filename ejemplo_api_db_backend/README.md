## Configuracion

Usamos un paquete llamado _Viper_ para gestionar la configuracion. Con Viper podemos externalizar la configuración a un archivo - yaml, o json -, y cargar las propiedades a partir de él.

```go
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
```

en main gestionamos el nombre del archivo de configuración, contemplamos el uso de una variable de entorno para indicar el entorno, y esperamos tener un archivo de configuracion con un nombre diferente según el entorno:

```go
func getConfigFileName() string {
	// buscamos la variable de entorno ENV
	env := os.Getenv("ENV")

	// si la variable de entorno ENV está definida, usamos un archivo de configuración específico
	if env != "" {
		return "runners-" + env
	}

	// si la variable de entorno ENV no está definida, usamos el archivo de configuración por defecto	
	return "runners"
}
```

la forma en que usamos las propiedades es muy sencilla. En el archivo de configuracion, un yaml en este ejemplo. Dentro del yaml definimos diferentes secciones donde ubicar las propiedades.

```yaml
###############################################################################
# Database configuration

# Connection string is in Go pq driver format:
# host=<host> port=<port> user=<databaseUser> password=<databaseUserPassword> dbname=<databaseName>

[database]

connection_string = "host=localhost port=5432 user=postgres password=postgres dbname=runners_db sslmode=disable"
max_idle_connections = 5
max_open_connections = 20
connection_max_lifetime = "60s"
driver_name = "postgres"
###############################################################################
# HTTP server configuration

[http]

server_address = ":8080"
###############################################################################
```

podemos ver una sección _database_ y varias propiedades definidas detro de ella. Para acceder a las propiedades haremos:

```go
connectionString := config.GetString("database.connection_string")
```

## instalación postgress en Ubuntu


[Guia de instalación](https://documentation.ubuntu.com/server/how-to/databases/install-postgresql/)

```sh
sudo apt install postgresql postgresql-contrib
```

Por defecto, PostgreSQL crea un usuario llamado "postgres"

sudo -u postgres psql
ALTER USER postgres WITH PASSWORD 'prueba';
CREATE USER egsmartin WITH PASSWORD 'prueba';
CREATE DATABASE runners_db OWNER egsmartin;
GRANT ALL PRIVILEGES ON DATABASE runners_db TO egsmartin;

\q


para instalar _pgadmin_, primero tenemos que configurar el repositorio:

```sh
# Install the public key for the repository (if not done previously):
curl -fsS https://www.pgadmin.org/static/packages_pgadmin_org.pub | sudo gpg --dearmor -o /usr/share/keyrings/packages-pgadmin-org.gpg

# Create the repository configuration file:
sudo sh -c 'echo "deb [signed-by=/usr/share/keyrings/packages-pgadmin-org.gpg] https://ftp.postgresql.org/pub/pgadmin/pgadmin4/apt/$(lsb_release -cs) pgadmin4 main" > /etc/apt/sources.list.d/pgadmin4.list && apt update'
```

ahora ya podemos instalarlo:

```sh
sudo apt install pgadmin4
```


systemctl status postgresql
systemctl restart postgresql
sudo systemctl start postgresql
sudo systemctl stop postgresql


sudo systemctl enable postgresql
sudo systemctl enable --now postgresql
sudo systemctl disable postgresql
sudo systemctl disable --now postgresql