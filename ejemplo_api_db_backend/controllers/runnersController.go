package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runners-postgresql/metrics"
	"runners-postgresql/models"
	"runners-postgresql/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

const ROLE_ADMIN = "admin"
const ROLE_RUNNER = "runner"

type RunnersController struct {
	runnersService *services.RunnersService
	usersService   *services.UsersService
}

func NewRunnersController(runnersService *services.RunnersService, usersService *services.UsersService) *RunnersController {
	return &RunnersController{
		runnersService: runnersService,
		usersService:   usersService,
	}
}

func (rc RunnersController) CreateRunner(ctx *gin.Context) {
	// actualizamos la metrica de contador de peticiones HTTP cada vez que se recibe una solicitud en el endpoint create runner
	metrics.HttpRequestsCounter.Inc()

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

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("Error while unmarshaling create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := rc.runnersService.CreateRunner(&runner)
	if responseErr != nil {
		// responde con el http status code y el payload, y detiene la ejecución del handler
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (rc RunnersController) UpdateRunner(ctx *gin.Context) {
	// actualizamos la metrica de contador de peticiones HTTP cada vez que se recibe una solicitud en el endpoint create runner
	metrics.HttpRequestsCounter.Inc()

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

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading update runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("Error while unmarshaling update runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	responseErr = rc.runnersService.UpdateRunner(&runner)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (rc RunnersController) DeleteRunner(ctx *gin.Context) {
	// actualizamos la metrica de contador de peticiones HTTP cada vez que se recibe una solicitud en el endpoint create runner
	metrics.HttpRequestsCounter.Inc()

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

	runnerId := ctx.Param("id")

	responseErr = rc.runnersService.DeleteRunner(runnerId)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (rc RunnersController) GetRunner(ctx *gin.Context) {
	// actualizamos la metrica de contador de peticiones HTTP cada vez que se recibe una solicitud en el endpoint create runner
	metrics.HttpRequestsCounter.Inc()

	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN, ROLE_RUNNER})
	if responseErr != nil {
		// actualizamos la metrica de contador informando también la etiqueta correspondiente al status code
		metrics.GetRunnerHttpResponsesCounter.WithLabelValues(
			strconv.Itoa(responseErr.Status)).Inc()
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	if !auth {
		// actualizamos la metrica de contador informando también la etiqueta correspondiente al status code
		metrics.GetRunnerHttpResponsesCounter.WithLabelValues("401").Inc()
		ctx.Status(http.StatusUnauthorized)
		return
	}

	// path parameter
	runnerId := ctx.Param("id")

	response, responseErr := rc.runnersService.GetRunner(runnerId)
	if responseErr != nil {
		// actualizamos la metrica de contador informando también la etiqueta correspondiente al status code
		metrics.GetRunnerHttpResponsesCounter.WithLabelValues(
			strconv.Itoa(responseErr.Status)).Inc()
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	// actualizamos la metrica de contador informando también la etiqueta correspondiente al status code
	metrics.GetRunnerHttpResponsesCounter.WithLabelValues("200").Inc()
	ctx.JSON(http.StatusOK, response)
}

func (rc RunnersController) GetRunnersBatch(ctx *gin.Context) {
	// actualizamos la metrica de contador de peticiones HTTP cada vez que se recibe una solicitud en el endpoint create runner
	metrics.HttpRequestsCounter.Inc()

	// Medimos la duración de la operación (percentiles, valor medio, desviacion estándar, etc.) utilizando un histograma de Prometheus. Para ello, creamos un timer al inicio del handler y lo detenemos al final del handler utilizando defer. El timer observará la duración de la operación y actualizará el histograma con ese valor.
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		metrics.GetAllRunnersTimer.Observe(f)
	}))

	defer func() {
		//termina la observación, para el cronómetro y actualiza el histograma con la duración de la operación
		timer.ObserveDuration()
	}()

	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN, ROLE_RUNNER})
	fmt.Println("Response error", responseErr)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	// obtenemos los query parameters
	params := ctx.Request.URL.Query()
	country := params.Get("country")
	year := params.Get("year")

	response, responseErr := rc.runnersService.GetRunnersBatch(country, year)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
