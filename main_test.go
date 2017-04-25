package main

import (
	"os"
	"testing"
	"net/http"
	"net/http/httptest"
	"fmt"
	"encoding/json"
	"github.com/antonholmquist/jason"
	"bytes"
	"strconv"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize()

	//ensureTableExists()

	code := m.Run()

	//clearTable()

	os.Exit(code)
}

/**/
func TestCreateAndDeleteTenant(t *testing.T) {

	// Create Tenant////////////////////////////////////////
	payload := []byte(`{"first_name":"http","last_name":"test-tenant","zaid":"0823414062123",
		"telephone":"0123456789", "mobile":"0833459876","move_in_date":"2017-05-01", "site":"mganka", "room":"A4" }`)

	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	b, _ := json.Marshal(result)
	v, _ := jason.NewObjectFromBytes(b)
	code, _ := v.GetInt64("code")
	ten, _ := v.GetObject("tenant")

	if code != int64(23) {
		t.Errorf("Expected result code 23. Got '%v'", code)
	}

	tenant_id,  _ :=ten.GetString("id")


	//Get tenant

	resource := fmt.Sprintf("/tenant/%s", tenant_id)
	req_get, _ := http.NewRequest("GET", resource, nil)

	response_get := executeRequest(req_get)

	checkResponseCode(t, http.StatusOK, response_get.Code)

	var result_get map[string]interface{}
	json.Unmarshal(response_get.Body.Bytes(), &result_get)

	b_get, _ := json.Marshal(result_get)
	v_get, _ := jason.NewObjectFromBytes(b_get)
	code_get, _ := v_get.GetInt64("code")
	tenant, _ := v_get.GetObject("tenant")

	if code_get != int64(21) {
		t.Errorf("Expected result code 21. Got '%v'", code)
	}

	name, _ :=tenant.GetString("name")
	if name != "http test-tenant" {
		t.Errorf("Expected tenants name to be 'http test tenant'. Got '%v'", name)
	}

	zaid, _ :=tenant.GetString("zaid")
	if zaid != "0823414062123" {
		t.Errorf("Expected tenant ZAID to be '0823414062123'. Got '%v'", zaid)
	}

	mobile, _ :=tenant.GetString("mobile")
	if mobile != "0833459876" {
		t.Errorf("Expected tenant mobile to be '0833459876'. Got '%v'", mobile)
	}

	telephone, _ :=tenant.GetString("telephone")
	if telephone != "0123456789" {
		t.Errorf("Expected tenant telephone to be '0123456789'. Got '%v'", zaid)
	}

	resource_del := fmt.Sprintf("/tenant/%s", tenant_id)

	req, _ = http.NewRequest("DELETE", resource_del, nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusAccepted, response.Code)
}

func TestCreateUpdateAndDeleteTenant(t *testing.T) {

	// Create Tenant////////////////////////////////////////
	payload := []byte(`{"first_name":"http","last_name":"test-tenant","zaid":"0823414062123",
		"telephone":"0123456789", "mobile":"0833459876","move_in_date":"2017-05-01", "site":"mganka", "room":"A4" }`)

	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)
	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	b, _ := json.Marshal(result)
	v, _ := jason.NewObjectFromBytes(b)
	code, _ := v.GetInt64("code")
	//message, _ := v.GetString("message")
	tenant, _ := v.GetObject("tenant")


	if code != int64(23) {
		t.Errorf("Expected tenant result code 203. Got '%v'", code)
	}
	tenant_id,  _ :=tenant.GetString("id")
	tenant_name,  _ := tenant.GetString("name")

	//Update tenant
	payload_updates := []byte(`{"name":"test tenant - update tenant"}`)

	resource_updated := fmt.Sprintf("/tenant/%s", tenant_id)
	req_updated, _ := http.NewRequest("PUT", resource_updated, bytes.NewBuffer(payload_updates))
	response_updated := executeRequest(req_updated)

	checkResponseCode(t, http.StatusAccepted, response_updated.Code)
	var result_updated map[string]interface{}
	json.Unmarshal(response_updated.Body.Bytes(), &result_updated)

	b_updated, _ := json.Marshal(result_updated)
	v_updated, _ := jason.NewObjectFromBytes(b_updated)
	code_updated, _ := v_updated.GetInt64("code")
	tenant_updated, _ := v_updated.GetObject("tenant")

	tenant_id_updated,  _ := tenant_updated.GetString("id")
	tenant_name_updated,  _ := tenant_updated.GetString("name")

	if code_updated != int64(21) {
		t.Errorf("Expected tenant update result code 21. Got '%v'", code_updated)
	}

	if tenant_id != tenant_id_updated {
		t.Errorf("Expected the id to remain the same (%v). Got %v", tenant_id, tenant_id_updated)
	}

	if tenant_name == tenant_name_updated {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", tenant_name, "test tenant - update tenant", tenant_name_updated)
	}

	// Delete
	resource_delete := fmt.Sprintf("/tenant/%s", tenant_id)

	req, _ = http.NewRequest("DELETE", resource_delete, nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusAccepted, response.Code)
}

