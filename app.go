package main

import (
	"net/http"
	"encoding/json"
	"strconv"
	"log"
	"github.com/gorilla/mux"
	"github.com/antonholmquist/jason"
	"github.com/khosimorafo/imiqasho"
	"github.com/khosimorafo/imiqashoserver"
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


/**********Tenants ***********************************************/

func (a *App) createTenant(w http.ResponseWriter, r *http.Request) {

	var p imiqasho.Tenant
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithJSON(w, http.StatusBadRequest, ResponseWrapper{Code:43, Message:"IInvalid request payload."})
		return
	}
	defer r.Body.Close()

	checksOut := checkIfMoveInDateIsValid(p.MoveInDate)

	if !checksOut {

		respondWithJSON(w, http.StatusBadRequest, ResponseWrapper{Code:43, Message:"Invalid tenant move in date."})
		return
	}

	var i imiqasho.EntityInterface
	i = p
	result, body, err := imiqasho.Create(i)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {

		if result == "success" {

			respondWithJSON(w, http.StatusCreated, ResponseWrapper{Code:23, Message:"success", Tenant:body})

		} else {

			respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:43, Message:err.Error()})
		}
	}
}

func (a *App) getTenant(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	tenant := imiqasho.Tenant{ID:id}
	var i imiqasho.EntityInterface
	i = tenant
	result, ten, error := imiqasho.Read(i)

	if error != nil {

		respondWithError(w, http.StatusInternalServerError, error.Error())
		return
	} else {

		if result == "success" {

			respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:21, Message:"success", Tenant:ten})

		} else {

			respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:43, Message:error.Error()})
		}
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
	result, tenant, error := imiqasho.Update(i)

	if error != nil {

		respondWithError(w, http.StatusInternalServerError, error.Error())
		return
	} else {

		if result == "success" {

			respondWithJSON(w, http.StatusAccepted, ResponseWrapper{Code:21, Message:"success", Tenant:tenant})
		} else {

			respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:43, Message:error.Error()})
		}
	}
}

func (a *App) deleteTenant(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	////fmt.Println("Deleting tenant with id : ",id)

	p := imiqasho.Tenant{ID: id}

	var i imiqasho.EntityInterface
	i = p
	result, error := imiqasho.Delete(i)

	if error != nil {

		respondWithError(w, http.StatusInternalServerError, error.Error())
		return
	} else {

		if result == "success" {

			respondWithJSON(w, http.StatusAccepted, map[string]string{"result": "success"})
		} else {

			respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:43, Message:error.Error()})
		}
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

	result, tenants, err := imiqasho.GetTenants(filters)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}else {

		if result == "success" {

			respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:21, Message:"success", Tenants:tenants})
		} else {

			respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:43, Message:err.Error()})
		}
	}

}

func (a *App) createFirstTenantInvoice(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	tenant := imiqasho.Tenant{ID:id}

	result, entity, err := imiqasho.Read(tenant)
	if err != nil {

		respondWithJSON(w, http.StatusNotAcceptable, "Please submit valid customer_id")
		return
	}

	if result == "failure" {

		respondWithJSON(w, http.StatusNotAcceptable, "Please submit valid customer_id")
		return
	}

	b, _ := json.Marshal(entity)
	v, _ := jason.NewObjectFromBytes(b)
	tenant_id, _ := v.GetString("id")
	in_date, _ := v.GetString("move_in_date")

	////fmt.Printf("CreateFirstTenantInvoice id : ", tenant_id)

	ten := imiqasho.Tenant{ID:tenant_id, MoveInDate:in_date}

	result, invoice, error := ten.CreateFirstTenantInvoice()


	if error != nil {

		respondWithError(w, http.StatusInternalServerError, error.Error())
		return
	} else {

		if result == "success" {
			respondWithJSON(w, http.StatusCreated, ResponseWrapper{Code:23, Message:"success", Invoice:invoice})

		} else {
			respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:43, Message:err.Error()})
		}
	}
}

func (a *App) createNextTenantInvoice(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	tenant := imiqasho.Tenant{ID:id}

	result, entity, _ := imiqasho.Read(tenant)

	b, _ := json.Marshal(entity)
	v, _ := jason.NewObjectFromBytes(b)
	tenant_id, _ := v.GetString("id")
	ten := imiqasho.Tenant{ID:tenant_id}

	result, invoice, error := ten.CreateNextTenantInvoice()

	if error != nil {

		respondWithError(w, http.StatusInternalServerError, error.Error())
		return
	} else {

		if result == "success" {
			respondWithJSON(w, http.StatusCreated, ResponseWrapper{Code:23, Message:"success", Invoice:invoice})

		} else {
			respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:43, Message:error.Error()})
		}
	}
}

