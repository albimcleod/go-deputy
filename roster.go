package godeputy

import "time"

//Roster is the struct for a Deputy roster
type Roster struct {
	ID                 int       `json:"Id"`
	Date               time.Time `json:"Date"`
	StartTime          int64     `json:"StartTime"`
	EndTime            int64     `json:"EndTime"`
	MetaData           MetaData  `json:"_DPMetaData"`
	MatchedByTimesheet int       `json:"MatchedByTimesheet"`
	Comment            string    `json:"Comment"`
}

// GetStartTime will
func (obj *Roster) GetStartTime() time.Time {
	return time.Unix(obj.StartTime, 0)
}

// GetEndTime will
func (obj *Roster) GetEndTime() time.Time {
	return time.Unix(obj.EndTime, 0)
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
	ID                  int    `json:"Id"`
	OperationalUnitName string `json:"OperationalUnitName"`
	Company             int    `json:"Company"`
	CompanyName         string `json:"CompanyName"`
}
