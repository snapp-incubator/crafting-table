package parser

//func ExtractGetVariables(vars string) *[]structure.GetVariable {
//	// TODO : parse variables with regex
//
//	newVar := vars[1 : len(vars)-1] // remove "[" and "]"
//
//	var varSlice []string
//	var temp string
//
//	for _, c := range newVar {
//
//		if temp != "" && string(temp[0]) == "(" && string(c) == ")" {
//			varSlice = append(varSlice, temp[1:])
//			temp = ""
//			continue
//		}
//
//		if temp != "" && string(temp[0]) == "(" && string(c) != ")" {
//			temp += string(c)
//			continue
//		}
//
//		if temp != "" && string(temp[0]) != "(" && string(c) == "," {
//			varSlice = append(varSlice, temp)
//			temp = ""
//			continue
//		}
//
//		if (temp != "" && string(temp[0]) != "(" && string(c) != ",") || (temp == "" && string(c) != ",") {
//			temp += string(c)
//			continue
//		}
//
//	}
//
//	if temp != "" {
//		varSlice = append(varSlice, temp)
//	}
//
//	result := make([]structure.GetVariable, 0)
//	for _, varTmp := range varSlice {
//		if strings.Contains(varTmp, ",") {
//			varSliceTmp := strings.Split(varTmp, ",")
//			result = append(result, structure.GetVariable{Conditions: varSliceTmp})
//			continue
//		}
//
//		result = append(result, structure.GetVariable{Conditions: []string{varTmp}})
//	}
//
//	return &result
//}
//
//func ExtractUpdateVariables(vars string) *[]structure.UpdateVariables {
//	cleanFlag := vars
//	for strings.Contains(cleanFlag, "  ") {
//		cleanFlag = strings.Replace(vars, "  ", " ", -1)
//	}
//	cleanFlag = strings.Replace(vars, " ", "", -1)
//
//	regex := regexp.MustCompile(`\[\(([a-zA-Z]+([a-zA-Z]*[0-9]*_*)*,)*|([a-zA-Z]+([a-zA-Z]*[0-9]*_*)*)?\),\(([a-zA-Z]+([a-zA-Z]*[0-9]*_*)*,)*|([a-zA-Z]+([a-zA-Z]*[0-9]*_*)*)?\)]`)
//	matches := regex.FindAllString(cleanFlag, -1)
//
//	arrays := make([]string, 0)
//	tmp := ""
//	for _, m := range matches {
//		if strings.Contains(m, "]") {
//			tmp += m
//			arrays = append(arrays, tmp)
//			tmp = ""
//			continue
//		}
//		tmp += m
//	}
//
//	result := make([]structure.UpdateVariables, 0)
//	for _, item := range arrays {
//		item = strings.Replace(item, "[", "", -1)
//		item = strings.Replace(item, "]", "", -1)
//		temp := strings.Split(item, "),(")
//
//		byVariables := strings.Split(strings.Replace(temp[0], "(", "", -1), ",")
//		filedVariables := strings.Split(strings.Replace(temp[1], ")", "", -1), ",")
//
//		itemUpdateVariables := structure.UpdateVariables{
//			Conditions: byVariables,
//			Fields:     filedVariables,
//		}
//		result = append(result, itemUpdateVariables)
//	}
//
//	return &result
//}
//
//func ExtractJoinVariables(join string) *[]structure.JoinVariables {
//	cleanFlag := join
//	// remove spaces
//	for strings.Contains(cleanFlag, "  ") {
//		cleanFlag = strings.Replace(join, "  ", " ", -1)
//	}
//	cleanFlag = strings.Replace(join, " ", "", -1)
//	cleanFlag = cleanFlag[1 : len(cleanFlag)-1] // remove brackets
//
//	tmps := strings.Split(cleanFlag, "],[")
//	for i, tmp := range tmps {
//		tmps[i] = strings.Replace(tmp, "[", "", -1)
//	}
//
//	results := make([]structure.JoinVariables, 0)
//	for _, tmp := range tmps {
//		parenthesis := strings.Split(tmp, "),(")
//		joinList := &structure.JoinVariables{}
//		for i, parenthesisTmp := range parenthesis {
//			parenthesis[i] = strings.Replace(parenthesisTmp, "(", "", -1)
//			parenthesis[i] = strings.Replace(parenthesisTmp, ")", "", -1)
//			vars := strings.Split(parenthesis[i], ",")
//			joinList.Fields = append(joinList.Fields, structure.JoinField{
//				JoinStructPath: vars[0],
//				JoinStructName: vars[1],
//				JoinFieldAs:    vars[2],
//				JoinOn:         vars[3],
//				JoinType:       vars[4],
//			})
//		}
//
//		results = append(results, *joinList)
//	}
//
//	return &results
//}
