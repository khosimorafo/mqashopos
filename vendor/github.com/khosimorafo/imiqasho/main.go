package imiqasho

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/antonholmquist/jason"
	"github.com/smallnest/goreq"
	"github.com/khosimorafo/imiqashoserver"
)

type App struct {}

func (a *App) Initialize() {
}

type EntityInterface interface {

	Create() (string, *EntityInterface, error)
	Read() (string, *EntityInterface, error)
	Update() (string, *EntityInterface, error)
	Delete() (string, error)
}

func Create(i EntityInterface) (string, *EntityInterface, error) {

	result, message, _ := i.Create()
	return result, message, nil
}

func Read(i EntityInterface) (string, *EntityInterface, error) {

	result, message, _ := i.Read()
	return result, message, nil
}

func Update(i EntityInterface) (string, *EntityInterface, error) {

	result, message, _ := i.Update()
	return result, message, nil
}

func Delete(i EntityInterface) (string, error) {

	result, err := i.Delete()
	return result, err
}

//****************************Tenants

type TenantInterface interface {

	CreateFirstTenantInvoice() (string, *EntityInterface, error)
	CreateNextTenantInvoice() (string, *EntityInterface, error)
	CreateInvoice() (string, *EntityInterface, error)
	CreatePayment(payload PaymentPayload) (string, *EntityInterface, error)
}

type InvoiceInterface interface {

	MakePaymentExtensionRequest() (string, error)
	UpdatePaymentExtensionStatus() (string, error)
}

func CreateFirstTenantInvoice(t TenantInterface) (string, *EntityInterface, error){

	result, message, _ := t.CreateFirstTenantInvoice()
	return result, message, nil
}

func CreateNextTenantInvoice(t TenantInterface) (string, *EntityInterface, error){

	result, message, _ := t.CreateNextTenantInvoice()
	return result, message, nil
}

func CreateInvoice(t TenantInterface) (string, *EntityInterface, error){

	result, message, _ := t.CreateInvoice()
	return result, message, nil
}

type TenantZoho struct {
	ID           	string       	`json:"contact_id,omitempty"`
	Name         	string        	`json:"contact_name,omitempty"`
	Telephone   	string        	`json:"telephone,omitempty"`
	Fax          	string        	`json:"fax,omitempty"`
	Mobile       	string        	`json:"mobile,omitempty"`
	Status       	string        	`json:"status,omitempty"`
	CustomFields   []CustomField 	`json:"custom_fields,omitempty"`
	ContactPersons []ContactPerson 	`json:"contact_persons,omitempty"`
}

type Tenant struct {
	ID          string  `json:"id,omitempty"`
	Salutation  string  `json:"salutation"`
	Name        string  `json:"name"`
	FirstName   string  `json:"first_name"`
	Surname     string  `json:"last_name"`
	ZAID        string  `json:"zaid"`
	Telephone   string  `json:"telephone"`
	Fax         string  `json:"fax"`
	Mobile      string  `json:"mobile"`
	Site        string  `json:"site"`
	Room        string  `json:"room"`
	MoveInDate  string  `json:"move_in_date"`
	MoveOutDate string  `json:"move_out_date"`
	Outstanding float64 `json:"outstanding"`
	Credits     float64 `json:"credit_available"`
	Status      string  `json:"status"`
	IsPrimary   bool    `json:"is_primary_contact,omitempty"`

}

type ContactPerson struct {

	Salutation      string        `json:"salutation,omitempty"`
	Name         	string        `json:"first_name,omitempty"`
	Surname    	string        `json:"last_name,omitempty"`
	Email       	string        `json:"email,omitempty"`
	Mobile      	string        `json:"mobile,omitempty"`
	Phone       	string        `json:"phone,omitempty"`
	IsPrimary   	bool          `json:"is_primary_contact,omitempty"`
}

//A method to create new tenant

func (tenant Tenant) Create() (string, *EntityInterface, error) {

	//fmt.Printf("Creating tenant - %s with zar_id %s\n", tenant.Name, tenant.ZAID)

	contacts := make([]ContactPerson, 0)
	contacts = append(contacts, ContactPerson{tenant.Salutation, tenant.FirstName, tenant.Surname,
		tenant.Fax, tenant.Mobile, tenant.Telephone, tenant.IsPrimary})

	cfs := make([]CustomField, 0)

	cfs = append(cfs, CustomField{Index: 4, Value: tenant.ZAID})
	cfs = append(cfs, CustomField{Index: 5, Value: tenant.Site})
	cfs = append(cfs, CustomField{Index: 6, Value: tenant.Room})
	cfs = append(cfs, CustomField{Index: 7, Value: tenant.MoveInDate})
	cfs = append(cfs, CustomField{Index: 8, Value: tenant.MoveOutDate})

	name := tenant.FirstName + " " + tenant.Surname

	tenant_zoho := TenantZoho{ID: tenant.ID, Name: name, Mobile: tenant.Mobile, Fax: tenant.Fax,
		Telephone: tenant.Telephone, ContactPersons:contacts, CustomFields: cfs}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(tenant_zoho)

	resp, _, err := goreq.New().
		Post(postUrl("contacts")).
		SetHeader("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8").
		SendRawString("JSONString=" + b.String()).End()

	result, entity, error := TenantResult(resp, err)

	return result, entity, error
}

func (tenant Tenant) Read() (string, *EntityInterface, error) {

	//fmt.Printf("Retrieving tenant - %s \n", tenant.ID)

	resp, _, err := goreq.New().Get(readUrl("contacts", tenant.ID)).End()

	result, entity, error := TenantResult(resp, err)

	return result, entity, error
}

