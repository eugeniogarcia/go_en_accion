package server

import (
	"database/sql"
	"log"
	"runners-postgresql/controllers"
	"runners-postgresql/repositories"
	"runners-postgresql/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Servidor HTTP que maneja las solicitudes entrantes
type HttpServer struct {
	config            *viper.Viper
	router            *gin.Engine
	runnersController *controllers.RunnersController
	resultsController *controllers.ResultsController
	usersController   *controllers.UsersController
}

func InitHttpServer(config *viper.Viper, dbHandler *sql.DB) HttpServer {
	// Crea el repositorio
	runnersRepository := repositories.NewRunnersRepository(dbHandler)
	resultRepository := repositories.NewResultsRepository(dbHandler)
	usersRepository := repositories.NewUsersRepository(dbHandler)

	// Crea los servicios
	runnersService := services.NewRunnersService(runnersRepository, resultRepository)
	resultsService := services.NewResultsService(resultRepository, runnersRepository)
	usersService := services.NewUsersService(usersRepository)

	// Crea el controller
	runnersController := controllers.NewRunnersController(runnersService, usersService)
	resultsController := controllers.NewResultsController(resultsService, usersService)
	usersController := controllers.NewUsersController(usersService)

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

	// devuelve el servidor HTTP configurado
	return HttpServer{
		config:            config,
		router:            router,
		runnersController: runnersController,
		resultsController: resultsController,
		usersController:   usersController,
	}
}

// implementa el método Start para el apigateway HTTP (router de Gin)
func (hs HttpServer) Start() {
	// arrancar significa arrancar el router en la dirección indicada en la configuración. Si en la configuración solo especificamos el puerto (por ejemplo, ":8080"), el servidor escuchará en todas las interfaces de red disponibles.
	err := hs.router.Run(hs.config.GetString("http.server_address"))
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