func TestCreateTenantWithInitialInvoice(t *testing.T) {

	// Create Tenant////////////////////////////////////////
	payload := []byte(`{"first_name":"http","last_name":"test-tenant","zaid":"0823414062123",
		"telephone":"0123456789", "mobile":"0833459876","move_in_date":"2017-05-01", "site":"mganka", "room":"A4" }`)

	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	//t.Log(response.Body)

	var result_ten map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result_ten)

	//fmt.Print(result_ten)

	b, _ := json.Marshal(result_ten)
	v, _ := jason.NewObjectFromBytes(b)
	code, _ := v.GetInt64("code")
	//message, _ := v.GetString("message")
	tenant, _ := v.GetObject("tenant")

	if code != int64(23) {
		t.Errorf("Expected result code 23. Got '%v'", code)
		return
	}
	tenant_id, _ := tenant.GetString("id")
	t.Log(tenant_id)

	//Create First Invoice////////////////////////////////////////

	resource := fmt.Sprintf("/tenant/%s/create_first_invoice", tenant_id)
	req_inv, _ := http.NewRequest("GET", resource, nil)

	response_inv := executeRequest(req_inv)

	checkResponseCode(t, http.StatusCreated, response.Code)

	//t.Log(response_inv.Body)

	var result_inv map[string]interface{}
	json.Unmarshal(response_inv.Body.Bytes(), &result_inv)

	//fmt.Print(result_inv)

	b_inv, _ := json.Marshal(result_inv)
	v_inv, _ := jason.NewObjectFromBytes(b_inv)
	code_inv, _ := v_inv.GetInt64("code")
	invoice, _ := v_inv.GetObject("invoice")


	if code_inv != int64(23) {
		t.Errorf("Expected result code 203. Got '%v'", code_inv)
	}
	invoice_id, _ :=invoice.GetString("id")

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

