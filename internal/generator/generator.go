package generator

import (
	"log"
)

type Variables struct {
	Name []string
}

type UpdateVariables struct {
	By     []string
	Fields []string
}

var createSyntax = ""
var updateSyntax = ""
var getSyntax = ""

func GenerateRepository(source, destination, packageName string, getVars *[]Variables, updateVars *[]UpdateVariables, create bool) error {
	log.Println("Generating repository")
	structure, err := bindStruct(source)
	if err != nil {
		log.Println("Error in bindStruct: ", err)
		return err
	}

	var functions = []string{}
	var function []string

	if create {
		var function string
		createSyntax, function, err = createFunctionRepository(structure)
		if err != nil {
			log.Println("Error in createFunctionRepository: ", err)
			return err
		}
		functions = append(functions, function)
	}

	if getVars != nil {
		getSyntax, function, err = getFunctionCreator(structure, getVars)
		if err != nil {
			log.Println("Error in getFunctionCreator: ", err)
			return err
		}
		functions = append(functions, function...)
	}

	if updateVars != nil {
		updateSyntax, function, err = updateFunctionCreator(structure, updateVars)
		if err != nil {
			log.Println("Error in updateFunctionCreator: ", err)
			return err
		}
		functions = append(functions, function...)
	}

	interfaceSyntax := interfaceSyntaxCreator(structure, functions)

	fileContent := createTemplate(structure, packageName, interfaceSyntax,
		createSyntax, updateSyntax, getSyntax)

	err = writeFile(fileContent, destination)
	if err != nil {
		log.Println("Error in writeFile: ", err)
		return err
	}

	err = linter(destination)
	if err != nil {
		log.Println("Error in linter: ", err)
		return err
	}

	return nil
}
