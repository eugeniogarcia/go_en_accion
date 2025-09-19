package main

//Devuelve un SimpleDataStore
func NewSimpleDataStore() SimpleDataStore {
	return SimpleDataStore{
		userData: map[string]string{
			"1": "Eugenio",
			"2": "Vera Camen",
			"3": "Eugenito",
			"4": "Clara",
			"5": "Leah",
			"6": "Nico",
		},
	}
}

// Implementación concreta de DataStore
type SimpleDataStore struct {
	userData map[string]string
}

func (sds SimpleDataStore) UserNameForID(userID string) (string, bool) {
	name, ok := sds.userData[userID]
	return name, ok
}
