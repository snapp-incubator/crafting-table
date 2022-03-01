package cmd

import (
	"log"
	"strings"

	"github.com/n25a/repogen/internal/generator"

	"github.com/spf13/cobra"
)

var (
	source      string
	destination string
	packageName string
	getVars     *[]generator.Variables
	get         string
	updateVars  *[]generator.UpdateVariables
	update      string
	create      bool
)

var generateCMD = &cobra.Command{
	Use:   "generate",
	Short: "Start generating repository",
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
	generateCMD.Flags().StringVarP(&update, "update", "u", "", "Get variables for UPDATE functions in repository.  ex: -g [ [(byPar1,byPar2,...), (field1, field2)], ... ]")
	generateCMD.Flags().BoolVarP(&create, "create", "c", false, "Set to create CREATE function in repository")
}

func generate(_ *cobra.Command, _ []string) {
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
		for strings.Contains(get, " ") {
			get = strings.Replace(get, " ", "", -1)
		}
		get = strings.Replace(get, " ", "", -1)
		if err := validateFlag(get); err != nil {
			panic(err)
		}
		getVars = parseVariables(get)
	}

	if update != "" {
		for strings.Contains(update, " ") {
			update = strings.Replace(update, " ", "", -1)
		}
		update = strings.Replace(update, " ", "", -1)
		if err := validateUpdateFlag(update); err != nil {
			panic(err)
		}
		updateVars = parseUpdateVariables(update)
	}

	if err := generator.GenerateRepository(source, destination, packageName, getVars, updateVars, create); err != nil {
		log.Println(err)
	}
}
