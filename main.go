package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"identwork-scripts/schemas"

	"github.com/xuri/excelize/v2"
)

const (
	OUTPUT_DIR = "output"
)

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

func employeesToTxt(employees *schemas.Employees, photoDirectoryPath string) error {
	timestamp := time.Now()
	fmt.Println(strconv.Itoa(timestamp.Nanosecond()))

	frontFileName := fmt.Sprintf(
		"%s/%s/%s/%s/%s/front-%s.txt",
		OUTPUT_DIR,
		strconv.Itoa(timestamp.Year()),
		timestamp.Month().String(),
		strconv.Itoa(timestamp.Day()),
		strconv.Itoa(timestamp.Nanosecond()),
		strconv.Itoa(timestamp.Nanosecond()),
	)
	backFileName := fmt.Sprintf(
		"%s/%s/%s/%s/%s/back-%s.txt",
		OUTPUT_DIR,
		strconv.Itoa(timestamp.Year()),
		timestamp.Month().String(),
		strconv.Itoa(timestamp.Day()),
		strconv.Itoa(timestamp.Nanosecond()),
		strconv.Itoa(timestamp.Nanosecond()),
	)
	missingPhotoFileName := fmt.Sprintf(
		"%s/%s/%s/%s/%s/missing-photo-%s.txt",
		OUTPUT_DIR,
		strconv.Itoa(timestamp.Year()),
		timestamp.Month().String(),
		strconv.Itoa(timestamp.Day()),
		strconv.Itoa(timestamp.Nanosecond()),
		strconv.Itoa(timestamp.Nanosecond()),
	)

	extensions, err := getPhotosExtensions(employees, photoDirectoryPath)
	if err != nil {
		return err
	}
	employeesWithPhoto, employeesWithoutPhoto := getEmployeesWithPhotos(*employees, extensions)
	if err := createDirectories(timestamp); err != nil {
		return err
	}
	if err := createFrontFile(&employeesWithPhoto, extensions, frontFileName, photoDirectoryPath); err != nil {
		return err
	}
	if err := createBackFile(&employeesWithPhoto, backFileName); err != nil {
		return err
	}
	if len(employeesWithoutPhoto) == 0 {
		return nil
	}
	if err := createMissingPhotoFile(&employeesWithoutPhoto, missingPhotoFileName); err != nil {
		return err
	}
	return nil
}

func getPhotosExtensions(employees *schemas.Employees, photoDirectoryPath string) (schemas.PhotoIdExtension, error) {
	extensions := make(schemas.PhotoIdExtension)
	if err := filepath.WalkDir(photoDirectoryPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}
		filename := d.Name()
		extension := filepath.Ext(filename)
		id := strings.TrimSuffix(filename, extension)
		extensions[id] = extension
		return nil
	}); err != nil {
		return nil, err
	}
	return extensions, nil
}

func getEmployeesWithPhotos(employees schemas.Employees, photosIDsExtensions schemas.PhotoIdExtension) (schemas.Employees, schemas.Employees) {
	employeesWithPhoto := make(schemas.Employees, 0, len(photosIDsExtensions))
	maxLen := len(photosIDsExtensions)
	if (len(employees) - len(photosIDsExtensions)) < 0 {
		maxLen = len(employees)
	}
	employeesWithoutPhoto := make(schemas.Employees, 0, len(employees)-maxLen)

	for _, employee := range employees {
		_, ok := photosIDsExtensions[employee.ID]
		if ok {
			employeesWithPhoto = append(employeesWithPhoto, employee)
		} else {
			employeesWithoutPhoto = append(employeesWithoutPhoto, employee)
		}
	}
	return employeesWithPhoto, employeesWithoutPhoto
}

func createDirectories(timestamp time.Time) error {
	filepath := OUTPUT_DIR
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(filepath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	filepath = filepath + "/" + strconv.Itoa(timestamp.Year())
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(filepath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	filepath = filepath + "/" + timestamp.Month().String()
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(filepath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	filepath = filepath + "/" + strconv.Itoa(timestamp.Day())
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(filepath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	filepath = filepath + "/" + strconv.Itoa(timestamp.Nanosecond())
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(filepath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func createFrontFile(employees *schemas.Employees, extensions schemas.PhotoIdExtension, path, photoDirectoryPath string) error {
	front, err := os.Create(path)
	if err != nil {
		return err
	}
	defer front.Close()

	fmt.Fprint(front, "matricula\tnome_guerra\tcargo\tlotacao\tfoto\tmostrar_foto\n")
	for _, employee := range *employees {
		fmt.Fprintf(front, "%s\t", employee.ID)
		fmt.Fprintf(front, "%s\t", employee.WarName)
		fmt.Fprintf(front, "%s\t", employee.Role)
		fmt.Fprintf(front, "%s\t", employee.Workplace)
		fmt.Fprintf(front, "%s/%s%s\t", photoDirectoryPath, employee.ID, extensions[employee.ID])
		fmt.Fprintf(front, "%s\n", "true")
	}
	return nil
}

func createBackFile(employees *schemas.Employees, path string) error {
	back, err := os.Create(path)
	if err != nil {
		return err
	}
	defer back.Close()

	fmt.Fprint(back, "matricula\tnome\tidentidade\tadmissao\n")
	for _, employee := range *employees {
		fmt.Fprintf(back, "%s\t", employee.ID)
		fmt.Fprintf(back, "%s\t", employee.Name)
		fmt.Fprintf(back, "%s\t", employee.Identification)
		fmt.Fprintf(back, "%s\n", employee.AdmissionDate)
	}
	return nil
}

func createMissingPhotoFile(employees *schemas.Employees, path string) error {
	missing, err := os.Create(path)
	if err != nil {
		return err
	}
	defer missing.Close()

	fmt.Fprint(missing, "matricula\tnome\tidentidade\tadmissao\n")
	for _, employee := range *employees {
		fmt.Fprintf(missing, "%s\t", employee.ID)
		fmt.Fprintf(missing, "%s\t", employee.Name)
		fmt.Fprintf(missing, "%s\t", employee.Identification)
		fmt.Fprintf(missing, "%s\n", employee.AdmissionDate)
	}
	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter the path to the photo directory")

	photoDirectoryPath, _ := reader.ReadString('\n')
	photoDirectoryPath = strings.TrimSpace(photoDirectoryPath)

	fmt.Println("Validating photo directory...")
	if err := validateDirectoryPath(photoDirectoryPath); err != nil {
		fmt.Println("Error validating photo directory", err)
		return
	}

	fmt.Println("Reading excel data...")
	employeesData, err := readExcel("employees.xlsx", "Cards")
	if err != nil {
		fmt.Println("Error reading excel data", err)
		return
	}

	fmt.Println("Parsing data to employees type...")
	employees := parseDataToEmployees(employeesData[1:])

	fmt.Println("Saving data in .txt files...")
	if err := employeesToTxt(&employees, photoDirectoryPath); err != nil {
		fmt.Println("Error saving employees", err)
		return
	}
	fmt.Println("Done!")
	fmt.Println("Click to exit")
	exit, _ := reader.ReadString('\n')
	fmt.Println(exit)
}