func (tenant Tenant) Update() (string, *EntityInterface, error) {

	//fmt.Printf("Updating tenant - %s\n", tenant.ID)

	cfs := make([]CustomField, 0)

	cfs = append(cfs, CustomField{Index: 4, Value: tenant.ZAID})
	cfs = append(cfs, CustomField{Index: 5, Value: tenant.Site})
	cfs = append(cfs, CustomField{Index: 6, Value: tenant.Room})

	tenant_zoho := TenantZoho{ID: tenant.ID, Name: tenant.Name}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(tenant_zoho)

	resp, _, err := goreq.New().
		Put(putUrl("contacts", tenant.ID)).
		SetHeader("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8").
		SendRawString("JSONString=" + b.String()).End()

	result, entity, error := TenantResult(resp, err)

	return result, entity, error
}

func (tenant Tenant) Delete() (string, error) {

	resp, _, err := goreq.New().Delete(deleteUrl("contacts", tenant.ID)).End()

	result, _ := jason.NewObjectFromReader(resp.Body)

	if err != nil {

		return "failure", errors.New("Failed to delete tenant. Api http error")
	} else {

		code, _ := result.GetInt64("code")

		if code == 0 {

			return "success", nil
		} else {

//			fmt.Print(result)
			return "failure", errors.New("Failed to delete tenant. Api interface error")
		}
	}
}

func (tenant Tenant) CreateFirstTenantInvoice() (string, *EntityInterface, error) {

	var i EntityInterface
	i = tenant
	result, _, read_err := Read(i)

	if read_err != nil {

		return "failure", nil, read_err
	}

	if result == "failure"{

		return "failure", nil, errors.New("Please supply valid tenant id")
	}

	layout := "2006-01-02"

//	fmt.Printf("Move in date : %v",tenant.MoveInDate)
	//ti := "2017-04-12"
	t, err := time.Parse(layout, tenant.MoveInDate)

	if err != nil {
		fmt.Println(err)
	}

//	fmt.Printf("Actual Date : %v",t)


	p := imiqashoserver.P{t}

	period, err1 := p.GetPeriod()

	if err1 != nil {

		return "failure", nil, err1
	}else {

		line_item := GetRentalLineItem()
		// Set pro rata item amount
		pr, _ := p.GetProRataDays()
		line_item.Rate = line_item.Rate * pr

		_, entity, error := tenant.CreateInvoice(tenant.MoveInDate, period, line_item)

		return "success", entity, error
	}
}

func (tenant Tenant) CreateNextTenantInvoice() (string, *EntityInterface, error) {

	//1. Retrieve and sort tenant invoices.
	filters := make(map[string]string)
	filters["customer_id"] = tenant.ID
	filters["sort_column"] = "due_date"

	_, invoices, error := tenant.GetInvoices(filters)
	if error != nil {

		return "failure", nil, errors.New("Invoice validation failure!")
	}

	var invoice Invoice
	for _, inv := range *invoices {

		invoice = inv
//		fmt.Printf("\nThen quiry period is >>> %v %v", inv.PeriodIndex, inv.CustomerName)
		break
	}

	period, _ := imiqashoserver.GetNextPeriodByName(invoice.PeriodName)

//	fmt.Printf("\n\n\n Then full period is >>> %v \n\n\n", period)


	//3. When no period exists error is derived, proceed to create new invoice
	line_item := GetRentalLineItem()
	result, entity, error := tenant.CreateInvoice("",period, line_item)

	return result, entity, error
}

func (tenant Tenant) CreateTenantInvoice(m string) (string, *EntityInterface, error) {

	period, _ := imiqashoserver.GetPeriodByName(m)

	//1. Retrieve tenant invoices.

	_, invoices, error := tenant.GetInvoices(map[string]string{})
	if error != nil {

		return "failure", nil, errors.New("Invoice validation failure!")
	}

	//2. Check if any of stored invoices has index matching intended new invoice. If there is a match, return error

	for _, invoice := range *invoices {

//		fmt.Printf("Period indexes %v %v", period.Index, invoice.PeriodIndex)

		if int64(period.Index) == invoice.PeriodIndex{

			return "failure", nil, errors.New("Invoice for the period already exists!")
		}


	}

	//3. When no period exists error is derived, proceed to create new invoice
	line_item := GetRentalLineItem()
	result, entity, error := tenant.CreateInvoice("",period, line_item)

	return result, entity, error
}

func (tenant Tenant) CreateInvoice(invoice_date string, period imiqashoserver.Period, item LineItem) (string, *EntityInterface, error) {

	layout := "2006-01-02"

	var date string
	var due string

	_, error := time.Parse(layout, invoice_date)

	if error == nil {

		date = invoice_date
		due = invoice_date
	} else {

		date, due = generateInvoiceDates(period.Start)
	}


	//Set invoice number
	length := len(tenant.ID) - 6

	var reference bytes.Buffer
	reference.WriteString(period.Name)
	reference.WriteString("-")
	reference.WriteString(tenant.ID[length:])

	item.Description = item.Description + "  " + period.Name

	//define items slice
	line_items := make([]LineItem, 0)
	line_items = append(line_items, item)

	//Must remove this hack
	var index int64
	index = int64(period.Index)

	invoice := Invoice{CustomerID: tenant.ID, InvoiceDate: date, DueDate: due,LineItems:line_items,
		ReferenceNumber: reference.String(), PeriodIndex: index, PeriodName:period.Name}

//	fmt.Printf("\n\n\nInvoice to create is >>> %v", invoice)

	var i EntityInterface
	i = invoice
	result, entity, error := Create(i)

	return result, entity, error
}

