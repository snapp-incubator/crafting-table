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
	PackageName       string
	DBName            string
	Name              string
	Fields            []field
	FieldNameToType   map[string]string
	FieldDBNameToName map[string]string
	FieldNameToDBName map[string]string
}

func bindStruct(src string) (*Structure, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	scanner := bufio.NewScanner(f)

	funcFinded := false
	structure := new(Structure)

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

			if !strings.Contains(tmp[2], "db") {
				return nil, errors.New("db tag not found for field " + structField.Name)
			}

			index := strings.Index(tmp[2], ":")
			if index == -1 {
				return nil, errors.New("db tag is not valid for filed " + structField.Name)
			}

			structField.DBName = tmp[2][index+2 : len(tmp[2])-2]
			structure.Fields = append(structure.Fields, structField)

			if structure.FieldDBNameToName == nil {
				structure.FieldDBNameToName = make(map[string]string)
			}
			if structure.FieldNameToDBName == nil {
				structure.FieldNameToDBName = make(map[string]string)
			}
			if structure.FieldNameToType == nil {
				structure.FieldNameToType = make(map[string]string)
			}

			structure.FieldDBNameToName[structField.DBName] = structField.Name
			structure.FieldNameToDBName[structField.Name] = structField.DBName
			structure.FieldNameToType[structField.Name] = structField.Type
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return structure, nil
}

func (s *Structure) GetDBFields(prefix string) string {
	result := "\""
	tmp := ""
	for _, field := range s.Fields {
		if len(tmp) > 80 {
			result += tmp[:len(tmp)-2] + "\"+\n\t\""
			tmp = ""
		}
		tmp += prefix + field.DBName + ", "
	}

	if tmp != "" {
		result += tmp[:len(tmp)-2] + "\""
	}

	return result
}
