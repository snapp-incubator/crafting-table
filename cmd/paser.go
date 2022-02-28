package cmd

import (
	"log"
	"regexp"
	"strings"

	"github.com/n25a/repogen/internal/generator"
)

func parseVariables(vars string) *[]generator.Variables {
	// TODO : parse variables with regex
	newVar := vars[1 : len(vars)-1] // remove "[" and "]"

	var varSlice []string
	var temp string

	for _, c := range newVar {

		if temp != "" && string(temp[0]) == "(" && string(c) == ")" {
			varSlice = append(varSlice, temp[1:])
			temp = ""
			continue
		}

		if temp != "" && string(temp[0]) == "(" && string(c) != ")" {
			temp += string(c)
			continue
		}

		if temp != "" && string(temp[0]) != "(" && string(c) == "," {
			varSlice = append(varSlice, temp)
			temp = ""
			continue
		}

		if (temp != "" && string(temp[0]) != "(" && string(c) != ",") || (temp == "" && string(c) != ",") {
			temp += string(c)
			continue
		}

	}

	if temp != "" {
		varSlice = append(varSlice, temp)
	}

	result := make([]generator.Variables, 0)
	for _, varTmp := range varSlice {
		if strings.Contains(varTmp, ",") {
			varSliceTmp := strings.Split(varTmp, ",")
			result = append(result, generator.Variables{Name: varSliceTmp})
			continue
		}

		result = append(result, generator.Variables{Name: []string{varTmp}})
	}

	return &result
}

func parseUpdateVariables(vars string) *[]generator.UpdateVariables {
	cleanFlag := vars
	for strings.Contains(cleanFlag, "  ") {
		cleanFlag = strings.Replace(vars, "  ", " ", -1)
	}
	cleanFlag = strings.Replace(vars, " ", "", -1)

	regex := regexp.MustCompile(`\[\(([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*,)*|([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*)?\),\(([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*,)*|([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*)?\)\]`)
	matches := regex.FindAllString(cleanFlag, -1)

	arrays := make([]string, 0)
	tmp := ""
	for _, m := range matches {
		if strings.Contains(m, "]") {
			tmp += m
			arrays = append(arrays, tmp)
			tmp = ""
			continue
		}
		tmp += m
	}

	firstParenthesesRegex := regexp.MustCompile(`\(([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*,)+|([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*)?\),`)
	secondParenthesesRegex := regexp.MustCompile(`\(([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*,)+|([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*)?\)`)
	//variablesRegex := regexp.MustCompile(`^[a-zA-Z]+[0-9]*$`)
	result := make([]generator.UpdateVariables, 0)
	for _, item := range arrays {
		byVariablesTmp := firstParenthesesRegex.FindAllString(item, -1)
		byVariables := ""
		for _, v := range byVariablesTmp {
			byVariables += v
		}
		byVariables = byVariables[:len(byVariables)-1]

		filedVariablesTmp := secondParenthesesRegex.FindAllString(item, -1)
		filedVariables := ""
		for _, v := range filedVariablesTmp {
			filedVariables += v
		}
		log.Fatal(byVariables, ":::::", filedVariablesTmp)
		//itemUpdateVariables := generator.UpdateVariables{
		//	By:     variablesRegex.FindAllString(byVariables, -1),
		//	Fields: variablesRegex.FindAllString(filedVariables, 1),
		//}
		//_ = fmt.Sprintf("%+v", itemUpdateVariables)
		//result = append(result, itemUpdateVariables)
	}

	return &result
}
