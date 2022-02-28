package cmd

import (
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
	regex := regexp.MustCompile(`^\[\((([a-zA-Z]+[0-9]*,)*([a-zA-Z]+[0-9]*)?)\),\((([a-zA-Z]+[0-9]*,)*([a-zA-Z]+[0-9]*)?)\)\]$`)
	arrays := regex.FindAllString(vars, -1)

	firstParenthesesRegex := regexp.MustCompile(`^\((([a-zA-Z]+[0-9]*,)*([a-zA-Z]+[0-9]*)?)\),$`)
	secondParenthesesRegex := regexp.MustCompile(`^\((([a-zA-Z]+[0-9]*,)*([a-zA-Z]+[0-9]*)?)\)$`)
	variablesRegex := regexp.MustCompile(`^[a-zA-Z]+[0-9]*$`)
	result := make([]generator.UpdateVariables, 0)
	for _, item := range arrays {
		byVariables := firstParenthesesRegex.FindString(item)
		filedVariables := secondParenthesesRegex.FindString(item)

		itemUpdateVariables := generator.UpdateVariables{
			By:     variablesRegex.FindAllString(byVariables, -1),
			Fields: variablesRegex.FindAllString(filedVariables, 1),
		}
		result = append(result, itemUpdateVariables)
	}

	return &result
}
