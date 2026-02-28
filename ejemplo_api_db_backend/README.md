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

- **Controller**. Esta capa proporciona las funciones que se asocian a los recursos de la api
- **Servicio**. Esta capa implementa la lógica de negocio. La función del controller hará uso de las funciones de negocio implementadas en el servicio. El acceso a datos se abstrae en la siguiente capa, el repositorio
- **Repository**. Implementa la lógica de acceso a datos. Se han considerado cuatro implementaciones alternativas:
	- Usando Postgres
	- Utilizando MySql
	- Usando MongoDB
	- utilizando DynamoDB

La construcción del proxy la hacemos utilizando el paquete Gin. Con Gin definimos para cada recurso/método el controlador aosociado:

```go
// instancia el router de Gin...
router := gin.Default()

// ...y define las rutas y los controladores asociados
router.POST("/runner", runnersController.CreateRunner)
router.PUT("/runner", runnersController.UpdateRunner)
router.DELETE("/runner/:id", runnersController.DeleteRunner)
router.GET("/runner/:id", runnersController.GetRunner)

[...]
```

una vez creado y configurado el router, para empezar a atender peticiones hacemos `router.Run([direccion de escucha])`. En cada ruta estamos indicando una función. Esas funciones estan definidas en la capa controller.

La forma en la que se implementa cada capa (controller, servicio, repositorio) es la misma:

- Se usa una **función factoría** para crear una instancia del objeto (controller, servicio o repositorio). `func NewXXXXX([parametros]) *XXXXX {...}`
- La capa se modela con un struct que contiene todas las propiedades necesarias, así como los métodos necesarios para gestionar la capa (en el ejemplo anterior, el struct sería `XXXXX`)

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

- **Abrir un cursor para leer datos**:

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

similar al caso anterior, pero **utilizando una transacción**:

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

- **ejecutamos una query sin cursor** (con o sin transacción; En este caso es con transacción):

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

**Hay que destacar que cuando usamos la conexión a bases de datos relacionales, los métodos son los mismos independientemente del driver que usemos (podemos ver en la variante MySql como en el repositorio el acceso se hace igual que en el caso de Postgres)**.


### Transacciones

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

en los repositorios se usará la transacción o directamente la conexión a la base de datos dependiendo de si queremos o no trabajar con transacciones:

```go
rr.transaction.Query(query, [argumentos])
```

```go
rr.dbHandler.Query(query, [argumentos])
```

donde se gestiona la transacción es en la capa superior a la de repositorio, es decir, en la capa de servicio. 

```go
// Inicia una trasacción
err = repositories.BeginTransaction(rs.runnersRepository, rs.resultsRepository)
if err != nil {
	return nil, &models.ResponseError{
		Message: "Failed to start transaction",
		Status:  http.StatusBadRequest,
	}
}

[...]

// Crear el resultado
response, responseErr := rs.resultsRepository.CreateResult(result)
// Si hay un error, hacemos rollback y retornamos el error
if responseErr != nil {
	repositories.RollbackTransaction(rs.runnersRepository, rs.resultsRepository)
	return nil, responseErr
}

[...]

// Si hemos llegado hasta aquí, todo ha ido bien y hacemos commit
repositories.CommitTransaction(rs.runnersRepository, rs.resultsRepository)
return response, nil
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

## Base de datos

Es interesante ver la definición del esquema de la base de datos. He incluido comentarios sobre las lineas del script. Como puntos destacables de este primer script:

- Extensiones. Incluimos una extensión para poder usar el lenguaje plsql de postgres, y para generar _uuid_
- Configuramos diferentes parametros en el esquema: timeouts, aspectos regionales
- Definimos primary keys y foreing keys. Si usamos la vista _ERD_ en postgres podemos ver el modelo entidad relación resultante
- Definimos varios índices
- Se incluyen diferentes constrains en campos, así como valores por defecto (`NOT NULL`, `DEFAULT`)

```sql
-- No definimos ningún timeout para la ejecución de las sentencias, bloqueos o transacciones inactivas
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;