func TestCreateTenantWithInvoiceAndPayment(t *testing.T) {

	// Create Tenant////////////////////////////////////////
	payload := []byte(`{"first_name":"http","last_name":"test-tenant","zaid":"0823414062123",
		"telephone":"0123456789", "mobile":"0833459876","move_in_date":"2017-05-01", "site":"mganka", "room":"A4" }`)

	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	//t.Log(response.Body)

	var result_ten map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result_ten)

	//fmt.Print(result_ten)

	b, _ := json.Marshal(result_ten)
	v, _ := jason.NewObjectFromBytes(b)
	code, _ := v.GetInt64("code")
	//message, _ := v.GetString("message")
	tenant, _ := v.GetObject("tenant")

	if code != int64(23) {
		t.Errorf("Expected result code 23. Got '%v'", code)
		return
	}
	tenant_id, _ := tenant.GetString("id")
	//t.Log(tenant_id)

	//Create First Invoice////////////////////////////////////////

	resource := fmt.Sprintf("/tenant/%s/create_first_invoice", tenant_id)
	req_inv, _ := http.NewRequest("GET", resource, nil)

	response_inv := executeRequest(req_inv)

	checkResponseCode(t, http.StatusCreated, response.Code)

	//t.Log(response_inv.Body)

	var result_inv map[string]interface{}
	json.Unmarshal(response_inv.Body.Bytes(), &result_inv)

	//fmt.Print(result_inv)

	b_inv, _ := json.Marshal(result_inv)
	v_inv, _ := jason.NewObjectFromBytes(b_inv)
	code_inv, _ := v_inv.GetInt64("code")
	invoice, _ := v_inv.GetObject("invoice")


	if code_inv != int64(23) {
		t.Errorf("Expected result code 203. Got '%v'", code_inv)
	}
	invoice_id, _ :=invoice.GetString("id")

	checkResponseCode(t, http.StatusCreated, response.Code)

	invoice_balance, _ := invoice.GetFloat64("balance")
	str_balance := strconv.FormatFloat((invoice_balance-30.0),  'E', -1, 64)

	//Create Payment////////////////////////////////////////
	str_pay := fmt.Sprintf(`{"invoice_id":"%s","amount":%s, "payment_date":"2017-04-13", "payment_mode":"cash"}`,
				invoice_id, str_balance)

	//t.Log("Payment String : ", str_pay)

	paymentload := []byte(str_pay)

	pay_resource := fmt.Sprintf("/tenant/%s/payments", tenant_id)
	req_pay, _ := http.NewRequest("POST", pay_resource,  bytes.NewBuffer(paymentload))

	response_pay := executeRequest(req_pay)
	var result_pay map[string]interface{}
	json.Unmarshal(response_pay.Body.Bytes(), &result_pay)

	//t.Log(result_pay)

	b_payment, _ := json.Marshal(result_pay)
	v_payment, _ := jason.NewObjectFromBytes(b_payment)
	code_payment, _ := v_payment.GetInt64("code")
	payment, _ := v_payment.GetObject("payment")

	checkResponseCode(t, http.StatusCreated, response_pay.Code)

	if code_payment != int64(23) {
		t.Errorf("Expected result code 23. Got '%v'", code_inv)
		return
	}
	checkResponseCode(t, http.StatusCreated, response.Code)

	payment_id, _ := payment.GetString("id")
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

