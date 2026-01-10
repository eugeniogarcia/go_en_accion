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

[Guia de instalación](https://documentation.ubuntu.com/server/how-to/databases/install-postgresql/). Instalamos postgres:

```sh
sudo apt install postgresql postgresql-contrib
```

Por defecto, PostgreSQL crea un usuario llamado "postgres". Podemos abrir una sesión de _psql_ para cambiar la contraseña del usuario, o para crear un nuevo usuario y/o base de datos:

```sh
sudo -u postgres psql
```

y luego en la sesión de psql:

```psql
ALTER USER postgres WITH PASSWORD 'prueba';

CREATE USER egsmartin WITH PASSWORD 'prueba';
CREATE DATABASE runners_db OWNER egsmartin;
GRANT ALL PRIVILEGES ON DATABASE runners_db TO egsmartin;
```

salimos con `\q`. Para instalar _pgadmin_, primero tenemos que configurar el repositorio:

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

### Servicio postgres

Postgres se ejecuta como un servicio. Podemos operar el servicio con los siguientes comandos:

```sh
systemctl status postgresql
systemctl restart postgresql
systemctl start postgresql
systemctl stop postgresql
```

con esto vemos el estado del servicio, lo re-arrancamos, arrancamos o paramos. Podemos habilitar que el servicio se arranque automáticamente al arrancar el linux:

```sh
sudo systemctl enable postgresql
sudo systemctl enable --now postgresql
```

o deshabilitarlo:

```sh
sudo systemctl disable postgresql
sudo systemctl disable --now postgresql
```

con la opción `--now` hacemos que el cambio se haga de forma inmediata

## api

Para construir la api identificamos tres capas independientes con el objeto de separar _concerns_:

- Controller. Esta capa proporciona las funciones que se asocian a los recursos de la api
- Servicio. Esta capa implementa la lógica de negocio. La función del controller hará uso de las funciones de negocio implementadas en el servicio. El acceso a datos se abstrae en la siguiente capa, el repositorio
- Repository. Implementa la lógica de acceso a datos. Se han considerado cuatro implementaciones alternativas:
	- Usando Postgres
	- Utilizando MySql
	- Usando MongoDB
	- utilizando DynamoDB

La construcción del proxy la hacemos utilizando el paquete Gin. Con Gin definimos para cada recurso/método el controlador aosociado.

La forma en la que se implementa cada capa es la misma:

- Se usa una función factoría para crear una instancia del objeto (controller, servicio o repositorio)
- La capa se modela con un struct que contiene todas las propiedades necesarias, así como los métodos necesarios para gestionar la capa

### Controller

Usamos Gin para implementar el router:

```go
// instancia el router de Gin...
router := gin.Default()

// ...y define las rutas y los controladores asociados
router.POST("/runner", runnersController.CreateRunner)
router.PUT("/runner", runnersController.UpdateRunner)
router.DELETE("/runner/:id", runnersController.DeleteRunner)
router.GET("/runner/:id", runnersController.GetRunner)
router.GET("/runner", runnersController.GetRunnersBatch)

router.POST("/result", resultsController.CreateResult)
router.DELETE("/result/:id", resultsController.DeleteResult)

router.POST("/login", usersController.Login)
router.POST("/logout", usersController.Logout)
```

los métodos asociados a cada recurso son el controler. Estos métodos tienen la misma firma `CreateRunner(ctx *gin.Context) {`. El argumento es el contexto Gin. El contexto se usa para acceder a todos los elementos de la request (cabecera, payload, path parameters y query parameters). El contexto tambien sirve para crear la respuesta. Para ello proporciona diferentes métodos:

```go
// recuperamos una cabecera de la petición
accessToken := ctx.Request.Header.Get("Token")

// contruye una respuesta con el http status code y el payload
ctx.JSON(responseErr.Status, responseErr)

// responde con el http status code y un payload, y detiene la ejecución del handler
ctx.AbortWithError(http.StatusInternalServerError, err)

// responde con el http status code
ctx.Status(http.StatusUnauthorized)

// obtenemos los query parameters
params := ctx.Request.URL.Query()
country := params.Get("country")
year := params.Get("year")

// path parameter
runnerId := ctx.Param("id")
```

hay que tener en cuenta que podemos encadenar varios métodos asociados a un recurso creando un pipeline. Cuando contestamos con un http status code/payload el pipeline continua con el siguiente elemento de la cadena. Sin embardo cuando contestamos con un _Abort_ el pipeline termina.

Como comentamos antes el controller se crea con una factoría, y tiene tantos métodos como recursos tiene la api.

```go
// el controlador tiene como propiedades los servicios que va a utilizar
type ResultsController struct {
	resultsService *services.ResultsService
	usersService   *services.UsersService
}

// factoria que crea el controlador
func NewResultsController(resultsService *services.ResultsService,
	userService *services.UsersService) *ResultsController {

	return &ResultsController{
		resultsService: resultsService,
		usersService:   userService,
	}
}

// los métodos de cada controler son los métodos que asociamos a los recursos de la api
func (rc ResultsController) CreateResult(ctx *gin.Context) {
[...]
}

func (rc ResultsController) DeleteResult(ctx *gin.Context) {
[...]
}
```

### Servicios

Implementa la lógica de negocio. Todos aquellos accesos que se precisen a la capa de datos se implementan en la capa Repositorio

### Repositorio

Implementa el acceso a datos. En la carpeta _Variantes_ tenemos ejemplos de implementación de la capa de datos con MySql, MongoDB y DynamoDB. Comentaré las pinceladas principales con la implementación de Postgres.

En primer lugar comentar como se gestionan las transacciones. Como cada repositorio gestiona el acceso a una tabla y hay lógica de negocio que trabaja con ambas tablas lo que haremos es a) crear una transacción, b) guardarla en los dos repositorios, de modo que cuando los métodos del repositorio accedan a los datos lo hagan usando la transacción. La transacción se convierte así en un elemento transversal para las dos tablas. Para gestionar la transacción usamos estos tres métodos:

