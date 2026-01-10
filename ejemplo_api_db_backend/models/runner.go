package models

type Runner struct {
	ID           string    `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Age          int       `json:"age,omitempty"`
	IsActive     bool      `json:"is_active"`
	Country      string    `json:"country"`
	PersonalBest string    `json:"personal_best,omitempty"` // se incluye el campo en el json solo si no es vacío
	SeasonBest   string    `json:"season_best,omitempty"`   // se incluye el campo en el json solo si no es vacío
	Results      []*Result `json:"results,omitempty"`       // se incluye el campo en el json solo si no es nulo o vacío
}
