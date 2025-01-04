package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Account struct {
	Name  string
	Email string
}

type AccountHandler struct {
	notifier []AccountNotifier
}

type EmailNotifier struct{}
type SMSNotifier struct{}

type AccountNotifier interface {
	NotifyAccountCreated(context context.Context, account Account) error
}

func (e *EmailNotifier) NotifyAccountCreated(context context.Context, account Account) error {
	fmt.Println("Account Notified by email:", account)
	return nil
}

func (s *SMSNotifier) NotifyAccountCreated(context context.Context, account Account) error {
	fmt.Println("Account Notified by sms :", account)
	return nil
}

func (h *AccountHandler) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var account Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if account.Name == "" || account.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}
	//Notify account creation
	for _, n := range h.notifier {
		if err := n.NotifyAccountCreated(r.Context(), account); err != nil {
			log.Fatalf("Error notifying account creation: %v", err)
		}
	}
	// Save account to database
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
	
}
func main() {

	// Create a new account handler
	accountHandler := AccountHandler{
		notifier: []AccountNotifier{
			&EmailNotifier{}, &SMSNotifier{},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /createaccount", accountHandler.handleCreateAccount)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err)
	}
}


package main

import (
	"fmt"
	"reflect"
)

type ID struct {
	id int
}
type Address struct {
	Street string
	City   string
}
type Employee struct {
	Name    string
	Age     int
	ID      ID
	Address Address
}
type Company struct {
	Name      string
	Employees []Employee
	Address   Address
	Rating    // anonymous field
}

type Rating struct {
	Stars int
}

func main() {

	company := Company{
		Name: "Golang Inc",
		Employees: []Employee{
			Employee{
				Name:    "John Doe",
				Age:     40,
				ID:      ID{id: 111},
				Address: Address{Street: "1234 Elm St", City: "Denver"},
			},
			Employee{
				Name:    "Max Doe",
				Age:     50,
				ID:      ID{id: 222},
				Address: Address{Street: "9999 Elm St", City: "Denver"},
			}},
		Address: Address{Street: "999", City: "NY"},
	}

	v := reflect.ValueOf(company)

	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.String:
			fmt.Println(v.Type().Field(i).Name, v.Field(i).String())
		case reflect.Map:
			for _, key := range v.Field(i).MapKeys() {
				value := v.Field(i).MapIndex(key)
				fmt.Println(key.Interface(), value.Interface())
			}
		case reflect.Slice, reflect.Array:
			for j := 0; j < v.Field(i).Len(); j++ {
				if v.Field(i).Index(j).Kind() == reflect.Struct {
					fmt.Println(v.Type().Field(i).Name, v.Field(i).Index(j).Interface())
				}
			}
		case reflect.Struct:
			for j := 0; j < v.Field(i).NumField(); j++ {
				fmt.Println(v.Field(i).Type().Field(j).Name, v.Field(i).Field(j).Interface())
			}
		default:
			fmt.Println(v.Field(i).Kind(), v.Type().Field(i).Name, v.Field(i).Interface())
		}
	}

	company.addEmployee("David Doe", 40, 333, "1234 Elm St", "Denver")

	for _, emp := range company.Employees {
		fmt.Println(emp)
	}

	company.printRating() // calling method on anonymous field

}

func (c *Company) addEmployee(name string, age int, id int, street string, city string) {
	//add employee to company
	c.Employees = append(c.Employees, Employee{ //(*c).Employees = append(c.Employees, Employee{
		Name: name,
		Age:  age,
		ID:   ID{id},
		Address: Address{
			Street: street,
			City:   city,
		},
	})
}

func (r *Rating) printRating() {
	r.Stars = 5
	fmt.Println("Rating:", r.Stars)
}



// Anonymous struct https://blog.boot.dev/golang/anonymous-structs-golang/
newCar := struct {
	make    string
	model   string
	mileage int
	}{
	make:    "Ford",
	model:   "Taurus",
	mileage: 200000,
}


package main

import (
	"fmt"
	"net/http"
)

type Shaper interface {
	area() float64
}

type Circle struct {
	radius float64
}

type Rectangle struct {
	length, width float64
}

type MyHandler struct{}

func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func (c *Circle) area() float64 {
	return 3.14 * c.radius * c.radius

}

func (r *Rectangle) area() float64 {
	return r.length * r.width
}

type SalaryCalculator interface {
	CalculateSalary() int
}

type Companydetails interface {
	PrintDetails() string
}

type Permanent struct {
	empId    int
	basicpay int
	pf       int
	company  Companydetails
}

type Contract struct {
	empId    int
	basicpay int
	company  Companydetails
}

func (p Permanent) CalculateSalary() int {
	return p.basicpay + p.pf
}

func (c Contract) CalculateSalary() int {
	return c.basicpay
}

func (p Permanent) PrintDetails() string {
	return fmt.Sprintf("Interface :: Employee ID: %d Basic Pay: %d PF: %d", p.empId, p.basicpay, p.pf)
}

func (c Contract) PrintDetails() string {
	return fmt.Sprintf("Interface :: Employee ID: %d Basic Pay: %d", c.empId, c.basicpay)
}

func totalExpense(s []SalaryCalculator) {
	expense := 0
	for _, v := range s {
		expense = expense + v.CalculateSalary()
	}
	fmt.Println("Total Expense Per Month $", expense)
}

func calculateArea(s Shaper) float64 {
	switch s.(type) { // type switch to check the type of the interface
	case *Circle:
		return s.area()
	case *Rectangle:
		return s.area()
	default:
		fmt.Print("Unknown shape")
	}
	return s.area()
}

func main() {
	c := Circle{radius: 5}
	r := Rectangle{length: 5, width: 10}
	shapes := []Shaper{&c, &r}
	for _, shape := range shapes {
		fmt.Println(shape.area())
	}

	//http.Handle("/", MyHandler{})
	// http.ListenAndServe(":8080", nil)

	perm := Permanent{1, 3000, 20, nil}
	cont := Contract{2, 3000, nil}

	//employees := []Permanent{perm}
	employees := []SalaryCalculator{perm, cont}
	for _, emp := range employees {
		switch emp.(type) {
		case Permanent:
			fmt.Println("Permanent Employee", emp)
		case Contract:
			fmt.Println("Contract Employee", emp)

		}
	}

	fmt.Println(perm.PrintDetails())
	fmt.Println(cont.PrintDetails())

}