func (tenant Tenant) GetInvoices(filters map[string]string) (string, *[]Invoice, error) {

	filters["customer_id"] = tenant.ID

	resp, _, _ := goreq.New().Get(listsUrl("invoices", filters)).End()

	result, error := jason.NewObjectFromReader(resp.Body)

	//define slice
	invoices := make([]Invoice, 0)

	if error != nil {

		return "failure", &invoices, errors.New("Invoice query failure. Api http error")
	} else {

		message, _ := result.GetString("message")

		if message == "success" {

			invs, _ := result.GetObjectArray("invoices")
			for _, inv := range invs {

				invoice_id, _ := inv.GetString("invoice_id")
				customer_id, _ := inv.GetString("customer_id")
				customer_name, _ := inv.GetString("customer_name")
				reference, _ := inv.GetString("reference_number")
				due_date, _ := inv.GetString("due_date")
				invoice_date, _ := inv.GetString("date")
				balance, _ := inv.GetFloat64("balance")
				total, _ := inv.GetFloat64("total")

				cfs, _ := inv.GetObject("custom_field_hash")

				p_index, _ := cfs.GetString("cf_periodindex")
				p_name, _ := cfs.GetString("cf_periodname")

				i, e := strconv.ParseInt(p_index, 10, 64)

				if e != nil { i = 0 }

				invoice := Invoice{CustomerID: customer_id, ID:invoice_id, CustomerName:customer_name,
					ReferenceNumber:reference,InvoiceDate:invoice_date, DueDate:due_date, PeriodIndex:int64(i),
					PeriodName:p_name, Total:total, Balance:balance}

				invoices = append(invoices, invoice)
			}

			return "success", &invoices, nil
		}

		return "failure", &invoices, errors.New("Invoice query failure. Api http error")
	}
}

func (tenant Tenant) CreatePayment(payload PaymentPayload) (string, *EntityInterface, error) {

	invoice := Invoice{ID: payload.InvoiceID}

	result, entity, error := invoice.MakePayment(payload)

	return result, entity, error
}

func (tenant Tenant) GetPayments(filters map[string]string) (string, *[]Payment, error) {

	filters["customer_id"] = tenant.ID

	resp, _, _ := goreq.New().Get(listsUrl("customerpayments", filters)).End()

	result, error := jason.NewObjectFromReader(resp.Body)

	//define slice
	payments := make([]Payment, 0)

	if error != nil {

		return "failure", &payments, errors.New("Payment query failure. Api http error")
	} else {

		message, _ := result.GetString("message")

		if message == "success" {

			pymnts, _ := result.GetObjectArray("payments")
			for _, pymnt := range pymnts {

				payment_id, _ := pymnt.GetString("payment_id")
				payment_number, _ := pymnt.GetString("payment_number")
				invoice_number, _ := pymnt.GetString("invoice_number")
				mode, _ := pymnt.GetString("payment_mode")
				customer_id := tenant.ID
				amount, _ := pymnt.GetFloat64("amount")
				date, _ := pymnt.GetString("date")

				payment := Payment{ID: payment_id, CustomerID: customer_id, PaymentAmount: amount,
					PaymentDate: date, PaymentMode: mode, PaymentNumber: payment_number,
					InvoiceNumber: invoice_number}

				payments = append(payments, payment)
			}

			return "success", &payments, nil
		}

		return "failure", &payments, errors.New("Payment query failure. Api http error")
	}
}

func GetTenants(filters map[string]string) (string, *[]Tenant, error) {

	//filters := map[string]string{}

	resp, _, _ := goreq.New().Get(listsUrl("contacts", filters)).End()

	result, error := jason.NewObjectFromReader(resp.Body)

	//define slice
	tenants := make([]Tenant, 0)

	if error != nil {

		return "failure", &tenants, errors.New("Tenant query failure. Api http error")
	} else {

		message, _ := result.GetString("message")

		if message == "success" {

			contacts, _ := result.GetObjectArray("contacts")
			for _, contact := range contacts {

				customer_id, _ := contact.GetString("contact_id")
				name, _ := contact.GetString("contact_name")
				zaid, _ := contact.GetString("cf_zar_id_no")


				site, _ := contact.GetString("cf_site")
				room, _ := contact.GetString("cf_room")

				in_date, _ := contact.GetString("cf_moveindate")
				out_date, _ := contact.GetString("cf_moveoutdate")

				outstanding, _ := contact.GetFloat64("outstanding_receivable_amount")
				credit_available, _ := contact.GetFloat64("unused_credits_receivable_amount")
				status, _ := contact.GetString("status")

				contact_persons, _ := contact.GetObjectArray("contact_persons")

				var first_name string
				var last_name string
				var telephone string
				var mobile string
				for _, person := range contact_persons {

					first_name, _ = person.GetString("first_name")
					last_name, _ = person.GetString("last_name")
					telephone, _ = person.GetString("phone")
					mobile, _ = person.GetString("mobile")
				}


				tenant := Tenant{ID: customer_id, Name: name, ZAID: zaid, Telephone: telephone, Mobile: mobile,
					Site: site, Room: room, Status: status, Outstanding: outstanding, FirstName:first_name,
					Surname:last_name, Credits: credit_available, MoveInDate:in_date, MoveOutDate:out_date}

				tenants = append(tenants, tenant)
			}

			return "success", &tenants, nil
		} else {

			return "failure", &tenants, errors.New("Tenant query failure. Api interface error")
		}
	}
}

