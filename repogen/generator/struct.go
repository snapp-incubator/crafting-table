package generator

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type field struct {
	Name   string
	Type   string
	DBName string
}

type Structure struct {
	PackageName string
	DBName      string
	Name        string
	Fields      []field
}

var ErrInvalidFileName = errors.New("invalid file name")

func getFileName(path string) (string, error) {
	index := strings.Index(path, ".go")
	if index == -1 {
		return "", ErrInvalidFileName
	}

	if strings.Contains(path, "/") {
		lastBackSlash := strings.LastIndex(path, "/")
		return path[lastBackSlash+1 : index-1], nil
	}

	return path[:index-1], nil
}

func bindStruct(src string) (*Structure, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	funcFinded := false
	structure := Structure{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.Contains(line, "package") {
			structure.PackageName = strings.Split(line, " ")[1]
			continue
		}

		if strings.Contains(line, "type") && strings.Contains(line, "struct") {
			funcFinded = true
			tmp := strings.Split(line, " ")
			structure.Name = tmp[1]
			continue
		}

		if funcFinded && string(line[0]) == "}" {
			break
		}

		if funcFinded {
			for strings.Contains(line, "  ") {
				line = strings.Replace(line, "  ", " ", -1)
			}

			tmp := strings.Split(line, " ")
			if len(tmp) < 3 {
				return nil, errors.New("invalid structure")
			}

			structField := field{}
			structField.Name = tmp[0]
			structField.Type = tmp[1]

			if !strings.Contains(tmp[1], "db") {
				return nil, errors.New("db tag not found for field " + structField.Name)
			}

			index := strings.Index(tmp[1], ":")
			if index == -1 {
				return nil, errors.New("db tag is not valid for filed " + structField.Name)
			}

			structField.DBName = tmp[1][index+2 : len(tmp[1])-1]

			structure.Fields = append(structure.Fields, structField)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &structure, nil
}