```go
func BeginTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error {
	// creamos un contexto para la transacción
	ctx := context.Background()
	// iniciamos la transacción
	transaction, err := resultsRepository.dbHandler.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	// asignamos la transacción a ambos repositorios
	runnersRepository.transaction = transaction
	resultsRepository.transaction = transaction

	return nil
}

func RollbackTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error {
	// toma la transacción de uno de los repositorios (los dos tienen la misma transacción así que da igual cual utilicemos)
	transaction := runnersRepository.transaction
	// limpiamos la transacción en ambos repositorios
	runnersRepository.transaction = nil
	resultsRepository.transaction = nil
	// hacemos el rollback
	return transaction.Rollback()
}

func CommitTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error {
	// toma la transacción de uno de los repositorios (los dos tienen la misma transacción así que da igual cual utilicemos)
	transaction := runnersRepository.transaction

	// limpiamos la transacción en ambos repositorios
	runnersRepository.transaction = nil
	resultsRepository.transaction = nil

	// hacemos el commit
	return transaction.Commit()
}
```

Abrir un cursor para leer datos:

```go
query := `
SELECT id, race_result, location, position, year
FROM results
WHERE runner_id = $1`

// ejecutamos la query (consulta)
rows, err := rr.dbHandler.Query(query, runnerId)
if err != nil {
	return nil, &models.ResponseError{
		Message: err.Error(),
		Status:  http.StatusInternalServerError,
	}
}

// aseguramos que se cierre el cursor
defer rows.Close()

results := make([]*models.Result, 0)
var id, raceResult, location string
var position, year int

// iteramos sobre el cursor
for rows.Next() {
	// capturamos los datos recuperados con el cursor
	err := rows.Scan(&id, &raceResult, &location, &position, &year)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	result := &models.Result{
		ID:         id,
		RunnerID:   runnerId,
		RaceResult: raceResult,
		Location:   location,
		Position:   position,
		Year:       year,
	}

[...]
```

similar al caso anterior, pero utilizando una transacción:

```go
query := `
	INSERT INTO results(runner_id, race_result, location, position, year)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

// ejecutamos la query dentro de una transaccion (estamos cambiando datos)
rows, err := rr.transaction.Query(query, result.RunnerID, result.RaceResult, result.Location, result.Position, result.Year)
if err != nil {
	return nil, &models.ResponseError{
		Message: err.Error(),
		Status:  http.StatusInternalServerError,
	}
}

// aseguramos que se cierre el cursor
defer rows.Close()

var resultId string
// iteramos sobre el cursor
for rows.Next() {
	// capturamos los datos recuperados con el cursor
	err := rows.Scan(&resultId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
}
```

ejecutamos una query sin cursor (con o sin transacción; En este caso es con transacción):

```go
query := `
	UPDATE runners
	SET
		personal_best = $1,
		season_best = $2
	WHERE id = $3`

//ejecutamos la query
res, err := rr.transaction.Exec(query, runner.PersonalBest, runner.SeasonBest, runner.ID)
if err != nil {
	return &models.ResponseError{
		Message: err.Error(),
		Status:  http.StatusInternalServerError,
	}
}

// vemos cuantas filas fueron afectadas
rowsAffected, err := res.RowsAffected()
if err != nil {
	return &models.ResponseError{
		Message: err.Error(),
		Status:  http.StatusInternalServerError,
	}
}
```

## Modelo

El payload que intercambiamos en las apis se modelará como una estructura indicando vía anotaciones como se mapeará al json correspondiente. Por ejemplo, en este caso filtramos el campo _Status_ e incluimos el campo _Message_ con el nombre _message_:

```go
type ResponseError struct {
	Message string `json:"message"`
	Status  int    `json:"-"` // El status no se incluye en la respuesta JSON - se informará en la cabecera http status code
}
```

```go
type Runner struct {
	ID           string    `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Age          int       `json:"age,omitempty"`
	IsActive     bool      `json:"is_active"`
	Country      string    `json:"country"`
	PersonalBest string    `json:"personal_best,omitempty"` // se incluye el campo en el json solo si no es vacío
	SeasonBest   string    `json:"season_best,omitempty"`   // se incluye el campo en el json solo si no es vacío
	Results      []*Result `json:"results,omitempty"`       // se incluye el campo en el json solo si no es nulo o vacío
}
```

```go
```

```go
```

```go
```