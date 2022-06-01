package public

import (
	"github.com/gorilla/mux"
	"net/http"
)

func GetRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/without-auth", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("no auth required\n"))
	}).Methods("GET")

	return r
}
