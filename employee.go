package godeputy

import "time"

//Employee is the struct for a Deputy Employee
type Employee struct {
	ID          int       `json:"Id"`
	FirstName   string    `json:"FirstName"`
	LastName    string    `json:"LastName"`
	DateOfBirth time.Time `json:"DateOfBirth"`
	Active      bool      `json:"Active"`
	Contact     Contact   `json:"ContactObject"`
}

//Employees is the struct for a list of Employee
type Employees []Employee

type Contact struct {
	Email1 string `json:"FirstName"`
}
