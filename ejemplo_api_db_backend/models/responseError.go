package models

type ResponseError struct {
	Message string `json:"message"`
	Status  int    `json:"-"` // El status no se incluye en la respuesta JSON - se informará en la cabecera http status code
}
