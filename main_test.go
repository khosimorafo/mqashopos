package main

import (
	"os"
	"testing"
	//"."
	"net/http"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"fmt"
	//"strconv"
	//"strconv"
)

var a App
var tenant_id string

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize()

	//ensureTableExists()

	code := m.Run()

	//clearTable()

	os.Exit(code)
}

/*
func TestCreateAndDeleteTenant(t *testing.T) {

	// Create Tenant
	payload := []byte(`{"first_name":"http","last_name":"test-tenant","zaid":"0823414062123",
		"telephone":"0123456789", "mobile":"0833459876","move_in_date":"2017-04-25"}`)

	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var tenant map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &tenant)

	fmt.Print(tenant)

	if tenant["name"] != "http test-tenant" {
		t.Errorf("Expected tenants name to be 'http test tenant'. Got '%v'", tenant["name"])
	}

	if tenant["zaid"] != "0823414062123" {
		t.Errorf("Expected tenant ZAID to be '0823414062123'. Got '%v'", tenant["zaid"])
	}

	if tenant["mobile"] != "0833459876" {
		t.Errorf("Expected tenant mobile to be '0833459876'. Got '%v'", tenant["mobile"])
	}

	if tenant["telephone"] != "0123456789" {
		t.Errorf("Expected tenant telephone to be '0123456789'. Got '%v'", tenant["telephone"])
	}

	tenant_id := tenant["id"]
	t.Log("Created id is : ", tenant_id)

	resource := fmt.Sprintf("/tenant/%s", tenant_id)

	req, _ = http.NewRequest("DELETE", resource, nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusAccepted, response.Code)
}


func TestCreateUpdateAndDeleteTenant(t *testing.T) {

	t.Log("Read id is : ", tenant_id)

	// Create Tenant
	payload := []byte(`{"first_name":"http","last_name":"test-tenant","zaid":"0823414062123",
		"telephone":"0123456789", "mobile":"0833459876","move_in_date":"2017-04-25"}`)

	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	tenant_id := m["id"]
	t.Log(tenant_id)
	resource := fmt.Sprintf("/tenant/%s", tenant_id)

	req1, _ := http.NewRequest("GET", resource, nil)

	response1 := executeRequest(req1)
	var originalTenant map[string]interface{}
	json.Unmarshal(response1.Body.Bytes(), &originalTenant)

	payload_updates := []byte(`{"name":"test tenant - update tenant"}`)

	resource_updated := fmt.Sprintf("/tenant/%s", tenant_id)
	req_updated, _ := http.NewRequest("PUT", resource_updated, bytes.NewBuffer(payload_updates))
	response_updated := executeRequest(req_updated)

	checkResponseCode(t, http.StatusAccepted, response_updated.Code)

	var m_updated map[string]interface{}
	json.Unmarshal(response_updated.Body.Bytes(), &m_updated)

	t.Log(originalTenant)
	t.Log("Updated ID is ", m_updated["id"])

	if m_updated["id"] != originalTenant["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalTenant["id"], m_updated["id"])
	}

	if m_updated["name"] == originalTenant["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalTenant["name"], m_updated["name"], m_updated["name"])
	}

	// Delete
	resource_delete := fmt.Sprintf("/tenant/%s", tenant_id)

	req, _ = http.NewRequest("DELETE", resource_delete, nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusAccepted, response.Code)
}
*/

/*
func TestCreateTenantWithInitialInvoice(t *testing.T) {

	// Create Tenant
	payload := []byte(`{"name":"http test tenant","zaid":"0823414062123", "move_in_date":"2017-05-13"}`)

	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	tenant_id := m["id"]

	t.Log(tenant_id)
	resource := fmt.Sprintf("/tenant/%s/create_first_invoice", tenant_id)

	req_inv, _ := http.NewRequest("GET", resource, nil)

	response_inv := executeRequest(req_inv)
	var invoice map[string]interface{}
	json.Unmarshal(response_inv.Body.Bytes(), &invoice)
	invoice_id := invoice["id"]

	t.Log(invoice)

	checkResponseCode(t, http.StatusCreated, response.Code)

	// Delete Invoice
	invoice_delete := fmt.Sprintf("/tenant/%s/invoice/%s", tenant_id, invoice_id)

	invoice_delete_req, _ := http.NewRequest("DELETE", invoice_delete, nil)
	invoice_delete_response := executeRequest(invoice_delete_req)

	checkResponseCode(t, http.StatusAccepted, invoice_delete_response.Code)

	// Delete Tenant
	tenant_delete := fmt.Sprintf("/tenant/%s", tenant_id)

	tenant_delete_request, _ := http.NewRequest("DELETE", tenant_delete, nil)
	tenant_delete_response := executeRequest(tenant_delete_request)

	checkResponseCode(t, http.StatusAccepted, tenant_delete_response.Code)

}
*/

