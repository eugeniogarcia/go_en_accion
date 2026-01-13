package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runners-postgresql/models"
	"runners-postgresql/repositories"
	"runners-postgresql/services"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetRunnersResponse(t *testing.T) {
	// usa sqlmock para crear una conexión a la base de datos, y el objeto mock
	dbHandler, mock, _ := sqlmock.New()
	// aseguramos que al final se cierre la conexión a la base de datos
	defer dbHandler.Close()

	// usamos mock para definir un mock. Indicamos las columnas que queremos que nos devuelva el mock...
	columnsUsers := []string{"user_role"}
	//...indicamos que query tiene que se mockeada, que columnas se tienen que devolver, y los valores - una sola final
	mock.ExpectQuery("SELECT user_role").WillReturnRows(
		sqlmock.NewRows(columnsUsers).AddRow("runner"),
	)

	// definimos otro mock para un select *; Indicamos las columnas que tiene que devolver el mock, y los valores - dos filas
	columns := []string{"id", "first_name", "last_name", "age", "is_active", "country", "personal_best", "season_best"}
	mock.ExpectQuery("SELECT *").WillReturnRows(
		sqlmock.NewRows(columns).
			AddRow("1", "John", "Smith", 30, true, "United States", "02:00:41", "02:13:13").
			AddRow("2", "Marijana", "Komatinovic", 30, true, "Serbia", "01:18:28", "01:18:28"))

	// definimos el router, usando la conexión a la base de datos mockeada
	router := initTestRouter(dbHandler)

	// crea una request (GET, al recurso /runner, con un payload nulo)
	request, _ := http.NewRequest("GET", "/runner", nil)
	// añade el header 'token' a la request
	request.Header.Set("token", "token")

	// El recorder implementa http.ResponseWriter y nos permite capturar lo que el handler escribe en la response
	recorder := httptest.NewRecorder()

	// usamos el http handler de router para probar. Como http writer usamos un recorder, y como request la request que hemos creado, Pasamos el request al handler
	router.ServeHTTP(recorder, request)

	// comprueba el status code de la respuesta
	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

	// comprueba el body de la respuesta
	var runers []*models.Runner
	json.Unmarshal(recorder.Body.Bytes(), &runers)

	assert.NotEmpty(t, runers)
	assert.Equal(t, 2, len(runers))
}

func initTestRouter(dbHandler *sql.DB) *gin.Engine {
	// apenas definimos las capas que queremos usar en el test. Estamos usando la base de datos mockeada
	runnersRepository := repositories.NewRunnersRepository(dbHandler)
	usersRepository := repositories.NewUsersRepository(dbHandler)
	// no usamos el repositorio de tokens en este test, por eso le pasamos nil
	runnersService := services.NewRunnersService(runnersRepository, nil)
	usersServices := services.NewUsersService(usersRepository)
	runnersController := NewRunnersController(runnersService, usersServices)

	router := gin.Default()
	// solo incluimos la ruta que queremos testear
	router.GET("/runner", runnersController.GetRunnersBatch)

	return router
}