func TenantResult(response goreq.Response, err []error) (string, *EntityInterface, error) {

	if err != nil {

		return "failure", nil, errors.New("Tenant operation failure. Api http error")
	} else {

		result, _ := jason.NewObjectFromReader(response.Body)

		code, _ := result.GetInt64("code")
		message, _ := result.GetString("message")

		if code == 0 {

			contact, _ := result.GetObject("contact")

			customer_id, _ := contact.GetString("contact_id")
			name, _ := contact.GetString("contact_name")
			zaid, _ := contact.GetString("cf_zar_id_no")


			site, _ := contact.GetString("cf_site")
			room, _ := contact.GetString("cf_room")
			in_date, _ := contact.GetString("cf_moveindate")
			out_date, _ := contact.GetString("cf_moveoutdate")

			outstanding, _ := contact.GetFloat64("outstanding_receivable_amount")
			credit_available, _ := contact.GetFloat64("unused_credits_receivable_amount")
			status, _ := contact.GetString("status")

			contact_persons, _ := contact.GetObjectArray("contact_persons")

			var first_name string
			var last_name string
			var telephone string
			var mobile string
			for _, person := range contact_persons {

				first_name, _ = person.GetString("first_name")
				last_name, _ = person.GetString("last_name")
				telephone, _ = person.GetString("phone")
				mobile, _ = person.GetString("mobile")
			}


			tenant := Tenant{ID: customer_id, Name: name, ZAID: zaid, Telephone: telephone, Mobile: mobile,
			Site: site, Room: room, Status: status, Outstanding: outstanding, FirstName:first_name,
			Surname:last_name, Credits: credit_available, MoveInDate:in_date, MoveOutDate:out_date}

			var i EntityInterface
			i = tenant
			return "success", &i, nil
		} else {

			return "failure", nil, errors.New(message)
		}
	}
}

//****************************Invoices

type InvoiceZoho struct {
	ID              string     `json:"invoice_id,omitempty"`
	CustomerID      string     `json:"customer_id,omitempty"`
	ReferenceNumber string     `json:"reference_number,omitempty"`
	InvoiceDate     string     `json:"date,omitempty"`
	DueDate         string     `json:"due_date,omitempty"`
	LineItems       []LineItem `json:"line_items,omitempty"`
	Discount	float64	   `json:"discount,omitempty"`
	CustomFields []CustomField `json:"custom_fields,omitempty"`
}

type Invoice struct {
	ID              string     	`json:"id,omitempty"`
	CustomerID      string     	`json:"customer_id"`
	CustomerName      string     	`json:"customer_name"`
	InvoiceNumber  	string     	`json:"invoice_number"`
	ReferenceNumber string     	`json:"reference_number"`
	Total		float64		`json:"total"`
	Balance		float64		`json:"balance"`
	Discount	float64		`json:"discount,omitempty"`
	InvoiceDate     string     	`json:"date"`
	DueDate         string     	`json:"due_date"`
	LineItems       []LineItem 	`json:"line_items,omitempty"`
	PeriodIndex	int64		`json:"period_index,omitempty"`
	PeriodName	string 		`json:"period_name,omitempty"`
	Status          string		`json:"status,omitempty"`
}

func (invoice Invoice) Create() (string, *EntityInterface, error) {

//	fmt.Printf("Creating invoice for customer %s\n", invoice.CustomerID)

	cfs := make([]CustomField, 0)

	cfs = append(cfs, CustomField{Index: 2, Value: invoice.PeriodIndex})
	cfs = append(cfs, CustomField{Index: 3, Value: invoice.PeriodName})

	invoice_zoho := InvoiceZoho{ID: invoice.ID, CustomerID:invoice.CustomerID, ReferenceNumber:invoice.ReferenceNumber,
		InvoiceDate:invoice.InvoiceDate, DueDate:invoice.DueDate, LineItems:invoice.LineItems, CustomFields: cfs}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(invoice_zoho)


	resp, _, err := goreq.New().
		Post(postUrl("invoices")).
		SetHeader("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8").
		SendRawString("JSONString=" + b.String()).End()

	if err != nil {

		return "failure", nil, errors.New("Failed to create invoice. Api http error.")
	}


	result, entity, error := InvoiceResult(resp, err)

	return result, entity, error
}

func (invoice Invoice) Read() (string, *EntityInterface, error) {

//	fmt.Printf("Read(), Retrieving invoice - %s \n", invoice.ID)

	resp, _, err := goreq.New().Get(readUrl("invoices", invoice.ID)).End()

	result, entity, error := InvoiceResult(resp, err)

	return result, entity, error
}

func (invoice Invoice) Update() (string, *EntityInterface, error) {


	invoice_zoho := InvoiceZoho{ID: invoice.ID, LineItems:invoice.LineItems}

	result, entity, error := invoice_zoho.Update()

	return result, entity, error
}

func (invoice Invoice) Delete() (string, error) {

	resp, _, err := goreq.New().Delete(deleteUrl("invoices", invoice.ID)).End()

	result, _ := jason.NewObjectFromReader(resp.Body)

	if err != nil {

		return "failure", errors.New("Failed to delete invoice. Api http error")
	} else {

		code, _ := result.GetInt64("code")

		if code == 0 {

			return "success", nil
		} else {

//			fmt.Print(result)
			return "failure", errors.New("Failed to delete invoice. Api interface error")
		}
	}
}

func (invoice Invoice) MakePayment(payload PaymentPayload) (string, *EntityInterface, error){

	result, entity, err := invoice.Read()

	if err != nil{

		return "failure", nil, errors.New("Failed to read invoice!")
	}

	if result != "success"{

		return "failure", nil, errors.New("Failed to read invoice. Please submit valid invoice")
	}

	invoice.ProcessDiscount()

	b, _ := json.Marshal(entity)
	inv, _ := jason.NewObjectFromBytes(b)
	invoice_id, _ := inv.GetString("id")
	invoice_number, _ := inv.GetString("invoice_number")
	customer_id, _ := inv.GetString("customer_id")
	customer_name, _ := inv.GetString("customer_name")

	payment := Payment{InvoiceID: invoice_id, InvoiceNumber: invoice_number, CustomerID: customer_id, CustomerName: customer_name,
		PaymentAmount:        payload.PaymentAmount, PaymentMode: payload.PaymentMode, PaymentDate: payload.PaymentDate}

	var p EntityInterface
	p = payment
	result, entity, error := Create(p)

	if result == "success" {

		go invoice.UpdatePaymentExtensionStatusToPaid()
	}

	return result, entity, error
}