/*
func TestCreateTenantWithInvoiceAndPayment(t *testing.T) {

	// Create Tenant

	payload := []byte(`{"first_name":"http","last_name":"test-tenant","zaid":"0823414062123",
		"telephone":"0123456789", "mobile":"0833459876","move_in_date":"2017-04-25"}`)

	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var tenant map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &tenant)
	tenant_id := tenant["id"]

	t.Log(tenant_id)
	// Create Invoice
	resource := fmt.Sprintf("/tenant/%s/create_first_invoice", tenant_id)
	req_inv, _ := http.NewRequest("GET", resource, nil)

	response_inv := executeRequest(req_inv)
	var invoice map[string]interface{}
	json.Unmarshal(response_inv.Body.Bytes(), &invoice)
	invoice_id := invoice["id"]
	//var invoice_balance float64
	invoice_balance, _ := invoice["balance"].(float64)
	//var str_balance string
	//str_balance = fmt.Sprintf( invoice_balance)
	str_balance := strconv.FormatFloat(invoice_balance,  'E', -1, 64)
	t.Log(invoice)

	checkResponseCode(t, http.StatusCreated, response.Code)

	// Create Payment
	str_pay := fmt.Sprintf(`{"invoice_id":"%s","amount":%s, "payment_date":"2017-06-13", "payment_mode":"cash"}`,
				invoice_id, str_balance)

	fmt.Printf("Payment String : ", str_pay)

	paymentload := []byte(str_pay)

	pay_resource := fmt.Sprintf("/tenant/%s/payments", tenant_id)
	req_pay, _ := http.NewRequest("POST", pay_resource,  bytes.NewBuffer(paymentload))

	response_pay := executeRequest(req_pay)
	var payment map[string]interface{}
	json.Unmarshal(response_pay.Body.Bytes(), &payment)
	payment_id := payment["id"]

	t.Log(payment)

	checkResponseCode(t, http.StatusCreated, response.Code)

	// Delete Payment
	payment_delete := fmt.Sprintf("/tenant/%s/payment/%s", tenant_id, payment_id)

	payment_delete_req, _ := http.NewRequest("DELETE", payment_delete, nil)
	payment_delete_response := executeRequest(payment_delete_req)

	checkResponseCode(t, http.StatusAccepted, payment_delete_response.Code)


	// Delete Invoice
	invoice_delete := fmt.Sprintf("/tenant/%s/invoice/%s", tenant_id, invoice_id)

	invoice_delete_req, _ := http.NewRequest("DELETE", invoice_delete, nil)
	invoice_delete_response := executeRequest(invoice_delete_req)

	checkResponseCode(t, http.StatusAccepted, invoice_delete_response.Code)

	// Delete Tenant
	tenant_delete := fmt.Sprintf("/tenant/%s", tenant_id)

	tenant_delete_request, _ := http.NewRequest("DELETE", tenant_delete, nil)
	tenant_delete_response := executeRequest(tenant_delete_request)

	checkResponseCode(t, http.StatusAccepted, tenant_delete_response.Code)

}

*/

func TestCreateTenantAndReadInvoices(t *testing.T) {

	// Create Tenant
	payload := []byte(`{"first_name":"http","last_name":"test-tenant","zaid":"0823414062123",
		"telephone":"0123456789", "mobile":"0833459876","move_in_date":"2017-04-25"}`)

	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	tenant_id := m["id"]

	t.Log(tenant_id)
	resource := fmt.Sprintf("/tenant/%s/create_first_invoice", tenant_id)
	req_inv, _ := http.NewRequest("GET", resource, nil)

	response_inv := executeRequest(req_inv)
	var invoice map[string]interface{}
	json.Unmarshal(response_inv.Body.Bytes(), &invoice)
	invoice_id := invoice["id"]

	t.Log(invoice)

	checkResponseCode(t, http.StatusCreated, response.Code)

	// Delete Invoice
	invoice_delete := fmt.Sprintf("/tenant/%s/invoice/%s", tenant_id, invoice_id)

	invoice_delete_req, _ := http.NewRequest("DELETE", invoice_delete, nil)
	invoice_delete_response := executeRequest(invoice_delete_req)

	checkResponseCode(t, http.StatusAccepted, invoice_delete_response.Code)

	// Delete Tenant
	tenant_delete := fmt.Sprintf("/tenant/%s", tenant_id)

	tenant_delete_request, _ := http.NewRequest("DELETE", tenant_delete, nil)
	tenant_delete_response := executeRequest(tenant_delete_request)

	checkResponseCode(t, http.StatusAccepted, tenant_delete_response.Code)

}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}