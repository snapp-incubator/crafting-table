package generator

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/iancoleman/strcase"
)

func getConditions(v []string, structure *Structure) string {
	var conditions []string

	for _, value := range v {
		res := structure.FieldDBNameToName[value]
		if res == "" {
			panic(fmt.Sprintf("%s is not valid db_field", value))
		}
	}

	for _, value := range v {
		conditions = append(conditions, fmt.Sprintf("%s = ?", value))
	}
	return "WHERE " + strings.Join(conditions, " AND ")
}

func getFunctionVars(v []string, structure *Structure) string {
	for _, value := range v {
		res := structure.FieldDBNameToName[value]
		if res == "" {
			panic(fmt.Sprintf("%s is not valid db_field", value))
		}
	}

	res := ""
	for _, value := range v {
		res += fmt.Sprintf("%s %s, ",
			strcase.ToLowerCamel(structure.FieldDBNameToName[value]),
			structure.FieldNameToType[structure.FieldDBNameToName[value]])
	}

	return res[:len(res)-2]
}

func getUpdateVariables(v []string, structure *Structure) string {
	for _, value := range v {
		res := structure.FieldDBNameToName[value]
		if res == "" {
			panic(fmt.Sprintf("%s is not valid db_field", value))
		}
	}

	res := ""
	for _, value := range v {
		res += fmt.Sprintf("%s, ",
			strcase.ToLowerCamel(structure.FieldDBNameToName[value]))
	}

	return res[:len(res)-2]
}

func interfaceSyntaxCreator(structure *Structure, functions []string) string {
	syntax := fmt.Sprintf(
		"type %s interface {",
		structure.Name,
	)

	for _, function := range functions {
		syntax += fmt.Sprintf("\n\t%s", function)
	}
	syntax += "\n}"

	return syntax
}

func linter(dst string) error {
	cmd := exec.Command("gofmt", "-s", "-w", dst)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("goimports", dst)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func writeFile(content, dst string) error {
	f, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_, err = f.WriteString(content)

	if err != nil {
		return err
	}

	return nil
}

func setKeys(fields []field) string {
	result := "\""
	tmp := ""
	for _, field := range fields {
		if len(tmp) > 80 {
			result += tmp[:len(tmp)-2] + "\"+\n\t\""
			tmp = ""
		}
		tmp += field.DBName + " = :" + field.DBName + ", "
	}

	if tmp != "" {
		result += tmp[:len(tmp)-2] + "\""
	}

	return result
}

func setKeysWithQuestion(vars []string) string {
	result := ""
	tmp := ""
	for _, varName := range vars {
		if len(tmp) > 80 {
			result += tmp[:len(tmp)-2] + "\"+\n\t\""
			tmp = ""
		}
		tmp += varName + " = ?, "
	}

	if tmp != "" {
		result += tmp[:len(tmp)-2] + "\""
	}

	return result
}

func execContextVariables(vars UpdateVariables, structure *Structure) string {
	result := ""
	for _, variable := range vars.Fields {
		result += fmt.Sprintf("%s, ", strcase.ToLowerCamel(structure.FieldDBNameToName[variable]))
	}

	for _, variable := range vars.By {
		result += fmt.Sprintf("%s, ", strcase.ToLowerCamel(structure.FieldDBNameToName[variable]))
	}

	return result[:len(result)-2]
}
