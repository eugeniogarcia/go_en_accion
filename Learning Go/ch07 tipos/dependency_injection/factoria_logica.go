package main

import "errors"

//Acepta interfaces y devuelve un SimpleLogic
func NewSimpleLogic(l Logger, ds DataStore) SimpleLogic {
	return SimpleLogic{
		l:  l,
		ds: ds,
	}
}

//define SimpleLogic
type SimpleLogic struct {
	l  Logger
	ds DataStore
}

func (sl SimpleLogic) SayHello(userID string) (string, error) {
	sl.l.Log("en SayHello para " + userID)
	name, ok := sl.ds.UserNameForID(userID)
	if !ok {
		return "", errors.New("Usuario desconocido")
	}
	return "Hola, " + name, nil
}
func (sl SimpleLogic) SayGoodbye(userID string) (string, error) {
	sl.l.Log("en SayGoodbye para " + userID)
	name, ok := sl.ds.UserNameForID(userID)
	if !ok {
		return "", errors.New("Usuario desconocido")
	}
	return "Adios, " + name, nil
}