func (invoice Invoice) ProcessDiscount() {

	//Check if payment qualifies for discount
	period, _ := imiqashoserver.GetPeriodByName("May-2017")

	_, can_discount := period.GetPeriodDiscountDate()
	if can_discount{

		 _, err_disc := invoice.DiscountInvoice()
		if err_disc != nil {

			// Allow payment to go through.
			//TODO: Add an offline handler for this error.
		}
	}
}

func (invoice Invoice) DiscountInvoice() (string, error){

	item := GetRentalDiscountLineItem()

	//define items slice
	line_items := make([]LineItem, 0)
	line_items = append(line_items, item)

	inv := Invoice{ID: invoice.ID, CustomerID: invoice.CustomerID}

	_, _, error_upd := inv.AddLineItems(line_items)
	if error_upd != nil{

		return "failure", error_upd
	}

	return "success", nil
}

func (invoice Invoice) AddLineItems(new_items []LineItem) (string, *EntityInterface, error) {

	read_result, inv, err := invoice.Read()
	if err != nil {
		return read_result, nil, err
	}

	read_inv, _ := json.Marshal(inv)
	v_inv, _ := jason.NewObjectFromBytes(read_inv)
	read_inv_line_items, _ := v_inv.GetObjectArray("line_items")
	for _, item := range read_inv_line_items {

		id, _ := item.GetString("item_id")
		name, _ := item.GetString("name")
		description, _ := item.GetString("description")
		rate, _ := item.GetFloat64("rate")
		quantity, _ := item.GetInt64("quantity")

		i := LineItem{ItemID: id, Name: name, Description: description, Rate: rate, Quantity: quantity}
		new_items = append(new_items, i)
	}

	invoice_zoho := InvoiceZoho{ID: invoice.ID, LineItems:new_items}

	result, entity, error := invoice_zoho.Update()
	if error != nil {

		var error_message bytes.Buffer
		error_message.WriteString("Failed to add line items. ")
		error_message.WriteString(error.Error())
		return "failure", nil, errors.New(error_message.String())
	}

	return result, entity, nil
}

func (invoice InvoiceZoho) Update() (string, *EntityInterface, error) {

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(invoice)

//	fmt.Println(b)

	resp, _, err := goreq.New().
		Put(putUrl("invoices", invoice.ID)).
		SetHeader("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8").
		SendRawString("JSONString=" + b.String()).End()

//	fmt.Printf(resp.Status)

	result, entity, error := InvoiceResult(resp, err)

	return result, entity, error
}

func (in Invoice) MakePaymentExtensionRequest(pay_by_date string) (string, error){

	p_date, _, err := imiqashoserver.DateFormatter(pay_by_date)

	if err!=nil {

		return "", errors.New("Please submit valid pay_by_date. ")
	}

	// 1. Query the invoice, for it detail.
	inv := Invoice{ID: in.ID}
	_, entity, err := inv.Read()

	if err != nil {

		return "", errors.New("Error querying invoice. ")
	}

	invoice := (*entity).(Invoice)

//	fmt.Printf("Invoice is ... %v", invoice)

	// 2. Fill in late payment struct

	var payment imiqashoserver.LatePayment

	payment.CustomerID = invoice.CustomerID
	payment.InvoiceID = invoice.ID
	payment.CustomerName = invoice.CustomerName
	payment.Period = invoice.PeriodName
	payment.Status = "approved"

	layout := "2006-01-02"
	date := time.Now()
	payment.Date = date.Format(layout)

	payment.MustPayBy = p_date

	result, err := payment.Create()

	if err != nil{

		return "", errors.New("Error. Could not create payment extension. ")
	}

	return result, nil
}

func (in Invoice) UpdatePaymentExtensionStatusToApproved() (string, error){

	lp := imiqashoserver.LatePayment{InvoiceID:in.ID}
	result, err := lp.RequestStatusAsApproved()

	return result, err
}

func (in Invoice) UpdatePaymentExtensionStatusToRejected() (string, error){

	lp := imiqashoserver.LatePayment{InvoiceID:in.ID}
	result, err := lp.RequestStatusAsRejected()

	return result, err
}

func (in Invoice) UpdatePaymentExtensionStatusToPaid() (string, error){

	lp := imiqashoserver.LatePayment{InvoiceID:in.ID}
	result, err := lp.RequestStatusAsPaid()

	return result, err
}

func (in Invoice) UpdatePaymentExtensionStatusToExpired() (string, error){

	lp := imiqashoserver.LatePayment{InvoiceID:in.ID}
	result, err := lp.RequestStatusAsExpired()

	return result, err
}

func (in Invoice) UpdatePaymentExtensionStatusToVoided() (string, error){

	lp := imiqashoserver.LatePayment{InvoiceID:in.ID}
	result, err := lp.RequestStatusAsVoided()

	return result, err
}

func GetInvoices(filters map[string]string) (string, *[]Invoice, error) {

	resp, _, _ := goreq.New().Get(listsUrl("invoices", filters)).End()

	result, error := jason.NewObjectFromReader(resp.Body)

	//define slice
	invoices := make([]Invoice, 0)

	if error != nil {

		return "failure", &invoices, errors.New("Invoice query failure. Api http error")
	} else {

		message, _ := result.GetString("message")

		if message == "success" {

			invs, _ := result.GetObjectArray("invoices")

			for _, inv := range invs {

				customer_id, _ := inv.GetString("customer_id")
				invoice_id, _ := inv.GetString("invoice_id")
				invoice_number, _ := inv.GetString("invoice_number")
				reference_number, _ := inv.GetString("reference_number")
				invoice_date, _ := inv.GetString("date")
				due_date, _ := inv.GetString("due_date")

				cfs, _ := inv.GetObject("custom_field_hash")

				p_index, _ := cfs.GetInt64("cf_periodindex")
				p_name, _ := cfs.GetString("cf_periodname")

				total, _ := inv.GetFloat64("total")
				balance, _ := inv.GetFloat64("balance")

				status, _ := inv.GetString("status")

				invoice := Invoice{CustomerID: customer_id, ID:invoice_id, InvoiceNumber:invoice_number, Status: status,
					ReferenceNumber:reference_number,InvoiceDate:invoice_date, DueDate:due_date, PeriodIndex:p_index,
					PeriodName:p_name, Total:total, Balance:balance}

				invoices = append(invoices, invoice)
			}

			return "success", &invoices, nil
		} else {

			return "failure", &invoices, errors.New("Invoice query failure. Api interface error")
		}
	}
}

