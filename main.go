package main

import (
	"fmt"

	"identwork-scripts/schemas"

	"github.com/xuri/excelize/v2"
)

const (
	ROOT_PHOTOS_FOLDER   = "C:/Users/Mareg/OneDrive/Maré Gráfica/Clientes/LAP/Crachás/Orgaos/"
	SUFFIX_PHOTOS_FOLDER = "fotos/3x4"
)

func readExcel(filename, sheet string) ([][]string, error) {
	file, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rows, err := file.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func parseDataToEmployees(employeesData [][]string) schemas.Employees {
	employees := make([]schemas.Employee, len(employeesData))

	for i, e := range employeesData {
		employee := schemas.Employee{
			ID:             e[0],
			Name:           e[1],
			WarName:        e[2],
			Role:           e[3],
			Identification: e[4],
			AdmissionDate:  e[5],
			Workplace:      e[6],
			Company:        e[7],
		}
		employees[i] = employee
	}
	return employees
}

func EmployeesToTxt(employees schemas.Employees) {
}

func main() {
	employeesData, err := readExcel("employees.xlsx", "Cards")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(employeesData)
	employees := parseDataToEmployees(employeesData[1:])
	fmt.Println(employees[0].ID)
}
