package repository

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/snapp-incubator/crafting-table/internal/structure"
)

var createSyntax = ""
var updateSyntax = ""
var getSyntax = ""

func Generate(source, destination, packageName string, getVars *[]structure.Variables, updateVars *[]structure.UpdateVariables, create bool) error {
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

	interfaceSyntax := interfaceCreator(s, functions)

	fileContent := createTemplate(s, packageName, interfaceSyntax,
		createSyntax, updateSyntax, getSyntax)

	err = exportRepository(fileContent, destination)
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

func interfaceCreator(structure *structure.Structure, functions []string) string {
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
	cmd := exec.Command("goimports", "-w", dst)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		err = errors.New(fmt.Sprintf("Error in goimports: %s", err.Error()))
		return err
	}

	cmd = exec.Command("gofmt", "-s", dst)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		err = errors.New(fmt.Sprintf("Error in gofmt: %s", err.Error()))
		return err
	}

	return nil
}

func exportRepository(content, dst string) error {
	f, err := os.Create(dst)

	if err != nil {
		err = errors.New(fmt.Sprintf("Error in creating file: %s", err.Error()))
		return err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_, err = f.WriteString(content)

	if err != nil {
		err = errors.New(fmt.Sprintf("Error in writing file: %s", err.Error()))
		return err
	}

	return nil
}