func InvoiceResult(response goreq.Response, err []error) (string, *EntityInterface, error) {

	if err != nil {

		return "failure", nil, errors.New("Invoice operation failure. Api http error")
	} else {

		result, _ := jason.NewObjectFromReader(response.Body)

		code, _ := result.GetInt64("code")
		message, _ := result.GetString("message")

		if code == 0 {

			inv, _ := result.GetObject("invoice")

			invoice_id, _ := inv.GetString("invoice_id")
			customer_id, _ := inv.GetString("customer_id")
			customer_name, _ := inv.GetString("customer_name")
			reference, _ := inv.GetString("reference_number")
			due_date, _ := inv.GetString("due_date")
			invoice_date, _ := inv.GetString("date")
			balance, _ := inv.GetFloat64("balance")
			total, _ := inv.GetFloat64("total")

			status, _ := inv.GetString("status")
			cfs, _ := inv.GetObject("custom_field_hash")

			period_index, _ := cfs.GetString("cf_periodindex")
			period_name, _ := cfs.GetString("cf_periodname")

			p_index, e := strconv.ParseInt(period_index, 10, 64)

			if e != nil { p_index = 0 }

			line_items, _ := inv.GetObjectArray("line_items")
			items := make([]LineItem, 0)

			for _, item := range line_items {

				id, _ := item.GetString("item_id")
				name, _ := item.GetString("name")
				description, _ := item.GetString("description")
				rate, _ := item.GetFloat64("rate")
				quantity, _ := item.GetInt64("quantity")

				i := LineItem{ItemID: id, Name: name, Description: description, Rate: rate, Quantity: quantity}
				items = append(items, i)
			}

			invoice := Invoice{ID: invoice_id, CustomerID: customer_id, CustomerName:customer_name,
				ReferenceNumber: reference, DueDate:       due_date, InvoiceDate: invoice_date,
				Balance:balance, Total:total, LineItems: items, PeriodIndex: p_index,
				PeriodName: period_name, Status:status}

			var i EntityInterface
			i = invoice

			return "success", &i, nil
		} else {

			return "failure", nil, errors.New(message)
		}
	}
}

//****************************Payments

type PaymentZoho struct {
	ID            string       `json:"id,omitempty"`
	CustomerID    string       `json:"customer_id"`
	PaymentAmount float64      `json:"amount"`
	PaymentMode   string       `json:"payment_mode"`
	Description   string       `json:"description"`
	PaymentDate   string   	   `json:"date"`
	Invoices      []PayInvoice `json:"invoices"`
}

type Payment struct {
	ID             		string   `json:"id,omitempty"`
	Description    		string   `json:"description,omitempty"`
	CustomerID     		string   `json:"customer_id,omitempty"`
	CustomerName   		string   `json:"customer_name,omitempty"`
	InvoiceID     	 	string   `json:"invoice_id,omitempty"`
	InvoiceNumber  		string   `json:"invoice_number,omitempty"`
	InvoiceAmount  		float64  `json:"invoice_amount,omitempty"`
	InvoiceBalance 		float64  `json:"invoice_balance,omitempty"`
	InvoiceAppliedAmount	float64  `json:"invoice_applied_amount,omitempty"`
	PaymentNumber  		string   `json:"payment_number,omitempty"`
	PaymentAmount  		float64  `json:"amount"`
	PaymentMode    		string   `json:"payment_mode,omitempty"`
	PaymentDate    		string   `json:"payment_date"`
	Status         		string   `json:"status,omitempty"`
}

/*Allow for a concise payment payload*/
type PaymentPayload struct {

	InvoiceID     	 	string   `json:"invoice_id,omitempty"`
	PaymentAmount  		float64  `json:"amount"`
	PaymentMode    		string   `json:"payment_mode,omitempty"`
	PaymentDate    		string   `json:"payment_date"`
}

type PayInvoice struct {
	InvoiceID     string  `json:"invoice_id"`
	InvoiceNumber string  `json:"invoice_number,omitempty"`
	AppliedAmount float64 `json:"amount_applied"`
}

func (payment Payment) Create() (string, *EntityInterface, error) {

//	fmt.Printf("\nCreating payment for customer %s, invoice %s\n", payment.CustomerID, payment.InvoiceID)

	payment_invoice := PayInvoice{InvoiceID:payment.InvoiceID, InvoiceNumber:payment.InvoiceNumber, AppliedAmount:payment.PaymentAmount}

	// By design, mqasho must have one payment must refer to one invoice.
	invoices := make([]PayInvoice, 1)
	invoices[0] = payment_invoice

	payment_zoho := PaymentZoho{CustomerID:payment.CustomerID, PaymentAmount:payment.PaymentAmount,
		PaymentMode:payment.PaymentMode, Description:payment.Description, PaymentDate:payment.PaymentDate,
		Invoices:invoices}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(payment_zoho)

	resp, _, err := goreq.New().
		Post(postUrl("customerpayments")).
		SetHeader("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8").
		SendRawString("JSONString=" + b.String()).End()

	result, entity, error := PaymentResult(resp, err)

	//fmt.Printf("\n Result is", entity)

	return result, entity, error
}

func (payment Payment) Read() (string, *EntityInterface, error) {

//	fmt.Printf("Retrieving payment - %s \n", payment.ID)

	resp, _, err := goreq.New().Get(readUrl("customerpayments", payment.ID)).End()

	result, entity, error := PaymentResult(resp, err)

	return result, entity, error
}

