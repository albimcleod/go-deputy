package godeputy

import (
	"time"
)

//Roster is the struct for a Deputy roster
type Roster struct {
	ID                 int64     `json:"Id"`
	Date               time.Time `json:"Date"`
	StartTime          int64     `json:"StartTime"`
	EndTime            int64     `json:"EndTime"`
	MetaData           MetaData  `json:"_DPMetaData"`
	MatchedByTimesheet int       `json:"MatchedByTimesheet"`
	Comment            string    `json:"Comment"`
	MatchedTimesheet   Timesheet `json:"MatchedByTimesheetObject"`
	Mealbreak          time.Time `json:"Mealbreak"`
	TotalTime          float64   `json:"TotalTime"`
}

// GetStartTime will
func (obj *Roster) GetStartTime() time.Time {
	date := time.Now()
	if obj.MatchedByTimesheet == 1 {
		date = time.Unix(obj.MatchedTimesheet.StartTime, 0)
	} else {
		date = time.Unix(obj.StartTime, 0)
	}

	timezone, _ := time.LoadLocation("Australia/Adelaide")

	return date.In(timezone)
}

// GetEndTime will
func (obj *Roster) GetEndTime() time.Time {
	date := time.Now()
	if obj.MatchedByTimesheet == 1 {
		date = time.Unix(obj.MatchedTimesheet.EndTime, 0)
	}

	date = time.Unix(obj.EndTime, 0)
	timezone, _ := time.LoadLocation("Australia/Adelaide")

	return date.In(timezone)
}

// GetTotalTime will
func (obj *Roster) GetTotalTime() float64 {
	if obj.MatchedByTimesheet == 1 {
		return obj.MatchedTimesheet.TotalTime * 60
	}

	return obj.TotalTime * 60
}

//Rosters is the struct for a list of Roster
type Rosters []Roster

//MetaData is the struct for a MetaData of Roster
type MetaData struct {
	OperationalUnit OperationalUnit `json:"OperationalUnitInfo"`
	EmployeeInfo    EmployeeInfo    `json:"EmployeeInfo"`
}

type EmployeeInfo struct {
	ID          int    `json:"Id"`
	DisplayName string `json:"DisplayName"`
	Employee    int    `json:"Employee"`
}

type OperationalUnit struct {
	ID                    int    `json:"Id"`
	OperationalUnitName   string `json:"OperationalUnitName"`
	Company               int    `json:"Company"`
	CompanyName           string `json:"CompanyName"`
	ParentOperationalUnit int    `json:"Company"`
}

type Timesheet struct {
	ID        int64     `json:"Id"`
	Date      time.Time `json:"Date"`
	StartTime int64     `json:"StartTime"`
	EndTime   int64     `json:"EndTime"`
	Mealbreak time.Time `json:"Mealbreak"`
	TotalTime float64   `json:"TotalTime"`
	//	MealbreakSlots string    `json:"MealbreakSlots"`
}
