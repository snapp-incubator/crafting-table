package structure

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
)

type CreateVariables struct {
	Enable              bool   `yaml:"enable"`
	CostumeFunctionName string `yaml:"costume_function_name"`
}

type GetVariable struct {
	Conditions          []string `yaml:"conditions"`
	CostumeFields       []string `yaml:"costume_fields"`
	CostumeFunctionName string   `yaml:"costume_function_name"`
}

type UpdateVariables struct {
	Conditions          []string `yaml:"conditions"`
	Fields              []string `yaml:"fields"`
	CostumeFunctionName string   `yaml:"costume_function_name"`
}

type JoinVariables struct {
	JoinStructPath      string `yaml:"join_struct_path"`
	JoinStructName      string `yaml:"join_struct_name"`
	JoinTableName       string `yaml:"join_table_name"`
	JoinFieldAs         string `yaml:"join_field_as"`
	CostumeFunctionName string `yaml:"costume_function_name"`
}

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

func (s *Structure) GetDBFields(prefix string) string {
	result := "\""
	tmp := ""
	for _, field := range s.Fields {
		if len(tmp) > 80 {
			result += tmp[:len(tmp)-2] + ", \"+\n\t\""
			tmp = ""
		}
		tmp += prefix + field.DBFlag + ", "
	}

	if tmp != "" {
		result += tmp[:len(tmp)-2] + "\""
	}

	return result
}

func (s *Structure) GetDBFieldsInQuotation() string {
	result := ""
	for _, field := range s.Fields {
		result += "\"" + field.DBFlag + "\",\n\t\t\t"
	}

	return result
}

func (s *Structure) GetVariableFields(prefix string) string {
	result := ""
	for _, field := range s.Fields {
		result += prefix + field.Name + ",\n\t\t\t"
	}

	return result
}
