package cmd

import (
	"log"
	"strings"

	"github.com/snapp-incubator/crafting-table/internal/parser"
	"github.com/snapp-incubator/crafting-table/internal/repository"
	"github.com/snapp-incubator/crafting-table/internal/structure"

	"github.com/spf13/cobra"
)

var (
	source      string
	destination string
	packageName string
	structName  string
	get         string
	update      string
	create      bool
	test        bool
)

var generateCMD = &cobra.Command{
	Use:   "generate",
	Short: "Start generating repository",
	Run:   generate,
}

func init() {
	generateCMD.Flags().StringVarP(&source, "source", "s", "", "Path of struct file")
	generateCMD.Flags().StringVarP(&destination, "destination", "d", "", "Path of destination to save repository file")

	// TODO: add flag for table name
	generateCMD.Flags().StringVarP(&packageName, "package", "p", "", "Name of repository package. default is 'repository'")
	generateCMD.Flags().StringVarP(&get, "get", "g", "", "Get variables for GET functions in repository. ex: -g [ (var1,var2), (var2,var4), var3 ]")
	generateCMD.Flags().StringVarP(&update, "update", "u", "", "Get variables for UPDATE functions in repository.  ex: -g [ [(byPar1,byPar2,...), (field1, field2)], ... ]")
	generateCMD.Flags().StringVarP(&structName, "struct-name", "n", "", "find struct with struct name in source file")
	generateCMD.Flags().BoolVarP(&create, "create", "c", false, "Set to create CREATE function in repository")
	generateCMD.Flags().BoolVarP(&test, "test", "t", false, "generate automatically tests for created repository")
}

func generate(_ *cobra.Command, _ []string) {
	// generate repository with cli
	if packageName == "" {
		packageName = "repository"
	} else {
		packageName = strings.Replace(packageName, " ", "", -1)
	}

	source = strings.Replace(source, " ", "", -1)
	destination = strings.Replace(destination, " ", "", -1)

	if get == "" && update == "" && !create {
		log.Fatal("you must set at least one flag for get, update or create")
	}

	var getVars *[]structure.GetVariable
	if get != "" {
		for strings.Contains(get, " ") {
			get = strings.Replace(get, " ", "", -1)
		}
		get = strings.Replace(get, " ", "", -1)
		if err := parser.ValidateGetFlag(get); err != nil {
			log.Fatal(err)
		}
		getVars = parser.ExtractGetVariables(get)
	}

	var updateVars *[]structure.UpdateVariables
	if update != "" {
		for strings.Contains(update, " ") {
			update = strings.Replace(update, " ", "", -1)
		}
		update = strings.Replace(update, " ", "", -1)
		if err := parser.ValidateUpdateFlag(update); err != nil {
			log.Fatal(err)
		}
		updateVars = parser.ExtractUpdateVariables(update)
	}

	if err := repository.Generate(source, destination, packageName, structName, getVars, updateVars, create, test); err != nil {
		log.Fatal(err)
	}
}
