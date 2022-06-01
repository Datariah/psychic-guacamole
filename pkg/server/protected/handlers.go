package protected

import (
	"encoding/json"
	"fmt"
	"github.com/Datariah/psychic-guacamole/pkg/server/entities"
	"github.com/Datariah/psychic-guacamole/pkg/server/helpers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func HandlePersonGet(w http.ResponseWriter, r *http.Request) {
	db := helpers.InitializeDBConnection()

	persons := &[]entities.Person{}

	tx := db.Find(persons)

	if tx.Error != nil {
		log.Errorf("error while querying for results: %v", tx.Error.Error())
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	message, err := json.Marshal(persons)
	if err != nil {
		log.Errorf("error mashaling response: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(message)
	return
}

func HandlePersonPost(w http.ResponseWriter, r *http.Request) {
	person := entities.Person{}

	db := helpers.InitializeDBConnection()

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		log.Errorf("error while decoding body into struct: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.Create(&person)

	if tx.Error != nil {
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	message, err := json.Marshal(map[string]interface{}{
		"message":     fmt.Sprintf("user with id %d created successfully", person.ID),
		"status_code": http.StatusCreated,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(message)

	return
}

func HandlePersonPut(w http.ResponseWriter, r *http.Request) {
	db := helpers.InitializeDBConnection()

	person := &entities.Person{}

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		log.Errorf("error while decoding body into struct: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If it does, then we update it
	tx := db.Updates(person)

	if tx.Error != nil {
		log.Errorf("error while executing query: %v", tx.Error.Error())
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	} else if tx.RowsAffected == 0 {
		log.Warnf("user with ID: %d not found", person.ID)
		w.WriteHeader(http.StatusNotFound)
		message, err := json.Marshal(map[string]interface{}{
			"message":     fmt.Sprintf("user with id %d not found", person.ID),
			"status_code": http.StatusNotFound,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write(message)

		return
	}

	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	message, err := json.Marshal(person)
	if err != nil {
		log.Errorf("error mashaling response: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(message)
	return
}

func HandlePersonDelete(w http.ResponseWriter, r *http.Request) {
	db := helpers.InitializeDBConnection()

	person := &entities.Person{}

	params := mux.Vars(r)

	uid, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Errorf("error while converting back uid to uint32: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx := db.Delete(person, uid)

	if tx.Error != nil {
		log.Errorf("error while executing query: %v", tx.Error.Error())
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	} else if tx.RowsAffected == 0 {
		log.Warnf("user with id %d not found", uid)
		w.WriteHeader(http.StatusNotFound)
		message, err := json.Marshal(map[string]interface{}{
			"message":     fmt.Sprintf("user with id %d not found", uid),
			"status_code": http.StatusNotFound,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write(message)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	message, err := json.Marshal(map[string]interface{}{
		"message":     fmt.Sprintf("user with id %d deleted successfully", uid),
		"status_code": http.StatusOK,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(message)

	return
}
