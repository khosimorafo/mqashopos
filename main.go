package main

import (
	"encoding/json"
	"fmt"
	"github.com/khosimorafo/imiqasho"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/tenants", tenantHandler)
	http.Handle("/", r)

	/*
	tenant := imiqasho.Tenant{Name: "M Tenant", Mobile: "0832345678", ZAID: "2222222222222", Site: "Mganka", Room: "3"}
	var i imiqasho.EntityInterface
	i = tenant
	result, body, err := imiqasho.Create(i)

	if result == "success" {

		generateEntityResponse(result, body)
	} else {


	}*/
}

func tenantHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Gorilla!\n"))
}

func generateEntityResponse(result string, entity *imiqasho.EntityInterface) string {

	response := EntityResponse{CodeNumber: 0, Message: result, Tenant: entity}

	resp, _ := json.Marshal(response)
	str := string(resp)

	fmt.Printf(str)

	return str
}

type EntityResponse struct {

	CodeNumber int 		`json:"code"`
	Message string 		`json:"message"`
	Data *imiqasho.EntityInterface 		`json:"data,omitempty"`

	Tenant *imiqasho.EntityInterface  		`json:"tenant,omitempty"`
	Invoice *imiqasho.EntityInterface  		`json:"invoice,omitempty"`
	Item *imiqasho.EntityInterface  		`json:"item,omitempty"`
	Payment *imiqasho.EntityInterface  		`json:"payment,omitempty"`
	/*
	Tenants string 		`json:"tenants,omitempty"`
	Invoices string 	`json:"invoices,omitempty"`
	Items string 		`json:"items,omitempty"`
	Payments string 	`json:"payments,omitempty"`
	*/
}