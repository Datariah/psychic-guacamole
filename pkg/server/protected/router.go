package protected

import (
	"github.com/Datariah/psychic-guacamole/pkg/server/entities"
	"github.com/Datariah/psychic-guacamole/pkg/server/helpers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func GetRouter() *mux.Router {
	r := mux.NewRouter()

	db := helpers.InitializeDBConnection()

	err := db.AutoMigrate(&entities.Person{})
	if err != nil {
		log.Panicf("error while auto-migrating entities: %v", err)
	}

	r.HandleFunc("/api/v1/persons", HandlePersonGet).Methods("GET")

	r.HandleFunc("/api/v1/persons", HandlePersonPost).Methods("POST")
	r.HandleFunc("/api/v1/persons", HandlePersonPut).Methods("PUT")
	r.HandleFunc("/api/v1/persons/{id}", HandlePersonDelete).Methods("DELETE")

	return r
}
