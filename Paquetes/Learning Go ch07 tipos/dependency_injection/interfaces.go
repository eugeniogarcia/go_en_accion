package main

// Interfaces para representar un DataStore, Logger y Ligica de negocio
type DataStore interface {
	UserNameForID(userID string) (string, bool)
}

type Logger interface {
	Log(message string)
}

type Logic interface {
	SayHello(userID string) (string, error)
}
