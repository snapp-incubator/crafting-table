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

	var signatures []string
	var signatureList []string

	if create {
		var signature string
		createSyntax, signature, err = createFunction(s)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in createFunctionRepository: %s", err.Error()))
			return err
		}
		signatures = append(signatures, signature)
	}

	if getVars != nil {
		getSyntax, signatureList, err = getFunction(s, getVars)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in getFunctionCreator: %s", err.Error()))
			return err
		}
		signatures = append(signatures, signatureList...)
	}

	if updateVars != nil {
		updateSyntax, signatureList, err = updateFunction(s, updateVars)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in updateFunctionCreator: %s", err.Error()))
			return err
		}
		signatures = append(signatures, signatureList...)
	}

	interfaceSyntax := interfaceCreator(s, signatures)

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

func interfaceCreator(structure *structure.Structure, signatures []string) string {
	syntax := fmt.Sprintf(
		"type %s interface {",
		structure.Name,
	)

	for _, signature := range signatures {
		syntax += fmt.Sprintf("\n\t%s", signature)
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
