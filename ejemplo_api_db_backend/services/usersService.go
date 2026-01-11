package services

import (
	"encoding/base64"
	"net/http"
	"runners-postgresql/models"
	"runners-postgresql/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	usersRepository *repositories.UsersRepository
}

func NewUsersService(usersRepository *repositories.UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

func (us UsersService) Login(username string, password string) (string, *models.ResponseError) {
	// Validaciones
	if username == "" || password == "" {
		return "", &models.ResponseError{
			Message: "Invalid username or password",
			Status:  http.StatusBadRequest,
		}
	}

	// Comprueba si el usuario y contraseña los tenemos en la base de datos, y si los tenemos obtenemos su id
	id, responseErr := us.usersRepository.LoginUser(username, password)
	if responseErr != nil {
		return "", responseErr
	}

	if id == "" {
		return "", &models.ResponseError{
			Message: "Login failed",
			Status:  http.StatusUnauthorized,
		}
	}

	// Crea un token de acceso para el usuario
	accessToken, responseErr := generateAccessToken(username)
	if responseErr != nil {
		return "", responseErr
	}
	// Guarda el token de acceso en la base de datos asociado al usuario
	us.usersRepository.SetAccessToken(accessToken, id)

	return accessToken, nil
}

func (us UsersService) Logout(accessToken string) *models.ResponseError {
	if accessToken == "" {
		return &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}
	// Elimina el token de acceso de la base de datos
	return us.usersRepository.RemoveAccessToken(accessToken)
}

func (us UsersService) AuthorizeUser(accessToken string, expectedRoles []string) (bool, *models.ResponseError) {
	if accessToken == "" {
		return false, &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}

	role, responseErr := us.usersRepository.GetUserRole(accessToken)
	if responseErr != nil {
		return false, responseErr
	}

	if role == "" {
		return false, &models.ResponseError{
			Message: "Failed to authorize user",
			Status:  http.StatusUnauthorized,
		}
	}

	for _, expectedRole := range expectedRoles {
		if expectedRole == role {
			return true, nil
		}
	}

	return false, nil
}

func generateAccessToken(username string) (string, *models.ResponseError) {
	// Creamos un token a partir del nombre de usuario. En la generación del token se utiliza el timestamp
	hash, err := bcrypt.GenerateFromPassword([]byte(username), bcrypt.DefaultCost)
	if err != nil {
		return "", &models.ResponseError{
			Message: "Failed to generate token",
			Status:  http.StatusInternalServerError,
		}
	}
	// codifica el token en base64 para que sea seguro para su transmisión
	return base64.StdEncoding.EncodeToString(hash), nil
}
