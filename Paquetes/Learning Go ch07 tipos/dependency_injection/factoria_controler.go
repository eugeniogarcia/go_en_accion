package main

import "net/http"

//factoria. Acepta un interface y devuelve un Controller
func NewController(l Logger, logic Logic) Controller {
	return Controller{
		l:     l,
		logic: logic,
	}
}

//Controller
type Controller struct {
	l     Logger
	logic Logic
}

//Tiene la firma que espera http.HandleFunc en el servidor http
func (c Controller) SayHello(w http.ResponseWriter, r *http.Request) {
	c.l.Log("En SayHello")
	userID := r.URL.Query().Get("user_id")
	message, err := c.logic.SayHello(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(message))
}