-- usamos el conjunto de caracteres UTF8
SET client_encoding = 'UTF8';
-- definimos la zona horaria por defecto como UTC
SET timezone = 'UTC';
-- definimos el formato de los números para que use el punto como separador decimal
SET numeric_std = 'on';
-- definimos el comportamiento de las comillas simples en las cadenas de texto  
SET standard_conforming_strings = on;
-- nivel de mensajes mínimos a mostrar
SET client_min_messages = warning;
-- desactivamos la seguridad a nivel de fila
SET row_security = off;

-- extensión que permite usar el lenguaje PL/pgSQL en funciones y triggers
CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
-- extensión que se utiliza para generar UUIDs. Incluye el tipo uuid que usamos en las columnas id de las tablas runners y results
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA pg_catalog;

SET search_path = public, pg_catalog;
SET default_tablespace = '';

-- runners
CREATE TABLE runners (
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(), -- generamos un UUID basado en la dirección MAC del servidor y la fecha/hora actual
    first_name text NOT NULL,
    last_name text NOT NULL,
    age integer,
    is_active boolean DEFAULT TRUE,
    country text NOT NULL,
    personal_best interval,
    season_best interval,
    CONSTRAINT runners_pk PRIMARY KEY (id) -- definimos la clave principal de la tabla
);

CREATE INDEX runners_country
ON runners (country); -- creamos un índice en la columna country para optimizar las consultas que filtren por país

CREATE INDEX runners_season_best
ON runners (season_best);

