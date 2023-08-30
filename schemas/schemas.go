package schemas

type Employee struct {
	ID             string
	Name           string
	WarName        string
	Role           string
	Identification string
	AdmissionDate  string
	Workplace      string
	Company        string
}

type Employees []Employee

type PhotoIdExtension map[string]string
