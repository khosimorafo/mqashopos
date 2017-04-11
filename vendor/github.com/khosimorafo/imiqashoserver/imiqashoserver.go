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

func CreateFinancialPeriodRange (start_date string, no_of_months int) (error) {

	collection := AppCollection().DB("feerlaroc").C("periods")

	t, err := now.Parse(start_date)

	if err != nil {

		log.Fatal("Date parsing error : ", err)
		return err
	}

	for i := 0; i < no_of_months; i++ {

		current := now.New(t).AddDate(0, i, 0)

		start := now.New(current).BeginningOfMonth().String()
		end := now.New(current).EndOfMonth().String()

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