-- results
CREATE TABLE results (
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(), -- generamos un UUID basado en la dirección MAC del servidor y la fecha/hora actual
    runner_id uuid NOT NULL,
    race_result interval NOT NULL,
    location text NOT NULL,
    position integer,
    year integer NOT NULL,
    CONSTRAINT results_pk PRIMARY KEY (id), -- definimos la clave principal de la tabla
    CONSTRAINT fk_results_runner_id FOREIGN KEY (runner_id) -- definimos una foreign key que referencia a la tabla runners. La columna runner_id de results referencia a la columna id de runners. Cuando se actualiza o elimina un registro en runners, no se realiza ninguna acción en results
        REFERENCES runners (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);
```

en este otro script destacar:
- Usamos la extensión que nos permite aplicar un salt y hashear. Si quisieramos crear usuarios desde el código go, tendríamos que usar una query del tipo `INSERT INTO users(username, user_password, user_role) VALUES ($1, crypt($2, gen_salt('bf')), $3)`

```sql
-- incluimos la extesión pgcrypto para hashear contraseñas
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA pg_catalog;

CREATE TABLE users (
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    username text NOT NULL UNIQUE,
    user_password text NOT NULL,
    user_role text NOT NULL,
    access_token text,
    CONSTRAINT users_pk PRIMARY KEY (id) -- definimos la clave principal de la tabla
);

CREATE INDEX user_access_token
ON users (access_token); -- creamos un índice en la columna access_token para optimizar las consultas que filtren por token de acceso

-- La función gen_salt('bf') genera una sal aleatoria para el algoritmo Blowfish. Produce una salida diferente cada vez que se llama. La salida tiene la siguiente forma: $2a$<cost>$<22 character salt>, esto es, un identificador del alfgortimo que se ha usado ($2a$ corresponde al algoritmo blowfish), el coste de computación (cost) y el salt propiamente dicho (que tendra 22 caracteres de largo).
-- La función crypt() toma dos argumentos: el valor a hashear  y la sal (salt). Del salt toma el algoritmo y la salt propiamente dicha para hashear el valor. El valor y la salt se combinan (es más complejo que una concatenación de ambos) y se aplica el algoritmo de hash especificado en la salt para producir el hash resultante. El resultado es una cadena que incluye el identificador del algoritmo, el coste, la salt y el hash resultante
-- guardamos la contraseñas de los dos usuarios que hemos creado hasheadas con un salt y utilizando el algoritmo Blowfish
INSERT INTO users(username, user_password, user_role)
VALUES
    ('admin', crypt('admin', gen_salt('bf')), 'admin'),
    ('runner', crypt('runner', gen_salt('bf')), 'runner');
```

La función `crypt (arg1, arg2)` se utiliza para generar el hash, o más exactamente para generar algo que en _bcrypt_ tiene la siguiente pinta: `$2[aby]$CC$<22c-salt><31c-hash>`, o lo que es lo mismo, `$algoritmo$coste$salt$hash$`. El salt son exactamente 22 caracteres, y el hash 31.

Lo que usa `crypt` del segundo argumento es el prefijo de algoritmo (`$2b$`), el coste y exactamente los 22 caracteres de salt. Lo que viene a continuación lo ignora

```
password=crypt('contraseña', gen_salt('bf'))
password2=crypt('contraseña', password)

entonces password2 y password son iguales
```

por este motivo en el repositorio para validar las credenciales de un usuario hacemos:

```go
func (ur UsersRepository) LoginUser(username string, password string) (string, *models.ResponseError) {
	query := `
		SELECT id
		FROM users
		WHERE username = $1 and user_password = crypt($2, user_password)`

	rows, err := ur.dbHandler.Query(query, username, password)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
```

y para insertar haríamos `INSERT INTO users(username, user_password, user_role) VALUES ($1, crypt($2, gen_salt('bf')), $3)`

## Seguridad

Vamos a dotar la aplicación de capacidades de autenticación y autorización. Para implementar la autenticación creamos una tabla users en la que guardaremos el user y password - asi como el rol asociado. En el apartado de Base de Datos se ha comentado ya, la contraseña se guarda en la tabla usando su hash salteado. Usamos el propio motor de base de datos para calcular el salt y el hash, así que en el código go simplemente trabajaremos con la contraseña en claro.

Complementamos la tabla con un controller (que proporciona las funciones que asociaremos al recurso `login` y al recurso `logout`), y un servicio que implementa la lógica correspondiente al login y logout. El servicio login toma las credenciales de la cabecera de autenticación básica, comprueba que las credenciales sean válidas - que coinciden con el usuario y contraseña que tenemos guardada en la base de datos), y en caso afirmativo crea un token de acceso (usando bcrypt) que guarda en la base de datos y se devuelve como resultado del login. El logout lo que hace es tomar un token de acceso de la cabecera, lo busca en la base de datos y lo borra.

Con esto tenemos implementado el mecanismo de autenticación que valida las credenciales y en caso de ser correctas genera un token de acceso.

El mecanismo de autorización lo implementamos en los controladores. El controlador tomaá el token de acceso de la cabecera, lo busca en la base de datos, y recupera los roles que ese token tenga asociados. Una vez recuperados se verifica si el rol asociado es "suficiente" para ejecutar el controlador. 

```go
func (rc ResultsController) CreateResult(ctx *gin.Context) {
	// recuperamos una cabecera de la petición
	accessToken := ctx.Request.Header.Get("Token")

	// verificamos que el token tenga asociado el role ROLE_ADMIN
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	if responseErr != nil {
		// contruye una respuesta con el http status code y el payload
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	if !auth {
		// responde con el http status code
		ctx.Status(http.StatusUnauthorized)
		return
	}

	[...]
```

## Observavilidad

Para observar la aplicación vamos a utilizar Prometheus y Grafana. Con Prometheus podemos definir métricas, e instrumentalizar los servicios/aplicaciones para que publique estas métricas en el repositorio central de Prometheus. La extracción de las metricas puede hacerse en modo pull (Prometheus extrae las metricas) o push (las aplicaciones/servicios publican las métricas). Típicamente se hace pull para asegurar que el repositorio central de Prometheus no se sature.

### Métricas

Las métricas en Prometheus se clasifican en tres tipos: contadores, histogramas y gauges

- Contadores se utilizan para contar eventos, como el número de peticiones HTTP o el número de errores, lo que nos permite monitorear el comportamiento de la aplicación y detectar posibles problemas
- Gauges se utilizan para medir valores que pueden subir y bajar, como la cantidad de memoria utilizada o el número de conexiones abiertas, pero en este caso no los vamos a utilizar
- Histogramas se utilizan para medir la distribución de los tiempos de ejecución de una operación, lo que nos permite entender mejor el rendimiento de la aplicación y detectar posibles cuellos de botella

otro concepto que se usa con las métricas es el de etiqueta. Cuando se publica una métrica se puede etiquetar con una o varias etiquetas (pares clave-valor). Cada etiquete representa una dimension en la que se pueden analizar las métricas.

En nuestro caso vamos a definir tres métricas, un contador (las gauges se manejan igual que los contadores) y un histograma. En una de las métricas vamos a utilizar una etiqueta. Las definimos en `metrics.go`:

```go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Definimos las metricas que vamos a utilizar en la aplicacion, estas metricas se pueden usar en cualquier parte de la aplicacion para medir el rendimiento y el comportamiento de las diferentes operaciones.
var (
	HttpRequestsCounter = promauto.NewCounter( // un contador
		prometheus.CounterOpts{
			Name: "runners_app_http_requests", // nombre de la metrica
			Help: "Número total de peticiones HTTP", // descripcion de la metrica
		},
	)

	GetRunnerHttpResponsesCounter = promauto.NewCounterVec( // un contador con etiquetas. Las etiquetas nos permiten clasificar las métricas en diferentes dimensiones. En este caso, vamos a clasificar las respuestas HTTP del endpoint get runner por su código de estado (status).
		prometheus.CounterOpts{
			Name: "runners_app_get_runner_http_responses",
			Help: "Número total de respuestas HTTP para el endpoint get runner",
		},
		[]string{"estado"}, // creamos una etiqueta llamada estado para clasificar las respuestas HTTP por su código de estado
	)

	GetAllRunnersTimer = promauto.NewHistogram( // un histograma. Un histograma nos permite medir la distribución de los tiempos de ejecución de una operación. En este caso, vamos a medir la duración de la operación get all runners.
		prometheus.HistogramOpts{
			Name: "runners_app_get_all_runners_duration",
			Help: "Duración de la operación get all runners en segundos",
		},
	)
)
```

### Exportar las métricas

Una vez tenemos las metricas definidas, tenemos que exportarlas. Para exportarlas se usa un cliente de Prometheus. En `prometheus.go` tenemos definido el exporter:

```go
func InitPrometheus() {
	http.Handle("/metrics", promhttp.Handler()) // endpoint en el que expondremos las métricas
	http.ListenAndServe(":9000", nil)
}
```

y lo arrancamos en `main.go` en una gorutina:

```go
// inicializamos Prometheus
go server.InitPrometheus()
```

### Informar las métricas

Una vez hemos definido las métricas y configurado el exportador de métricas, tenemos que identificar en la lógica de la aplicación donde debemos darles un valor.

En el controlador, cada vez que recibimos una llamada actualizamos el contador:

```go
func (rc RunnersController) GetRunnersBatch(ctx *gin.Context) {
	// actualizamos la metrica de contador de peticiones HTTP cada vez que se recibe una solicitud en el endpoint create runner	
	metrics.HttpRequestsCounter.Inc()

[...]
```

en las metricas en las que usamos etiquetas, además de informar el valor de la métrica tenemos que informar las etiquetas:

```go
[...]

if responseErr != nil {
	// actualizamos la metrica de contador informando también la etiqueta correspondiente al status code
	metrics.GetRunnerHttpResponsesCounter.WithLabelValues(
		strconv.Itoa(responseErr.Status)).Inc()
	ctx.JSON(responseErr.Status, responseErr)
	return
}

// actualizamos la metrica de contador informando también la etiqueta correspondiente al status code
metrics.GetRunnerHttpResponsesCounter.WithLabelValues("200").Inc()

[...]
```

en este ejemplo estamos informando la etiqueta con el http status code.

Podemos usar un timer para informar una métrica de tipo histograma y medir la duración. El timer registra automáticamente el tiempo que se ha tardado en procesar la peticion y permite el cálculo del p50, p95, p99 a lo largo del tiempo.

```go
[...]

// Medimos la duración de la operación (percentiles, valor medio, desviacion estándar, etc.) utilizando un histograma de Prometheus. Para ello, creamos un timer al inicio del handler y lo detenemos al final del handler utilizando defer. El timer observará la duración de la operación y actualizará el histograma con ese valor.
timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
	metrics.GetAllRunnersTimer.Observe(f)
}))

defer func() {
	//termina la observación, para el cronómetro y actualiza el histograma con la duración de la operación
	timer.ObserveDuration()
}()

[...]
```
 
### Grafana

Se describe en el apartado de [Kubernetes](#kubernetes).

## Automatización de tests

Se utiliza el paquete de tests para automatizar las pruebas. No hay nada especial salvo el uso de mocks para emular la base de datos - durante la ejecución de los tests.

### comandos

Para ejecutar los test de un determinado directorio. Aquí se estarían ejecutando todos los tests definidos en el directorio actual (_vervose_):

```ps
go test -v
```

si queremos ejecutar todos los tests incluidos en sub-directorios del raiz:

```ps
go test -v ./…
```

ejecuta todos los test incluidos en el directorio `services`:

```ps
go test -v ./services
```

si queremos analizar también la cobertura de código de los tests ejecutados tenemos dos flags, `-cover` y `coverpkg`. El primer flag verifica la cobertura local (esto es la cobertura que suponen los tests sobre el paquete que se está probando), el segundo es global (la cobertura del código independientemente del paquete que se esté probando). Si desde el directorio raiz hacemos:

```ps
go test -v -cover ./services
```

se ejecutan los tests definidos en el directorio `services` y la cobertura sobre el código definido en `services`. Si hacemos

```ps
go test -v -coverpkg=./... ./services
```

se ejecutan los tests definidos en el directorio `services` y la cobertura sobre todos los paquetes (si por ejemplo se importa en el código otros paquetes, se incluyen en el % de cobertura). Con `coverpkg` es necesario indicar el directorio/patrón que sobre el que se aplicará la cobertura, independientemente del directorio/patrón sobre el que se apliquen los tests. Con `-coverprofile=[nombre fichero]` hacemos que el resultado de los tests no salga por pantalla sino que se diija a un fichero

```ps
go test -v -coverpkg=./… ./… -coverprofile=coverage.out
```

podemos visualizar la cobertura con 

```ps
go tool cover -html=coverage.out
```

### Mocks de base de datos

Usamos el paquete `github.com/DATA-DOG/go-sqlmock` para definir un mock para la base de datos. El método `dbHandler, mock, error := sqlmock.New()` crea una conexión mock a la base de datos, nos devuelve un objeto con el que definir el mock, y el error.

- Creamos la conexión _mockeada_ y un _mock_:

```go
// usa sqlmock para crear una conexión a la base de datos, y el objeto mock
dbHandler, mock, _ := sqlmock.New()
// aseguramos que al final se cierre la conexión a la base de datos
defer dbHandler.Close()
```

- Con el mock definimos los casos a simular. Para cada caso tenemos que indicar la query (o parte de la query. Por ejemplo en este caso un select para recuperar el `user_role`; Podemos ver que la query no está escrita completamente, apenas tener un extracto de ella). Se indican las columnas que se deben recuperar, y los datos:

```go
// usamos mock para definir un mock. Indicamos las columnas que queremos que nos devuelva el mock...
columnsUsers := []string{"user_role"}
//...indicamos que query tiene que se mockeada, que columnas se tienen que devolver, y los valores - una sola final
mock.ExpectQuery("SELECT user_role").WillReturnRows(
	sqlmock.NewRows(columnsUsers).AddRow("runner"),
)
```

otro ejemplo:

```go
// definimos otro mock para un select *; Indicamos las columnas que tiene que devolver el mock, y los valores - dos filas
columns := []string{"id", "first_name", "last_name", "age", "is_active", "country", "personal_best", "season_best"}
mock.ExpectQuery("SELECT *").WillReturnRows(
	sqlmock.NewRows(columns).
		AddRow("1", "John", "Smith", 30, true, "United States", "02:00:41", "02:13:13").
		AddRow("2", "Marijana", "Komatinovic", 30, true, "Serbia", "01:18:28", "01:18:28"))
```

- Con esto, si ahora creamos el controlador usando la conexión a base de datos mockeada, cuando llegue el controller llame a alguna función de la capa de servicios, y está llame a alguna función de la capa repositorio, cuando la función del repositorio haga un `select * from runners" no se consultará la base de datos real sino la mockeada, se buscará entre los casos que hemos definido, y recuperar los datos definidos en el caso.

```go
// definimos el router, usando la conexión a la base de datos mockeada
router := initTestRouter(dbHandler)
```

En este punto combiene enteder que es un ruter y como se utilizan los http handlers en el contexto de un servidor http. Es combeniente repasar el [ejemplo de middleware](../Learning%20Go/ch13%20libreria%20estandard/middleware/main.go).

el `router` de Gin es un http handler que como tal puede usarse para gestionar peticiones http en un servidor http. El handler http incluye un método `ServeHTTP([http writer], [http request])`. Típicamente no vemos la llamada a este método, cuando hacemos `router.Run` con un router Gin, o `ListenAndServe()` con un servidor http clásico, se invoca bajo bambalinas a este método. El primer argumento es el mecanismo por el que se prepara la respuesta y se responde vía http. Si estamos en modo test podemos usar un http especial, `httptest.NewRecorder()`: 

```go
// El recorder implementa http.ResponseWriter y nos permite capturar lo que el handler escribe en la response
recorder := httptest.NewRecorder()

// usamos el http handler de router para probar. Como http writer usamos un recorder, y como request la request que hemos creado, Pasamos el request al handler
router.ServeHTTP(recorder, request)
```

con esto estamos simulando la recepción y gestion de una petición http, y en `recorder` podremos ver la respuesta que el router ha proporcionado:

```go
// comprueba el status code de la respuesta
assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

// comprueba el body de la respuesta
var runers []*models.Runner
json.Unmarshal(recorder.Body.Bytes(), &runers)

assert.NotEmpty(t, runers)
assert.Equal(t, 2, len(runers))
```

## docker

informacion de docker y de la máquina donde corre docker:

```ps
docker version

docker info
```

información de los contenedores que están corriendo y de las imagenes disponibles:

```
docker ps

docker images
```

descargar una imagen desde un repositorio, borrar una imagen, arrancar y parar un contenedor, y borrar un contenedor:

```ps
docker_pull image_name:tag

docker rmi image_name:tag

docker start container_name

docker stop container_name

docker rm container_name
```

construir una imagen:

```
docker build -f name_of_dockerfile path_to_dockerfile
-t image_name:tag
```

para construir nuestra imagen hacemos:

```ps
docker build -f .\Dockerfile .\ -t egsmartin/runners-app:latest
```

nuestro dockerfile:

```dockerfile
# imagen base, con Go y Alpine Linux
# Start from golang alpine base image
FROM golang:1.25-alpine

# Creamos el directorio de trabajo en la imagen
WORKDIR /app

# copiamos todo al directorio de trabajo
COPY . .

# descargamos las dependencias
RUN go mod download

# construimos la aplicación
RUN go build -o runners-app main.go

# exponemos el puerto 8080
EXPOSE 8080

# Comando para ejecutar cuando se inicie el contenedor
CMD ["/app/runners-app"]
```

la arrancamos:

```ps
docker run -p 8080:8080 runners-app
```

La aplicación se ejecuta dentro de un contenedor Docker mientras que la base de datos está en la máquina local. En Linux podemos hacer que el contenedor comparta el mismo espacio de direcciones que el host utilizando el `host networking mode`:

```ps
docker run --network host -p 8080:8080 runners-app
```

esto solo funciona en Linux. En windows lo que se crea es una entrada en la resolución de nombres del contenedor que apunta al host: `host.docker.internal`. De este modo cuando queremos conectarnos con el postgress que tenemos instalado en el host tendremos que referirnos a `host.docker.internal` en lugar de a `localhost`. Podemos incluir la variable de entorno en el _dockerfile_ haciendo incluyendo una línea `ENV ENV=k8s`, o podemos pasar la variable de entorno al arrancar:

```ps
docker run -p 8080:8080 -e ENV=k8s runners-app
```

podemos arrancar el contenedor _detachado_ (notese que hemos indicado el nombre de la imagen):

```ps
docker run -d --name mi-app -p 8080:8080 -e ENV=k8s runners-app
```

podemos ver los logs:

```ps
docker logs mi-app
```

### Optimizacion de la imagen

Podemos mejorar la imagen que hemos creado, tanto en el espacio y recursos que emplean como en el perfil de ataque que expone. Vamos a usar un pipeline en el que utilizamos una imagen para construir la salida que necesitamos, y una segunda imagen, la final, en la que copiaremos el resultado

```dockerfile
########################################
# Etapa 1: Build (builder)
########################################
# Imagen base con Go (Alpine, musl). Rápida y pequeña
FROM golang:1.25-alpine AS imagen_constructora

# Herramientas necesarias (git para módulos privados si aplica)
RUN apk add --no-cache git ca-certificates

# creamos el directorio de trabajo
WORKDIR /src

# no copiamos toda la aplicación todavía. Copiamos aquellas partes menos susceptibles a cambios, de modo que es más improbable que esta capa tenga que reconstruirse a menudo. Esto optimiza el uso de la caché de Docker y mejora los tiempos de construcción.
# Optimiza caché de módulos copiando primero los manifests
COPY go.mod go.sum ./

# Descarga dependencias
RUN go mod download

# Ahora ya copiamos el resto del código
COPY . .

# Variables para cross-compilación y binario estático
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH

# Compila el binario, reduciendo tamaño
RUN go build -trimpath -ldflags "-s -w" -o /out/runners-app ./main.go


########################################
# Etapa 2: Runtime (mínimo, no-root)
########################################
# Distroless estático, seguro y muy pequeño (sin shell)
FROM gcr.io/distroless/static:nonroot

# Creamos directorio de aplicación
WORKDIR /app

# Copiamos binario
COPY --from=imagen_constructora /out/runners-app /app/runners-app

# Copiamos archivos de configuración necesarios
COPY --from=imagen_constructora /src/runners.toml /app/runners.toml
COPY --from=imagen_constructora /src/runners-k8s.toml /app/runners-k8s.toml

EXPOSE 8080

# Variable de entorno para seleccionar el archivo de configuración
ENV ENV=k8s

# Ejecuta como usuario no root
USER nonroot

# Entrada del contenedor
ENTRYPOINT ["/app/runners-app"]
```

construimos la imagen usando otro nombre para poder así comparar una con otra:

```ps
docker build -f .\Dockerfile .\ -t mi-runners-app
```

vemos que la nueva imagen ocupa 6.77MB en lugar de 167MB, y el espacio en disco 27.3MB frente a 806MB:

```ps
docker images

IMAGE                   ID             DISK USAGE   CONTENT SIZE   EXTRA
mi-runners-app:latest   985bb5a18554       27.3MB         6.77MB
runners-app:latest      30aa4dce5595        806MB          167MB
```

La imagen nueva que hemos creado es distrolless, no tiene siquiere un shell, y se ejecuta con un `USER nonroot`. La imagen que se usa para compilar el programa podemos eliminarla:

```ps
docker system prune
```

### Publicar una imagen

Para publicar la imagen primero la tageamos, indicando el prefijo del repositorio al que querremos publicarla. Voy a publicarla al dockerhub:

```ps
docker tag runners-app egsmartin/runners-app:v1.0
```

vemos la imagen tageada (tenemos las dos imagenes en el repositorio local, aunque no ocupan espacio doble):

```ps
docker images

IMAGE                        ID             DISK USAGE   CONTENT SIZE   EXTRA
egsmartin/runners-app:v1.0   f914b331558d       27.3MB         6.77MB
runners-app:latest           f914b331558d       27.3MB         6.77MB
```

ahora hacemos login y publicamos:

```ps
docker login

docker push egsmartin/runners-app:v1.0
```

### docker compose

```ps
docker-compose up
docker-compose down
```

o detachadas:

```ps
docker-compose -d up
docker-compose -d down
```

podemos añadir varias instanacias:

docker compose up -d --scale runners-app=3

## Kubernetes

He estudiado [kuernetes](https://github.com/eugeniogarcia/kubernetes_up_running/blob/main/readme.md) en otro repositorio.

Para probar esta aplicacion voy a utilizar el cluster _kind_ que se incluye con Docker Desktop. Una vez creado el cluster utilizo el script [`.\instalar.ps1`](https://github.com/eugeniogarcia/kubernetes_up_running/blob/main/instalar.ps1) que se incluye en mi repositorio _kubernetes up and running_. Este script instala el ingress de _contour/envoy_, el metrics server, y el operador _k6_ (para hacer pruebas de rendimiento).

Crearemos dentro del cluster un [cluster de Postgress](https://github.com/eugeniogarcia/kubernetes_up_running/blob/main/postgress/readme.md), y también las instancias de [grafana y prometheus](https://github.com/eugeniogarcia/kubernetes_up_running/blob/main/prometheus/readme.md).

Tenemos que configurar el [archivo de configuración](runners-k8s.toml) para que apunte a la instancia de Postgres que hemos creado. El cluster de Postgres expone tres servicios:

- `mi-postgres-rw.database.svc`. Este servicio ataca a la instancia principal, y admite lecturas y escrituras
- `mi-postgres-ro.database.svc`. Este servicio ataca las réplicas y solo sirve para hacer consultas. El servicio balance las peticiones entre todas las réplicas
- `mi-postgres-r.database.svc`. Este servicio permite acceder e solo lectura a una réplica concreta. Para llamar hay que hacer `[nombre replica].mi-postgres-r.database.svc`

En nuestro caso usaremos el servicio de lectura y escritura.

En al apartado anterior creamos la imagen de la aplicación con:

```ps
docker build -f .\Dockerfile .\ -t egsmartin/runners-app:latest
```

la publicamos en el docker hub:

```ps
docker push egsmartin/runners-app:latest
```

y a continuación ya podemos desplegar la aplicacióne en el cluster kubernetes

```ps
kubectl apply -f .\runners-app-deployment.yaml

kubernetes> kubectl apply -f .\runners-app-service.yaml
```

En este punto podemos comprobar que todo funciona. El servicio tipo cluster que hemos creado se expone el `localhost:8080`. En primer lugar creamos un token usando basic auth. El usuario `admin` contraseña `admin` tiene el role necesario para crear runners:

```ps
curl --location --request POST 'localhost:8080/login' \
--header 'Authorization: ••••••' \
--data ''
```

Usamos le token para llamar al resto de servicios, pasandolo en una cabecera llamada `Token`. Creamos runners:

```ps
curl --location 'localhost:8090/runner' \
--header 'Token: JDJhJDEwJG5ONmFXQUdDMXprSWJmZWFYV1RxSy5KcnRMVWZ4dU43Q25pR1lmQmxCZlF0T2F3MFN5cEJH' \
--header 'Content-Type: application/json' \
--data '{
"first_name":"Nicolas",
"last_name":"Garcia Zach",
"age":14,
"is_active":true,
"country":"España",
"personal_best":"",
"season_best":""
}
'
```

consultamos los runners creados:

```ps
curl --location 'localhost:8090/runner' \
--header 'Token: JDJhJDEwJC94UWdCMGUzb3pVLkdWRUZsMGlDeXU1cmM2UG9uYUNTZGxRZWouVkd1ZzdXeERlTHcyc1BP'
```

en el script [simulamos carga para viaualizar en grafana](./kubernetes/load_generator.ps1).
