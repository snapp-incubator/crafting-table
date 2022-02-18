package cmd

import (
	"errors"
	"strings"

	"github.com/n25a/repogen/repogen/generator"

	"github.com/spf13/cobra"
)

var (
	source      string
	destination string
	packageName string
	getVars     *[]generator.Variables
	get         string
	updateVars  *[]generator.Variables
	update      string
	create      bool
)

var generateCMD = &cobra.Command{
	Use:   "genarate",
	Short: "Start generating reposiory",
	Run:   generate,
}

func init() {
	generateCMD.Flags().StringVarP(&source, "source", "s", "", "Path of struct file")
	generateCMD.Flags().StringVarP(&destination, "destination", "d", "", "Path of destination to save repository file")

	err := generateCMD.MarkFlagRequired("source")
	if err != nil {
		panic(err)
	}

	err = generateCMD.MarkFlagRequired("destination")
	if err != nil {
		panic(err)
	}

	// TODO: add flag for table name

	generateCMD.Flags().StringVarP(&packageName, "package", "p", "", "Name of repository package. default is 'repository'")
	generateCMD.Flags().StringVarP(&get, "get", "g", "", "Get variables for GET functions in repository. ex: -g [ (var1,var2), (var2,var4), var3 ]")
	generateCMD.Flags().StringVarP(&update, "update", "u", "", "Get variables for UPDATE functions in repository.  ex: -g [ (var1,var2), (var2,var4), var3 ]")
	generateCMD.Flags().BoolVarP(&create, "create", "c", false, "Set to create CREATE function in repository")
}

func parseVariables(vars string) *[]generator.Variables {
	newVar := vars[0 : len(vars)-2] // remove "[" and "]"

	varSlice := strings.Split(newVar, ",")

	result := make([]generator.Variables, 0)
	for _, varTmp := range varSlice {
		if string(varTmp[0]) == "(" && string(varTmp[len(varTmp)-1]) == ")" {
			varSliceTmp := strings.Split(varTmp, ",")
			result = append(result, generator.Variables{Name: varSliceTmp})
			continue
		}

		result = append(result, generator.Variables{Name: []string{varTmp}})
	}

	return &result
}

func validateFlag(flag string) error {
	if string(flag[0]) != "[" && string(flag[len(flag)-1]) != "]" {
		return errors.New("You must set get variables in format of [ (var1,var2), (var2,var4), var3 ]")
	}

	openPrantheses := false
	for index, char := range flag {
		if openPrantheses && char == '(' {
			return errors.New("Open parentheses are not closed")
		}

		if !openPrantheses && char == ')' {
			return errors.New("Close parentheses are not opened")
		}

		if openPrantheses && char == ')' && flag[index-1] == ',' {
			return errors.New("Close parentheses must not be followed by comma")
		}

		if openPrantheses && char == ')' && flag[index+1] != ']' && flag[index+1] != ',' {
			return errors.New("Close parentheses must be followed by comma")
		}

		if char == '(' {
			openPrantheses = true
		}

		if char == ')' {
			openPrantheses = false
		}
	}
	return nil
}

func generate(cmd *cobra.Command, args []string) {
	if packageName == "" {
		packageName = "repository"
	} else {
		packageName = strings.Replace(packageName, " ", "", -1)
	}

	source = strings.Replace(source, " ", "", -1)
	destination = strings.Replace(destination, " ", "", -1)

	if get == "" && update == "" && !create {
		panic("You must set at least one flag for get, update or create")
	}

	if get != "" {
		get = strings.Replace(get, " ", "", -1)
		if err := validateFlag(get); err != nil {
			panic(err)
		}
		getVars = parseVariables(get)
	}

	if update != "" {
		update = strings.Replace(update, " ", "", -1)
		if err := validateFlag(update); err != nil {
			panic(err)
		}
		updateVars = parseVariables(update)
	}

	if err := generator.GenerateRepository(source, destination, packageName, getVars, updateVars, create); err != nil {
		panic(err)
	}
}
