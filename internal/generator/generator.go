package generator

import (
	"errors"
	"fmt"

	"github.com/n25a/repogen/internal/structure"
)

var createSyntax = ""
var updateSyntax = ""
var getSyntax = ""

func GenerateRepository(source, destination, packageName string, getVars *[]structure.Variables, updateVars *[]structure.UpdateVariables, create bool) error {
	s, err := structure.BindStruct(source)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error in bindStruct: %s", err.Error()))
		return err
	}

	var functions []string
	var function []string

	if create {
		var function string
		createSyntax, function, err = createFunctionRepository(s)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in createFunctionRepository: %s", err.Error()))
			return err
		}
		functions = append(functions, function)
	}

	if getVars != nil {
		getSyntax, function, err = getFunctionCreator(s, getVars)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in getFunctionCreator: %s", err.Error()))
			return err
		}
		functions = append(functions, function...)
	}

	if updateVars != nil {
		updateSyntax, function, err = updateFunctionCreator(s, updateVars)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in updateFunctionCreator: %s", err.Error()))
			return err
		}
		functions = append(functions, function...)
	}

	interfaceSyntax := interfaceSyntaxCreator(s, functions)

	fileContent := createTemplate(s, packageName, interfaceSyntax,
		createSyntax, updateSyntax, getSyntax)

	err = writeFile(fileContent, destination)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error in writeFile: %s", err.Error()))
		return err
	}

	err = linter(destination)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error in linter: %s", err.Error()))
		return err
	}

	return nil
}
