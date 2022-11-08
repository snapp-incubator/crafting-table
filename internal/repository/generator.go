package repository

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/snapp-incubator/crafting-table/internal/structure"
)

func Generate(source, destination, packageName, structName string, getVars *[]structure.GetVariable,
	updateVars *[]structure.UpdateVariables, joinVars *[]structure.JoinVariables,
	aggregateVars *[]structure.AggregateField, create, test bool) error {
	createSyntax := ""
	updateSyntax := ""
	getSyntax := ""
	joinSyntax := ""
	aggregateSyntax := ""

	createTestSyntax := ""
	updateTestSyntax := ""
	getTestSyntax := ""
	joinTestSyntax := ""
	aggregateTestSyntax := ""

	var testDestination string
	if test {
		testDestination = destination[:len(destination)-3] + "_test.go"
	}

	s, err := structure.BindStruct(source, structName)
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
			err = errors.New(fmt.Sprintf("Error in createFunction: %s", err.Error()))
			return err
		}
		signatures = append(signatures, signature)

		if test {
			createTestSyntax = createTestFunction(s)
		}
	}

	if getVars != nil && len(*getVars) > 0 {
		getSyntax, signatureList, err = getFunction(s, getVars)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in getFunction: %s", err.Error()))
			return err
		}
		signatures = append(signatures, signatureList...)

		if test {
			getTestSyntax = getTestFunction(s, getVars)
		}
	}

	if updateVars != nil && len(*updateVars) > 0 {
		updateSyntax, signatureList, err = updateFunction(s, updateVars)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in updateFunction: %s", err.Error()))
			return err
		}
		signatures = append(signatures, signatureList...)

		if test {
			updateTestSyntax = updateTestFunction(s, updateVars)

		}
	}

	if joinVars != nil {
		syntax, signatureList, err := joinFunction(s, joinVars)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in joinFunction: %s", err.Error()))
			return err
		}

		joinSyntax += syntax
		signatures = append(signatures, signatureList...)

		if test {
			joinTestSyntax += joinTestFunction(s, joinVars)
		}
	}

	if aggregateVars != nil && len(*aggregateVars) > 0 {
		aggregateSyntax, signatureList, err = aggregateFunction(s, aggregateVars)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in aggregateSyntax: %s", err.Error()))
			return err
		}
		signatures = append(signatures, signatureList...)

		//if test {
		//	aggregateTestSyntax = aggregateTestFunction(s, getVars)
		//}
	}

	interfaceSyntax := interfaceCreator(s, signatures)

	fileContent := createTemplate(s, packageName, interfaceSyntax,
		createSyntax, updateSyntax, getSyntax, joinSyntax, aggregateSyntax)

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

	if test {
		testFileContent := createTestTemplate(s, packageName, createTestSyntax, updateTestSyntax,
			getTestSyntax, joinTestSyntax, aggregateTestSyntax)

		err = exportRepository(testFileContent, testDestination)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in writeTestFile: %s", err.Error()))
			return err
		}

		err = linter(testDestination)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in linter: %s", err.Error()))
			return err
		}
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
