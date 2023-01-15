package build

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	internalStruct "github.com/snapp-incubator/crafting-table/internal/structure"
)

func Generate(repo Repo) error {
	//testDestination := repo.Destination[:len(repo.Destination)-3] + "_test.go"

	s, err := internalStruct.BindStruct(repo.Source, repo.StructName)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error in bindStruct: %s", err.Error()))
		return err
	}

	tableName := s.TableName
	if repo.TableName != "" {
		tableName = repo.TableName
	}

	var signatureList []string
	var functionList []string

	// Select
	for _, r := range repo.Select {
		if r.Type == SelectTypeGet {
			function, signature := BuildGetFunction(
				s,
				repo.Dialect,
				tableName,
				r.Fields,
				r.WhereConditions,
				r.AggregateFields,
				&r.OrderBy,
				&r.OrderType,
				&r.Limit,
				r.GroupBy,
				r.JoinFields,
				r.FunctionName,
			)
			functionList = append(functionList, function)
			signatureList = append(signatureList, signature)
		} else if r.Type == SelectTypeSelect {
			function, signature := BuildSelectFunction(
				s,
				repo.Dialect,
				tableName,
				r.Fields,
				r.WhereConditions,
				r.AggregateFields,
				&r.OrderBy,
				&r.OrderType,
				&r.Limit,
				r.GroupBy,
				r.JoinFields,
				r.FunctionName,
			)
			functionList = append(functionList, function)
			signatureList = append(signatureList, signature)
		}
	}

	// Update
	for _, r := range repo.Update {
		function, signature := BuildUpdateFunction(
			s,
			repo.Dialect,
			tableName,
			r.Fields,
			r.WhereConditions,
			r.FunctionName,
		)
		functionList = append(functionList, function)
		signatureList = append(signatureList, signature)
	}

	repoTemplate := BuildRepository(signatureList, functionList, repo.PackageName, s.TableName, s.Name)

	err = exportRepository(repoTemplate, repo.Destination)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error in writeFile: %s", err.Error()))
		return err
	}

	err = linter(repo.Destination)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error in linter: %s", err.Error()))
		return err
	}

	// TODO: add tests

	return nil
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
