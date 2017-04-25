package imiqashoserver

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"os"
	"log"
	"github.com/jinzhu/now"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/aodin/date"
	"github.com/dariubs/percent"
	"github.com/pkg/errors"
)

type App struct {

	Session *mgo.Session
	now.Now
}

func (a *App) Initialize() {

	a.Session = AppCollection()
}

func AppCollection() (*mgo.Session) {

	uri := "mongodb://mqasho:mqasho@ds137540.mlab.com:37540/feerlaroc"
	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	//defer sess.Close()

	//sess.SetSafe(&mgo.Safe{})

	return sess;
}

type PeriodInterface interface {

	CreateFinancialPeriodRange (start_date string, no_of_months int) (error)
	ReadFinancialPeriodRange (status string) ([]Period, error)
}

type PaymentInterface interface {

	RequestStatusAsPaid() (string, error)
	RequestStatusAsApproved() (string, error)
	RequestStatusAsRejected() (string, error)
}

type EntityInterface interface {

	Create() (string, error)
	Read() (string, *EntityInterface, error)
	Update() (string, error)
	Delete() (string, error)
}

func Create(i EntityInterface) (string, error) {

	result, _ := i.Create()
	return result, nil
}

func Read(i EntityInterface) (string, *EntityInterface, error) {

	result, message, _ := i.Read()
	return result, message, nil
}

func Update(i EntityInterface) (string, error) {

	result, _ := i.Update()
	return result, nil
}

func Delete(i EntityInterface) (string, error) {

	result, err := i.Delete()
	return result, err
}

//**************************Financial Period *******************************//

type P struct {

	Date time.Time
}

type Period struct {

	Index int 	`json:"index,omitempty"`
	Name string 	`json:"name,omitempty"`
	Status string 	`json:"status,omitempty"`

	Start string 	`json:"start_date,omitempty"`
	End string 	`json:"end_date,omitempty"`
	Year int	`json:"year,omitempty"`
	Month int	`json:"month,omitempty"`
}

type LatePayment struct {

	CustomerName string 	`json:"customername,omitempty"`
	CustomerID   string 	`json:"customerid,omitempty"`
	InvoiceID    string 	`json:"invoiceid,omitempty"`
	Period 	     string 	`json:"periodname,omitempty"`
	Status 	     string 	`json:"status,omitempty"`
	Date         string 	`json:"reportdate,omitempty"`
	MustPayBy    string 	`json:"mustpaybydate,omitempty"`
}

func CreateFinancialPeriodRange (start_date string, no_of_months int) (error) {

	collection := AppCollection().DB("feerlaroc").C("periods")

	t, err := now.Parse(start_date)

	if err != nil {

		log.Fatal("Date parsing error : ", err)
		return err
	}

	for i := 0; i < no_of_months; i++ {

		current := now.New(t).AddDate(0, i, 0)

		//t.Format(time.RFC3339)
		//current := t.Format("2006-01-02")

		start := now.New(current).BeginningOfMonth().Format("2006-01-02")
		end := now.New(current).EndOfMonth().Format("2006-01-02")

		month := now.New(current).Month()
		year := now.New(current).Year()

		name := fmt.Sprintf("%s-%d", month, year)

		period := Period{i, name, "open", start,end, year, int(month)}

		collection.Insert(period)

	}

	return nil
}

func ReadFinancialPeriodRange (status string) ([]Period, error) {

	collection := AppCollection().DB("feerlaroc").C("periods")

	ps := []Period{}
	err := collection.Find(bson.M{}).All(&ps)

	if err != nil {

		return nil, err
	}

	return ps, nil
}

func RemoveFinancialPeriodRange() error {

	collection := AppCollection().DB("feerlaroc").C("periods")

	collection.RemoveAll(bson.M{})

	return nil
}

func (p *P) GetProRataDays() (float64, error)  {

	days, all, err := p.GetDaysLeft()

	if err != nil {

		return -1, err
	}

	perc := percent.PercentOf(days, all)

	return perc/100, nil
}

