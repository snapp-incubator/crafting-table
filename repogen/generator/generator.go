package generator

import "log"

type Variables struct {
	Name []string
}

var createSyntax = ""
var updateSyntax = ""
var getSyntax = ""
var functions = []string{}

func GenerateRepository(source, destination, packageName string, getVars, updateVars *[]Variables, create bool) error {
	fileName, err := getFileName(source)
	if err != nil {
		log.Println("Error in getFileName: ", err)
		return err
	}

	structure, err := bindStruct(source)
	if err != nil {
		log.Println("Error in bindStruct: ", err)
		return err
	}

	var functions = []string{}
	var function string

	if create {
		createSyntax, function, err = createFunctionRepository(structure)
		if err != nil {
			log.Println("Error in createFunctionRepository: ", err)
			return err
		}
		functions = append(functions, function)
	}

	if getVars != nil {
		getSyntax, function, err = getFunctionCreator(getVars)
		if err != nil {
			log.Println("Error in getFunctionCreator: ", err)
			return err
		}
		functions = append(functions, function...)
	}

	//if updateVars != nil {
	//	updateSyntax,err = updateFunctionCreator(fileContent, updateVars)
	//	if err != nil {
	//		log.Println("Error in updateFunctionCreator: ", err)
	//		return err
	//	}
	//functions = append(functions, funcs...)
	//}
	//
	//interfaceSyntax, err := interfaceSyntaxCreator(structure, getVars, updateVars, create)
	//if err != nil {
	//	log.Println("Error in interfaceSyntaxCreator: ", err)
	//	return err
	//}

	//fileContent := createTemplate(fileName, packageName, interfaceSyntax, structure)

	//err = writeFile(fileContent, destination)
	//if err != nil {
	//	log.Println("Error in writeFile: ", err)
	//	return err
	//}

	//err = linter(distanationFileName)
	//if err != nil {
	//	log.Println("Error in linter: ", err)
	//	return err
	//}

	return nil
}
