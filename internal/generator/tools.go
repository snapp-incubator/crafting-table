package generator

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/snapp-incubator/crafting-table/internal/structure"
)

func interfaceSyntaxCreator(structure *structure.Structure, functions []string) string {
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
		fmt.Println("goooooooooooooooooooooooooooooooooooooooooooooooooooz")
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

func writeFile(content, dst string) error {
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
