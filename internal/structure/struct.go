package structure

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
)

type Field struct {
	Name   string
	Type   string
	DBFlag string
}

type Structure struct {
	PackageName          string
	TableName            string
	Name                 string
	Fields               []Field
	FieldMapNameToType   map[string]string
	FieldMapDBFlagToName map[string]string
	FieldMapNameToDBFlag map[string]string
}

func BindStruct(src, structName string) (*Structure, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	scanner := bufio.NewScanner(f)

	funcFound := false
	structure := new(Structure)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.Contains(line, "package") {
			structure.PackageName = strings.Split(line, " ")[1]
			continue
		}

		if strings.Contains(line, "type") && strings.Contains(line, "struct") {
			if structName != "" && !strings.Contains(line, structName) {
				continue
			}
			funcFound = true
			tmp := strings.Split(line, " ")
			structure.Name = tmp[1]
			structure.TableName = strcase.ToSnake(tmp[1])
			continue
		}

		if funcFound && string(line[0]) == "}" {
			break
		}

		if funcFound {
			for strings.Contains(line, "  ") {
				line = strings.Replace(line, "  ", " ", -1)
			}
			tmp := strings.Split(line, " ")
			if len(tmp) < 3 {
				return nil, errors.New("invalid structure")
			}

			structField := Field{}
			structField.Name = tmp[0]
			structField.Type = tmp[1]

			if !strings.Contains(tmp[2], "db") {
				return nil, errors.New("db tag not found for field " + structField.Name)
			}

			index := strings.Index(tmp[2], ":")
			if index == -1 {
				return nil, errors.New("db tag is not valid for filed " + structField.Name)
			}

			structField.DBFlag = tmp[2][index+2 : len(tmp[2])-2]
			structure.Fields = append(structure.Fields, structField)

			if structure.FieldMapDBFlagToName == nil {
				structure.FieldMapDBFlagToName = make(map[string]string)
			}
			if structure.FieldMapNameToDBFlag == nil {
				structure.FieldMapNameToDBFlag = make(map[string]string)
			}
			if structure.FieldMapNameToType == nil {
				structure.FieldMapNameToType = make(map[string]string)
			}

			structure.FieldMapDBFlagToName[structField.DBFlag] = structField.Name
			structure.FieldMapNameToDBFlag[structField.Name] = structField.DBFlag
			structure.FieldMapNameToType[structField.Name] = structField.Type
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return structure, nil
}
