package rModels

import (
	"gorm.io/gorm"
)

type EmployeeResponse struct {
	Page        int `json:"page"`
	TotalPages  int `json:"totalPages"`
	PageResult  int `json:"pageResult"`
	TotalResult int `json:"totalResult"`
	//RequestInfo RequestInfo `json:"requestInfo"`
	Data []Employee `json:"data"`
}

// type RequestInfo struct {
// 	CompanyName           string `json:"companyName"`
// 	CompanyIdentification string `json:"companyIdentification"`
// 	CompanyNickname       string `json:"companyNickname"`
// 	UserEmail             string `json:"userEmail"`
// }

type Employee struct {
	gorm.Model
	Code             string  `json:"code" gorm:"uniqueIndex"`
	DeviceCode       int64   `json:"deviceCode"`
	Identification   string  `json:"identification"`
	Name             string  `json:"name"`
	SecondName       string  `json:"secondName"`
	LastName         string  `json:"lastName"`
	SecondLastName   string  `json:"secondLastName"`
	BranchOfficeCode string  `json:"branchOfficeCode"`
	BranchOfficeName string  `json:"branchOfficeName"`
	DepartmentCode   string  `json:"departmentCode"`
	DepartmentName   string  `json:"departmentName"`
	EmployeeStatus   string  `json:"employeeStatus"`
	Genre            string  `json:"genre"`
	BirthDate        string  `json:"birthDate"`
	CivilStatus      string  `json:"civilStatus"`
	Address          *string `json:"address"`
	PersonalPhone    *int64  `json:"personalPhone"`
	PersonalMail     string  `json:"personalMail"`
	CorporateMail    *string `json:"corporateMail"`
	CorporatePhone   *int64  `json:"corporatePhone"`
	Favorite         bool    `json:"favorite"`
	Comment          *string `json:"comment"`
}