func checkIfMoveInDateIsValid(date string) (bool) {

	_, _, err := imiqashoserver.DateFormatter(date)

	if err != nil {
		return false
	}
	return true
}

/**********Invoices ***********************************************/

func (a *App) getInvoices(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	tenant := imiqasho.Tenant{ID:id}

	result, invoices, err := tenant.GetInvoices(map[string]string{})

	if err != nil {

		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	if result == "success" {

		//respondWithJSON(w, http.StatusOK, invoices)
		respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:21, Message:"success", Invoices:invoices})

	} else {

		respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:43, Message:err.Error()})
	}
}

func (a *App) deleteInvoice(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	//fmt.Println("Deleting invoice with id : ",id)

	p := imiqasho.Invoice{ID: id}

	var i imiqasho.EntityInterface
	i = p
	result, err := imiqasho.Delete(i)

	if result == "success" {

		respondWithJSON(w, http.StatusAccepted, map[string]string{"result": "success"})
	} else {

		respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:43, Message:err.Error()})
	}
}

func (a *App) makePaymentExtensionRequest(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	//fmt.Println("Deleting invoice with id : ",id)

	p := imiqasho.Invoice{ID: id}

	var i imiqasho.EntityInterface
	i = p
	result, err := imiqasho.Delete(i)

	if result == "success" {

		respondWithJSON(w, http.StatusAccepted, map[string]string{"result": "success"})
	} else {

		respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:43, Message:err.Error()})
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

	var p imiqasho.PaymentPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	result, payment, err := tenant.CreatePayment(p)

	//fmt.Printf("createPayment() result is %v", result)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {

		if result == "success" {

			respondWithJSON(w, http.StatusCreated, ResponseWrapper{Code:23, Message:"success", Payment:payment})
		} else {

			respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:43, Message:err.Error()})
		}
	}
}

func (a *App) getPayment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	p := imiqasho.Payment{ID:id}

	result, payment, err := p.Read()

	if err != nil {

		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	if result == "success" {

		//respondWithJSON(w, http.StatusOK, payments)
		respondWithJSON(w, http.StatusOK, ResponseWrapper{Code: 21, Message: "success", Payment: payment})

	} else {

		respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:43, Message:err.Error()})
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
	result, body, err := imiqasho.Update(i)

	if result == "success" {

		respondWithJSON(w, http.StatusAccepted, body)
	} else {

		respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:43, Message:err.Error()})
	}
}

func (a *App) deletePayment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["payment_id"]

	//fmt.Println("Deleting payment with id : ",id)

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

	vars := mux.Vars(r)
	id := vars["id"]

	tenant := imiqasho.Tenant{ID:id}

	result, payments, err := tenant.GetPayments(map[string]string{})

	if err != nil {

		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	if result == "success" {

		//respondWithJSON(w, http.StatusOK, payments)
		respondWithJSON(w, http.StatusOK, ResponseWrapper{Code: 21, Message: "success", Payments: payments})

	} else {

		respondWithJSON(w, http.StatusOK, ResponseWrapper{Code:43, Message:err.Error()})
	}
}

/*****************************************************************/

func (a *App) initializeRoutes() {

	a.Router.HandleFunc("/tenants", a.getTenants).Methods("GET")
	a.Router.HandleFunc("/tenants", a.createTenant).Methods("POST")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}", a.getTenant).Methods("GET")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}", a.updateTenant).Methods("PUT")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}", a.deleteTenant).Methods("DELETE")

	a.Router.HandleFunc("/tenant/{id:[0-9]+}/invoices", a.getInvoices).Methods("GET")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}/invoices", a.createNextTenantInvoice).Methods("POST")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}/create_first_invoice", a.createFirstTenantInvoice).Methods("GET")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}/create_next_invoice", a.createNextTenantInvoice).Methods("GET")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}/invoice/{id:[0-9]+}", a.deleteInvoice).Methods("DELETE")
	a.Router.HandleFunc("/tenant/{id:[0-9]+}/invoice/{id:[0-9]+}", a.makePaymentExtensionRequest).Methods("POST")

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

type ResponseWrapper struct {

	Code    	int    		`json:"code,omitempty"`
	Message 	string    	`json:"message,omitempty"`
	Tenant  	interface{}    	`json:"tenant,omitempty"`
	Tenants  	interface{}    	`json:"tenants,omitempty"`
	Invoice  	interface{}    	`json:"invoice,omitempty"`
	Invoices  	interface{}    	`json:"invoices,omitempty"`
	Item  	 	interface{}    	`json:"item,omitempty"`
	Items  	 	interface{}    	`json:"items,omitempty"`
	Payment  	interface{}    	`json:"payment,omitempty"`
	Payments  	interface{}    	`json:"payments,omitempty"`

	Data  	 	interface{}    	`json:"data,omitempty"`
}