func (p *P) GetDaysLeft() (int, int,error)  {

	period, err := p.GetPeriod()

	if err != nil {

		return -1, -1, err
	}

	end, err1 := now.Parse(period.End)
	if err1 != nil {

		return -1, -1,  err
	}

	var no_of_days date.Range
	no_of_days.Start = date.New(p.Date.Date())
	no_of_days.End = date.New(end.Date())

	start, err2 := now.Parse(period.Start)
	if err2 != nil {

		return -1, -1,  err2
	}
	var days_in_month date.Range
	days_in_month.Start = date.New(start.Date())
	days_in_month.End = date.New(end.Date())

	return no_of_days.Days(), days_in_month.Days(), nil

}

func (p *P) GetPeriod () (Period, error) {

	actual_date := date.New(p.Date.Date())

	ps, err := ReadFinancialPeriodRange("open")

	if err != nil {

		return Period{}, err
	}

	for _, period := range ps {

		p_range := date.EntireMonth(period.Year, time.Month(period.Month))
		if actual_date.Within(p_range){

			return period, nil
		}
	}

	return Period{}, nil
}

func GetPeriodByName (name string) (Period, error) {

	ps, err := ReadFinancialPeriodRange("open")

	if err != nil {

		return Period{}, err
	}

	for _, period := range ps {

		//p_range := date.EntireMonth(period.Year, time.Month(period.Month))
		if period.Name == name{

			return period, nil
		}
	}

	return Period{}, nil
}

func GetPeriodByIndex (index int) (Period, error) {

	ps, err := ReadFinancialPeriodRange("open")

	if err != nil {

		return Period{}, err
	}

	for _, period := range ps {

		//p_range := date.EntireMonth(period.Year, time.Month(period.Month))
		if period.Index == index{

			return period, nil
		}
	}

	return Period{}, nil
}

func (period Period) GetPeriodDiscountDate() (time.Time, bool)  {

	_, p_start_t, _ := DateFormatter(period.Start)

	//Create a stub(holder) date that navigates to the previous month.
	d := time.Duration(-int(p_start_t.Day())-5) * 24 * time.Hour
	stub_date := p_start_t.Add(d)

	//Go to the beginning of the previous month and add 25 day. The result is the cut-off date/time.
	d_date := now.New(stub_date).BeginningOfMonth().AddDate(0,0,25)

	//Use cut-off date/time to create a before range
	beforeCutoff := date.Range{End: date.New(d_date.Year(), d_date.Month(), d_date.Day())}

	//Test against today. Assumes that today's date/time is the actual test date.
	today := date.FromTime(time.Now())
	if (today.Within(beforeCutoff)){

		return d_date, true
	}

	return d_date, false
}

func GetNextPeriodByName (name string) (Period, error) {

	ps, err := ReadFinancialPeriodRange("open")

	if err != nil {

		return Period{}, err
	}

	var isnext bool
	isnext = false
	for _, period := range ps {

		if isnext {

			return period, nil
		}
		//p_range := date.EntireMonth(period.Year, time.Month(period.Month))
		if period.Name == name{

			isnext = true
		}
	}

	return Period{}, nil
}

func (payment LatePayment) Create() (string, error) {

	collection := AppCollection().DB("feerlaroc").C("late_payments")

	_, err := now.Parse(payment.Date)
	if err != nil {

		log.Fatal("Date parsing error : ", err)
		return "", err
	}

	collection.Insert(payment)

	return "success", nil
}

func (payment LatePayment) Read() (string, *EntityInterface, error){

	collection := AppCollection().DB("feerlaroc").C("late_payments")

	var lp LatePayment
	err := collection.Find(bson.M{"invoiceid": payment.InvoiceID}).One(&lp)

	if err != nil{
		return "failure", nil, errors.New("Late payment request record not found.")
	}

	var p EntityInterface
	p = lp

	return "success", &p, nil
}

