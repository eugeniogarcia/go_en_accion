package controllers

import (
	"log"
	"net/http"
	"runners-postgresql/services"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	usersService *services.UsersService
}

func NewUsersController(usersService *services.UsersService) *UsersController {
	return &UsersController{
		usersService: usersService,
	}
}

func (uc UsersController) Login(ctx *gin.Context) {
	// Obtiene las credenciales de autenticación básica
	username, password, ok := ctx.Request.BasicAuth()
	if !ok {
		log.Println("Error while reading credentials")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Valida el usuario y contraseña contra lo que tenemos guardado en la base de datos, y si son correctos genera un token de acceso (que se guarda en la base de datos) y se obtiene aqui
	accessToken, responseErr := uc.usersService.Login(username, password)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	// Devuelve el token de acceso al cliente
	ctx.JSON(http.StatusOK, accessToken)
}

func (uc UsersController) Logout(ctx *gin.Context) {
	// Obtiene el token de acceso de la cabecera Token
	accessToken := ctx.Request.Header.Get("Token")

	// Llama al servicio que elimina el token de acceso de la base de datos
	responseErr := uc.usersService.Logout(accessToken)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}
