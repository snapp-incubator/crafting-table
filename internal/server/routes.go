package server

import (
	_ "net/http"

	"github.com/snapp-incubator/crafting-table/internal/server/handlers"

	"github.com/gorilla/mux"
)

func initRoutes(r *mux.Router) {
	r.HandleFunc(
		"/api/db/get-structs",
		handlers.GetStructs(),
	).Methods("POST")

	r.HandleFunc(
		"/api/db/get-struct-fields",
		handlers.GetFieldsOFStruct(),
	).Methods("POST")

	r.HandleFunc(
		"/api/db/generate",
		handlers.GenerateManifestAndRepo(),
	).Methods("POST")

}
