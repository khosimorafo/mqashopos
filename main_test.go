package main_test

import (
	"os"
	"testing"

	"."
	"net/http"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"fmt"
)

var a main.App
var tenant_id string

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize()

	//ensureTableExists()

	code := m.Run()

	//clearTable()

	os.Exit(code)
}


func TestCreateAndDeleteTenant(t *testing.T) {

	// Create Tenant
	payload := []byte(`{"name":"test tenant","zaid":"0823414062123"}`)

	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	fmt.Print(m)

	if m["name"] != "test tenant" {
		t.Errorf("Expected tenants name to be 'test product'. Got '%v'", m["name"])
	}

	if m["zaid"] != "0823414062123" {
		t.Errorf("Expected tenant mobile to be '0823414062123'. Got '%v'", m["zaid"])
	}

	tenant_id := m["id"]
	t.Log(tenant_id)

	resource := fmt.Sprintf("/tenant/%s", tenant_id)

	req, _ = http.NewRequest("DELETE", resource, nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestCreateUpdateAndDeleteProduct(t *testing.T) {

	t.Log(tenant_id)

	// Create Tenant
	payload := []byte(`{"name":"test tenant","zaid":"0823414062123"}`)

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