func (payment Payment) Update() (string, *EntityInterface, error) {

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(payment)

//	fmt.Println(b)

	resp, _, err := goreq.New().
		Put(putUrl("customerpayment", payment.ID)).
		SetHeader("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8").
		SendRawString("JSONString=" + b.String()).End()

	result, entity, error := PaymentResult(resp, err)

	return result, entity, error
}

func (payment Payment) Delete() (string, error) {

	resp, _, err := goreq.New().Delete(deleteUrl("customerpayments", payment.ID)).End()

	result, _ := jason.NewObjectFromReader(resp.Body)

	if err != nil {

		return "failure", errors.New("Failed to delete payment. Api http error")
	} else {

		code, _ := result.GetInt64("code")

		if code == 0 {

			return "success", nil
		} else {

//			fmt.Print(result)
			return "failure", errors.New("Failed to delete payment. Api interface error")
		}
	}
}

func PaymentResult(response goreq.Response, err []error) (string, *EntityInterface, error) {

	if err != nil {

		return "failure", nil, errors.New("Payment operation failure. Api http error")
	} else {

		result, _ := jason.NewObjectFromReader(response.Body)

		code, _ := result.GetInt64("code")
		message, _ := result.GetString("message")


		if code == 0 {

			record, _ := result.GetObject("payment")

			id, _ := record.GetString("payment_id")
			customer_id, _ := record.GetString("customer_id")
			amount, _ := record.GetFloat64("amount")
			date, _ := record.GetString("date")
			mode, _ := record.GetString("payment_mode")
			status, _ := record.GetString("status")
			description, _ := record.GetString("description")
			customer_name, _ := record.GetString("customer_name")

			invoices, _ := record.GetObjectArray("invoices")

			var invoice_id string
			var invoice_number string
			var invoice_amount float64
			var invoice_balance float64
			var invoice_applied_amount float64

			for _, invoice := range invoices {

				invoice_id, _ = invoice.GetString("invoice_id")
				invoice_number, _ = invoice.GetString("invoice_number")
				invoice_amount, _ = invoice.GetFloat64("invoice_amount")
				invoice_balance, _ = invoice.GetFloat64("balance_amount")
				invoice_applied_amount, _ = invoice.GetFloat64("applied_amount")
			}

			payment := Payment{ID: id, CustomerID: customer_id, InvoiceID: invoice_id, PaymentAmount: amount,
				PaymentDate: date, PaymentMode: mode, Status: status, Description: description,
				CustomerName: customer_name, InvoiceNumber:invoice_number, InvoiceBalance:invoice_balance,
				InvoiceAmount:invoice_amount, InvoiceAppliedAmount:invoice_applied_amount}

			var i EntityInterface
			i = payment
			return "success", &i, nil
		} else {

			return "failure", nil, errors.New(message)
		}
	}
}

//****************************Item*************************************************************//

type Item struct {
	Name        string  `json:"name"`
	Description string  `json:"description                                                                  "`
	Rate        float64 `json:"rate"`
}

