package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"identwork-scripts/schemas"

	"github.com/xuri/excelize/v2"
)

const (
	OUTPUT_DIR = "output"
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

func employeesToTxt(employees schemas.Employees) error {
	timestamp := time.Now().Format("2-Jan-2006 15:04:05")
	frontName := fmt.Sprintf("%s/%s/front-%s.txt", OUTPUT_DIR, timestamp, timestamp)
	backName := fmt.Sprintf("%s/%s/back-%s.txt", OUTPUT_DIR, timestamp, timestamp)

	if err := createDirectories(timestamp); err != nil {
		return err
	}
	if err := createFrontFile(employees, frontName); err != nil {
		return err
	}
	if err := createBackFile(employees, backName); err != nil {
		return err
	}
	return nil
}

func createDirectories(timestamp string) error {
	if _, err := os.Stat(OUTPUT_DIR); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(OUTPUT_DIR, os.ModePerm)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(OUTPUT_DIR + "/" + timestamp); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(OUTPUT_DIR+"/"+timestamp, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func createFrontFile(employees schemas.Employees, path string) error {
	front, err := os.Create(path)
	if err != nil {
		return err
	}
	defer front.Close()

	fmt.Fprint(front, "matricula\tnome_guerra\tcargo\tlotacao\tfoto\tmostrar_foto\n")
	for _, employee := range employees {
		fmt.Fprintf(front, "%s\t", employee.ID)
		fmt.Fprintf(front, "%s\t", employee.WarName)
		fmt.Fprintf(front, "%s\t", employee.Role)
		fmt.Fprintf(front, "%s\t", employee.Workplace)
		fmt.Fprintf(front, "%s\t", "true")
		fmt.Fprintf(front, "%s\t", "true")
		fmt.Fprintln(front)
	}
	return nil
}

func createBackFile(employees schemas.Employees, path string) error {
	back, err := os.Create(path)
	if err != nil {
		return err
	}
	defer back.Close()

	fmt.Fprint(back, "matricula\tnome\tidentidade\tadmissao\n")
	for _, employee := range employees {
		fmt.Fprintf(back, "%s\t", employee.ID)
		fmt.Fprintf(back, "%s\t", employee.Name)
		fmt.Fprintf(back, "%s\t", employee.Identification)
		fmt.Fprintf(back, "%s\t", employee.AdmissionDate)
		fmt.Fprintln(back)
	}
	return nil
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

	if err := employeesToTxt(employees); err != nil {
		fmt.Println(err)
		return
	}
}
