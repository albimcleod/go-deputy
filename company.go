package godeputy

//Company is the struct for a Deputy Site
type Company struct {
	ID          int    `json:"Id"`
	CompanyName string `json:"CompanyName"`
}

//Companies is the struct for a list of Company
type Companies []Company