func TestCreateTenantAndReadInvoices(t *testing.T) {

	// Create Tenant////////////////////////////////////////
	payload := []byte(`{"first_name":"http","last_name":"test-tenant","zaid":"0823414062123",
		"telephone":"0123456789", "mobile":"0833459876","move_in_date":"2017-04-21"}`)

	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	//t.Log(response.Body)

	var result_ten map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result_ten)

	//fmt.Print(result_ten)

	b, _ := json.Marshal(result_ten)
	v, _ := jason.NewObjectFromBytes(b)
	code, _ := v.GetInt64("code")
	tenant, _ := v.GetObject("tenant")

	if code != int64(23) {
		t.Errorf("Expected result code 203. Got '%v'", code)
		return
	}

	tenant_id, _ := tenant.GetString("id")

	//t.Log(tenant_id)

	//Create First Invoice////////////////////////////////////////

	resource := fmt.Sprintf("/tenant/%s/create_first_invoice", tenant_id)
	req_inv, _ := http.NewRequest("GET", resource, nil)

	response_inv := executeRequest(req_inv)

	checkResponseCode(t, http.StatusCreated, response.Code)

	//t.Log(response_inv.Body)

	var result_inv map[string]interface{}
	json.Unmarshal(response_inv.Body.Bytes(), &result_inv)

	//fmt.Print(result_inv)

	b_inv, _ := json.Marshal(result_inv)
	v_inv, _ := jason.NewObjectFromBytes(b_inv)
	code_inv, _ := v_inv.GetInt64("code")
	invoice, _ := v_inv.GetObject("invoice")


	if code_inv != int64(23) {
		t.Errorf("Expected result code 203. Got '%v'", code_inv)
	}
	invoice_id, _ :=invoice.GetString("id")

	//Create Next Invoice////////////////////////////////////////

	resource_inv_nxt := fmt.Sprintf("/tenant/%s/create_next_invoice", tenant_id)
	req_inv_nxt, _ := http.NewRequest("GET", resource_inv_nxt, nil)

	response_inv_nxt := executeRequest(req_inv_nxt)

	checkResponseCode(t, http.StatusCreated, response.Code)

	//t.Log(response_inv_nxt.Body)

	var result_inv_nxt map[string]interface{}
	json.Unmarshal(response_inv_nxt.Body.Bytes(), &result_inv_nxt)

	//fmt.Print(result_inv_nxt)

	b_inv_nxt, _ := json.Marshal(result_inv_nxt)
	v_inv_nxt, _ := jason.NewObjectFromBytes(b_inv_nxt)
	code_inv_nxt, _ := v_inv_nxt.GetInt64("code")
	invoice_nxt, _ := v_inv_nxt.GetObject("invoice")


	if code_inv_nxt != int64(23) {
		t.Errorf("Expected result code 203. Got '%v'", code_inv_nxt)
	}
	invoice_id_nxt, _ :=invoice_nxt.GetString("id")

	// Delete Next Invoice////////////////////////////////////////
	invoice_delete_nxt := fmt.Sprintf("/tenant/%s/invoice/%s", tenant_id, invoice_id_nxt)

	invoice_delete_req_nxt, _ := http.NewRequest("DELETE", invoice_delete_nxt, nil)
	invoice_delete_response_nxt := executeRequest(invoice_delete_req_nxt)

	checkResponseCode(t, http.StatusAccepted, invoice_delete_response_nxt.Code)

	// Delete First Invoice////////////////////////////////////////
	invoice_delete := fmt.Sprintf("/tenant/%s/invoice/%s", tenant_id, invoice_id)

	invoice_delete_req, _ := http.NewRequest("DELETE", invoice_delete, nil)
	invoice_delete_response := executeRequest(invoice_delete_req)

	checkResponseCode(t, http.StatusAccepted, invoice_delete_response.Code)

	// Delete Tenant////////////////////////////////////////
	tenant_delete := fmt.Sprintf("/tenant/%s", tenant_id)

	tenant_delete_request, _ := http.NewRequest("DELETE", tenant_delete, nil)
	tenant_delete_response := executeRequest(tenant_delete_request)

	checkResponseCode(t, http.StatusAccepted, tenant_delete_response.Code)

}

func TestCreateTenantWithInvalidDate(t *testing.T){

	payload := []byte(`{"first_name":"http","last_name":"test-tenant","zaid":"0823414062123",
		"telephone":"0123456789", "mobile":"0833459876","move_in_date":"2017-18-18"}`)

	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}
/**/
func TestReadSingleTenant(t *testing.T){

	resource := fmt.Sprintf("/tenant/%s", "256831000000046005")
	request, _ := http.NewRequest("GET", resource, nil)

	response := executeRequest(request)

	//t.Log(response.Body)

	checkResponseCode(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	j, _ := json.Marshal(result)
	e, _ := jason.NewObjectFromBytes(j)
	code, _ := e.GetInt64("code")

	if code != int64(21) {
		t.Errorf("Expected result code 21. Got '%v'", code)
	}
}

func TestReadSingleInvoice(t *testing.T){

	resource := fmt.Sprintf("/invoices/%s", "256831000000048033")
	request, _ := http.NewRequest("GET", resource, nil)

	response := executeRequest(request)

	//t.Log(response.Body)

	checkResponseCode(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	j, _ := json.Marshal(result)
	e, _ := jason.NewObjectFromBytes(j)
	code, _ := e.GetInt64("code")

	if code != int64(21) {
		t.Errorf("Expected result code 21. Got '%v'", code)
	}
}

func TestReadSinglePayment(t *testing.T){

	resource := fmt.Sprintf("/payments/%s", "256831000000141001")
	request, _ := http.NewRequest("GET", resource, nil)

	response := executeRequest(request)

	//t.Log(response.Body)

	checkResponseCode(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	j, _ := json.Marshal(result)
	e, _ := jason.NewObjectFromBytes(j)
	code, _ := e.GetInt64("code")

	if code != int64(21) {
		t.Errorf("Expected result code 21. Got '%v'", code)
	}
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