type LineItem struct {
	ItemID      string  `json:"item_id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Rate        float64 `json:"rate,omitempty"`
	Quantity    int64   `json:"quantity,omitempty"`
}

func GetRentalLineItem() LineItem {

	item_id, rate, _ := getRentalItem()
	line := LineItem{ItemID: item_id, Rate: rate, Quantity: 1}

	return line
}

func GetRentalFineLineItem() LineItem {

	item_id, rate, _ := getRentalFineItem()
	line := LineItem{ItemID: item_id, Rate: rate, Quantity: 1}

	return line
}

func GetRentalDiscountLineItem() LineItem {

	item_id, rate, _ := getRentalDiscountItem()
	line := LineItem{ItemID: item_id, Rate: rate, Quantity: 1}

	return line
}

func getRentalItem() (string, float64, error) {

	resp, _, _ := goreq.New().Get(readUrl("items", "256831000000046017")).End()

	//fmt.Printf(body)

	result, error := jason.NewObjectFromReader(resp.Body)

	if error != nil {

		return "", 0, error
	} else {

		code, _ := result.GetInt64("code")
		if code == 0 {

			item, _ := result.GetObject("item")
			id, _ := item.GetString("item_id")
			rate, _ := item.GetFloat64("rate")
			return id, rate, nil
		}

		return "", 0, nil
	}
}

func getRentalFineItem() (string, float64, error) {

	resp, _, _ := goreq.New().Get(readUrl("items", "256831000000223043")).End()

	//fmt.Printf(body)

	result, error := jason.NewObjectFromReader(resp.Body)

	if error != nil {

		return "", 0, error
	} else {

		code, _ := result.GetInt64("code")
		if code == 0 {

			item, _ := result.GetObject("item")
			id, _ := item.GetString("item_id")
			rate, _ := item.GetFloat64("rate")
			return id, rate, nil
		}

		return "", 0, nil
	}
}

func getRentalDiscountItem() (string, float64, error) {

	resp, _, _ := goreq.New().Get(readUrl("items", "256831000000223405")).End()

	//fmt.Printf(body)

	result, error := jason.NewObjectFromReader(resp.Body)

	if error != nil {

		return "", 0, error
	} else {

		code, _ := result.GetInt64("code")
		if code == 0 {

			item, _ := result.GetObject("item")
			id, _ := item.GetString("item_id")
			rate, _ := item.GetFloat64("rate")
			return id, rate, nil
		}

		return "", 0, nil
	}
}

//****************************Monthly Runs*********************************//

func DoMonthlyLatePaymentFines(period_name string) (int, int, []string) {

	// create a slice for the errors
	var errstrings []string
	var no_of_succesful int
	no_of_succesful = 0

	var no_of_invoices int
	no_of_invoices = 0

	period, err_p := imiqashoserver.GetPeriodByName(period_name)
	if err_p != nil {

		err := fmt.Errorf("The period_name submitted is invalid. ")
		errstrings = append(errstrings, err.Error())
		return no_of_invoices, no_of_succesful, errstrings
	}

	invoice_date, _, _ := imiqashoserver.DateFormatter(period.Start)


	//1. Retrieve and sort tenant invoices.
	filters := make(map[string]string)
	filters["due_date_after"] = invoice_date

	_, invoices, error := GetInvoices(filters)
	if error != nil {

		err := fmt.Errorf("Failed to read invoices")
		errstrings = append(errstrings, err.Error())
		return no_of_invoices, no_of_succesful, errstrings
	}

	no_of_invoices = len(*invoices)

	item := GetRentalFineLineItem()

	requests, err := imiqashoserver.GetLatePaymentRequests(period_name)
	if err != nil {

		err := fmt.Errorf("Error while requesting late payment requests. ")
		errstrings = append(errstrings, err.Error())
		return no_of_invoices, no_of_succesful, errstrings

	}

	for _, invoice := range *invoices {

		if invoice.Status == "paid" {
			break
		}

		if invoice.Status == "partially_paid" {
			break
		}

		if invoice.PeriodName != period_name {
			break
		}

		for _, request := range *requests {
			if request.InvoiceID == invoice.ID {	break	}
		}

		//define items slice
		line_items := make([]LineItem, 0)
		line_items = append(line_items, item)

		invoice := Invoice{ID: invoice.ID, CustomerID: invoice.CustomerID}

		_, _, error_upd := invoice.AddLineItems(line_items)
		if error_upd != nil{

			err := fmt.Errorf("Error on invoice : ", invoice.ID)
			errstrings = append(errstrings, err.Error())
		} else{
			no_of_succesful++
		}
	}
	return no_of_invoices, no_of_succesful, errstrings
}

func DoMonthlyInvoiceCreation(period_name string) (string, string, error){

	filters := map[string]string{}

	result, tenants, _ := GetTenants(filters)

	if result == "success"{

		var invoices_succesfully_created int
		invoices_succesfully_created = 0

		//line_item := GetRentalLineItem()

		for _, tenant := range *tenants{

			_, _, err := tenant.CreateTenantInvoice(period_name)

			if(err == nil) {

				invoices_succesfully_created++
			}
		}

		if len(*tenants) != invoices_succesfully_created{

			return "failure", "Not all tenants were processed", nil

		}else {

			return "success", "All valid tenants processed",  nil
		}

	}

	return "failure", "", errors.New("Failed to create invoice.")
}

//****************************Common**********************************************************//

type CustomField struct {
	Index int64  `json:"index,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

func generateInvoiceDates(cur string) (string, string) {

	layout := "2006-01-02"


//	fmt.Printf("Move in date : %v",cur)
	//ti := "2017-05-12"
	t, err := time.Parse(layout, cur)

	if err != nil {
		fmt.Println(err)
	}

	// Derive invoice date
	t.Format(time.RFC3339)
	current := t.Format("2006-01-02")

	// Derive 5th of the next month. Due date
	//t2 := t.AddDate(0, 1, 0)
	d := time.Duration(-int(t.Day())+5) * 24 * time.Hour

	due := t.Add(d).Format("2006-01-02")

	return current, due
}

func generatePeriod() string {

	var buffer bytes.Buffer

	month := time.Now().AddDate(0, 1, 0).Month().String()
	year := time.Now().Year()

	buffer.WriteString(month)
	buffer.WriteString("-")
	buffer.WriteString(strconv.Itoa(year))

	return buffer.String()
}

func postUrl(entity string) string {

	apiUrl := "https://invoice.zoho.com"
	resource := fmt.Sprintf("/api/v3/%s/", entity)
	data := url.Values{}
	data.Set("authtoken", "23d96588d022f48fe2ce16dfd2b69c71")
	data.Add("organization_id", "163411778")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	u.RawQuery = data.Encode()
	urlStr := fmt.Sprintf("%v", u)

	return urlStr
}

func putUrl(entity string, id string) string {

	apiUrl := "https://invoice.zoho.com"
	resource := fmt.Sprintf("/api/v3/%s/%s", entity, id)
	data := url.Values{}
	data.Set("authtoken", "23d96588d022f48fe2ce16dfd2b69c71")
	data.Add("organization_id", "163411778")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	u.RawQuery = data.Encode()
	urlStr := fmt.Sprintf("%v", u)

	return urlStr
}

func readUrl(entity string, id string) string {

	apiUrl := "https://invoice.zoho.com"
	resource := fmt.Sprintf("/api/v3/%s/%s", entity, id)
	data := url.Values{}
	data.Set("authtoken", "23d96588d022f48fe2ce16dfd2b69c71")
	data.Add("organization_id", "163411778")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	u.RawQuery = data.Encode()
	urlStr := fmt.Sprintf("%v", u)

	return urlStr
}

func deleteUrl(entity string, id string) string {

	apiUrl := "https://invoice.zoho.com"
	resource := fmt.Sprintf("/api/v3/%s/%s", entity, id)
	data := url.Values{}
	data.Set("authtoken", "23d96588d022f48fe2ce16dfd2b69c71")
	data.Add("organization_id", "163411778")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	u.RawQuery = data.Encode()
	urlStr := fmt.Sprintf("%v", u)

	return urlStr
}

func listsUrl(entity string, filters map[string]string) string {

	apiUrl := "https://invoice.zoho.com"
	resource := fmt.Sprintf("/api/v3/%s/", entity)
	data := url.Values{}
	data.Set("authtoken", "23d96588d022f48fe2ce16dfd2b69c71")
	data.Add("organization_id", "163411778")

	for k, v := range filters {
		data.Add(k, v)
	}

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	u.RawQuery = data.Encode()
	urlStr := fmt.Sprintf("%v", u)

	return urlStr
}
