package main

import (
	"net/http"
	"encoding/json"
	"github.com/khosimorafo/imiqasho"

	"strconv"
	"log"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/antonholmquist/jason"
)

type App struct {

	Router *mux.Router
}

func (a *App) Initialize() {

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(port string) {

	log.Fatal(http.ListenAndServe(port, a.Router))
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

func (a *App) createTenantInvoice(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	tenant := imiqasho.Tenant{ID:id}

	result, entity, _ := imiqasho.Read(tenant)

	b, _ := json.Marshal(entity)
	v, _ := jason.NewObjectFromBytes(b)
	tenant_id, _ := v.GetString("id")
	//period_name, _ := v.GetString("period_name")

	//fmt.Printf("CreateFirstTenantInvoice id : ", tenant_id)

	ten := imiqasho.Tenant{ID:tenant_id}

	result, invoice, error := ten.CreateTenantInvoice("")

	if error != nil {

		respondWithError(w, http.StatusInternalServerError, error.Error())
		return
	} else {

		if result == "success" {

			respondWithJSON(w, http.StatusCreated, invoice)
		} else {

			respondWithJSON(w, http.StatusNotAcceptable, result)
		}
	}
}

func (a *App) createFirstTenantInvoice(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	tenant := imiqasho.Tenant{ID:id}

	result, entity, _ := imiqasho.Read(tenant)

	b, _ := json.Marshal(entity)
	v, _ := jason.NewObjectFromBytes(b)
	tenant_id, _ := v.GetString("id")
	in_date, _ := v.GetString("move_in_date")

	//fmt.Printf("CreateFirstTenantInvoice id : ", tenant_id)

	ten := imiqasho.Tenant{ID:tenant_id, MoveInDate:in_date}

	result, invoice, error := ten.CreateFirstTenantInvoice()

	if error != nil {

		respondWithError(w, http.StatusInternalServerError, error.Error())
		return
	} else {

		if result == "success" {

			respondWithJSON(w, http.StatusCreated, invoice)
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

func (a *App) getInvoices(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	tenant := imiqasho.Tenant{ID:id}

	result, invoices, err := tenant.GetInvoices(map[string]string{})

	if err != nil {

		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	if result == "success" {

		respondWithJSON(w, http.StatusOK, invoices)
	} else {

		respondWithError(w, http.StatusNotFound, "Tenant not found")
	}
}

func (a *App) deleteInvoice(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Println("Deleting invoice with id : ",id)

	p := imiqasho.Invoice{ID: id}

	var i imiqasho.EntityInterface
	i = p
	result, _ := imiqasho.Delete(i)

	if result == "success" {

		respondWithJSON(w, http.StatusAccepted, map[string]string{"result": "success"})
	} else {

		respondWithError(w, http.StatusNotFound, "Invoice not found")
	}
}

/**********Payment ***********************************************/

func (a *App) createPayment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	tenant := imiqasho.Tenant{ID:id}

	result, _, _ := imiqasho.Read(tenant)

	if result != "success" {

		respondWithError(w, http.StatusBadRequest, "Invalid tenant submitted")
		return
	}

	//b, _ := json.Marshal(entity)
	//v, _ := jason.NewObjectFromBytes(b)
	//tenant_id, _ := v.GetString("id")

	//ten := imiqasho.Tenant{ID:tenant_id}

	var p imiqasho.PaymentPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	result, payment, err := tenant.CreatePayment(p)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {

		if result == "success" {

			respondWithJSON(w, http.StatusCreated, payment)
		} else {

			respondWithJSON(w, http.StatusNotAcceptable, result)
		}
	}
}

func (a *App) getPayment(w http.ResponseWriter, r *http.Request) {

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

func (a *App) updatePayment(w http.ResponseWriter, r *http.Request) {

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

func (a *App) deletePayment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["payment_id"]

	fmt.Println("Deleting payment with id : ",id)

	p := imiqasho.Payment{ID: id}

	var i imiqasho.EntityInterface
	i = p
	result, _ := imiqasho.Delete(i)

	if result == "success" {

		respondWithJSON(w, http.StatusAccepted, map[string]string{"result": "success"})
	} else {

		respondWithError(w, http.StatusNotFound, "Payment not found")
	}
}

func (a *App) getPayments(w http.ResponseWriter, r *http.Request) {

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

/*****************************************************************/

func (a *App) initializeRoutes() {

	a.Router.HandleFunc("/tenants", a.getTenants).Methods("GET")
	a.Router.HandleFunc("/tenants", a.createTenant).Methods("POST")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}", a.getTenant).Methods("GET")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}", a.updateTenant).Methods("PUT")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}", a.deleteTenant).Methods("DELETE")

	a.Router.HandleFunc("/tenant/{id:[0-9]+}/invoices", a.getInvoices).Methods("GET")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}/invoices", a.createTenantInvoice).Methods("POST")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}/create_first_invoice", a.createFirstTenantInvoice).Methods("GET")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}/invoice/{id:[0-9]+}", a.deleteInvoice).Methods("DELETE")

	a.Router.HandleFunc("/tenant/{id:[0-9]+}/payments", a.getPayments).Methods("GET")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}/payments", a.createPayment).Methods("POST")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}/payment/{payment_id:[0-9]+}", a.getPayment).Methods("GET")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}/payment/{payment_id:[0-9]+}", a.updateTenant).Methods("PUT")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}/payment/{payment_id:[0-9]+}", a.deletePayment).Methods("DELETE")

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