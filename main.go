package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"identwork-scripts/schemas"

	"github.com/xuri/excelize/v2"
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

func validateDirectoryPath(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist")
		}
		return err
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("this path is not a directory")
	}
	filePaths, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		if len(filePaths) == 0 {
			return fmt.Errorf("directory is empty")
		}
		return err
	}
	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter the path to the photo directory")

	photoDirectoryPath, _ := reader.ReadString('\n')
	photoDirectoryPath = strings.TrimSpace(photoDirectoryPath)

	if err := validateDirectoryPath(photoDirectoryPath); err != nil {
		fmt.Println(err)
		return
	}

	employeesData, err := readExcel("employees.xlsx", "Cards")
	if err != nil {
		fmt.Println(err)
		return
	}

	employees := parseDataToEmployees(employeesData[1:])
	fmt.Println(employees)
}
