package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runners-postgresql/models"
	"runners-postgresql/services"

	"github.com/gin-gonic/gin"
)

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
	// recuperamos una cabecera de la petición
	accessToken := ctx.Request.Header.Get("Token")

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

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create result request body", err)
		// responde con el http status code y un payload, y detiene la ejecución del handler
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var result models.Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Error while unmarshaling create result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := rc.resultsService.CreateResult(&result)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (rc ResultsController) DeleteResult(ctx *gin.Context) {
	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	resultId := ctx.Param("id")

	responseErr = rc.resultsService.DeleteResult(resultId)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}
