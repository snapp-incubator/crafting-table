package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/snapp-incubator/crafting-table/internal/structure"
)

func ExtractGetVariables(vars string) *[]structure.Variables {
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

	result := make([]structure.Variables, 0)
	for _, varTmp := range varSlice {
		if strings.Contains(varTmp, ",") {
			varSliceTmp := strings.Split(varTmp, ",")
			result = append(result, structure.Variables{Name: varSliceTmp})
			continue
		}

		result = append(result, structure.Variables{Name: []string{varTmp}})
	}

	return &result
}

func ExtractUpdateVariables(vars string) *[]structure.UpdateVariables {
	cleanFlag := vars
	for strings.Contains(cleanFlag, "  ") {
		cleanFlag = strings.Replace(vars, "  ", " ", -1)
	}
	cleanFlag = strings.Replace(vars, " ", "", -1)

	regex := regexp.MustCompile(`\[\(([a-zA-Z]+([a-zA-Z]*[0-9]*_*)*,)*|([a-zA-Z]+([a-zA-Z]*[0-9]*_*)*)?\),\(([a-zA-Z]+([a-zA-Z]*[0-9]*_*)*,)*|([a-zA-Z]+([a-zA-Z]*[0-9]*_*)*)?\)]`)
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

	result := make([]structure.UpdateVariables, 0)
	for _, item := range arrays {
		item = strings.Replace(item, "[", "", -1)
		item = strings.Replace(item, "]", "", -1)
		temp := strings.Split(item, "),(")

		byVariables := strings.Split(strings.Replace(temp[0], "(", "", -1), ",")
		filedVariables := strings.Split(strings.Replace(temp[1], ")", "", -1), ",")

		itemUpdateVariables := structure.UpdateVariables{
			By:     byVariables,
			Fields: filedVariables,
		}
		_ = fmt.Sprintf("%+v", itemUpdateVariables)
		result = append(result, itemUpdateVariables)
	}

	return &result
}