func (payment LatePayment) Update() (string, error){

	collection := AppCollection().DB("feerlaroc").C("late_payments")

	result, _, e := payment.Read()

	if e != nil {

		return result, e
	}

	err := collection.Remove(bson.M{"invoiceid": payment.InvoiceID})

	if err != nil{
		return "failure", errors.New("Failed to remove late payment request record.")
	}

	return "success", nil
}

func (payment LatePayment) Delete() (string, error){

	collection := AppCollection().DB("feerlaroc").C("late_payments")

	result, _, e := payment.Read()

	if e != nil {

		return result, e
	}

	err := collection.Remove(bson.M{"invoiceid": payment.InvoiceID})

	if err != nil{
		return "failure", errors.New("Failed to remove late payment request record.")
	}

	return "success", nil
}

func (payment LatePayment) RequestStatusAsApproved() (string, error){

	//Check if payment request exists
	result, _, e := payment.Read()
	if e != nil {

		return result, e
	}

	collection := AppCollection().DB("feerlaroc").C("late_payments")

	collection.Update(bson.M{"invoiceid": payment.InvoiceID}, bson.M{"$set": bson.M{"status": "approved"}})

	return "success", nil
}

func (payment LatePayment) RequestStatusAsExpired() (string, error){

	//Check if payment request exists
	result, _, e := payment.Read()
	if e != nil {

		return result, e
	}

	collection := AppCollection().DB("feerlaroc").C("late_payments")

	collection.Update(bson.M{"invoiceid": payment.InvoiceID}, bson.M{"$set": bson.M{"status": "expired"}})

	return "success", nil
}

func (payment LatePayment) RequestStatusAsPaid() (string, error){

	//Check if payment request exists
	result, _, e := payment.Read()
	if e != nil {

		return result, e
	}

	collection := AppCollection().DB("feerlaroc").C("late_payments")

	collection.Update(bson.M{"invoiceid": payment.InvoiceID}, bson.M{"$set": bson.M{"status": "paid"}})

	return "success", nil
}

func (payment LatePayment) RequestStatusAsRejected() (string, error){

	//Check if payment request exists
	result, _, e := payment.Read()
	if e != nil {

		return result, e
	}

	collection := AppCollection().DB("feerlaroc").C("late_payments")

	collection.Update(bson.M{"invoiceid": payment.InvoiceID}, bson.M{"$set": bson.M{"status": "rejected"}})

	return "success", nil
}

func (payment LatePayment) RequestStatusAsVoided() (string, error){

	//Check if payment request exists
	result, _, e := payment.Read()
	if e != nil {

		return result, e
	}

	collection := AppCollection().DB("feerlaroc").C("late_payments")

	collection.Update(bson.M{"invoiceid": payment.InvoiceID}, bson.M{"$set": bson.M{"status": "void"}})

	return "success", nil
}

func GetLatePaymentRequests(period_name string)(*[]LatePayment, error){

	collection := AppCollection().DB("feerlaroc").C("late_payments")

	_, error := GetNextPeriodByName(period_name)
	if error != nil {

		return nil, errors.New("Failed to validate submitted period_name. ")
	}

	requests := []LatePayment{}
	err := collection.Find(bson.M{"period": period_name}).All(&requests)

	if err != nil {

		return nil, err
	}

	return &requests, nil
}

//Utilities
func DateFormatter(date string) (string, time.Time, error)  {

	layout := "2006-01-02"

	t, err := time.Parse(layout, date)

	ret_t := t.Format(layout)

	if err != nil {
		fmt.Println(err)
		return "", t, errors.New("Date submitted is invalid. ")
	}


	return ret_t, t, nil
}

func RemoveLatePaymentRequests() error {

	collection := AppCollection().DB("feerlaroc").C("late_payments")

	collection.RemoveAll(bson.M{})

	return nil
}
