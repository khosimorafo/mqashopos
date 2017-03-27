package main

import (
	"github.com/gorilla/mux"

	"net/http"
	"encoding/json"
	"github.com/khosimorafo/imiqasho"

	"strconv"
	"log"
	"fmt"
)

type App struct {

	Router *mux.Router
}

func (a *App) Initialize() {

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {

	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func (a *App) createTenant(w http.ResponseWriter, r *http.Request) {

	var p imiqasho.Tenant
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	var i imiqasho.EntityInterface
	i = p
	result, body, err := imiqasho.Create(i)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {

		if result == "success" {

			respondWithJSON(w, http.StatusCreated, body)
		} else {

			respondWithJSON(w, http.StatusNotAcceptable, result)
		}
	}
}

func (a *App) getTenant(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	tenant := imiqasho.Tenant{ID:id}
	var i imiqasho.EntityInterface
	i = tenant
	result, body, _ := imiqasho.Read(i)

	if result == "success" {

		respondWithJSON(w, http.StatusOK, body)
	} else {

		respondWithError(w, http.StatusNotFound, "Tenant not found")
	}
}

func (a *App) updateTenant(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	var p imiqasho.Tenant
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	var i imiqasho.EntityInterface
	i = p
	result, body, _ := imiqasho.Update(i)

	if result == "success" {

		respondWithJSON(w, http.StatusAccepted, body)
	} else {

		respondWithJSON(w, http.StatusNotAcceptable, result)
	}

	//respondWithJSON(w, http.StatusAccepted, body)
}

func (a *App) deleteTenant(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Println("Deleting tenant with id : ",id)

	p := imiqasho.Tenant{ID: id}

	var i imiqasho.EntityInterface
	i = p
	result, _ := imiqasho.Delete(i)

	if result == "success" {

		respondWithJSON(w, http.StatusAccepted, map[string]string{"result": "success"})
	} else {

		respondWithError(w, http.StatusNotFound, "Tenant not found")
	}
}

func (a *App) getTenants(w http.ResponseWriter, r *http.Request) {

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	filters := map[string]string{}

	_ , tenants, err := imiqasho.GetTenants(filters)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tenants)
}

func (a *App) initializeRoutes() {

	a.Router.HandleFunc("/tenants", a.getTenants).Methods("GET")
	a.Router.HandleFunc("/tenants", a.createTenant).Methods("POST")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}", a.getTenant).Methods("GET")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}", a.updateTenant).Methods("PUT")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}", a.deleteTenant).Methods("DELETE")
}

func respondWithError(w http.ResponseWriter, code int, message string) {

	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}