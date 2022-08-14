package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/snapp-incubator/crafting-table/internal/parser"
	"github.com/snapp-incubator/crafting-table/internal/repository"
	"github.com/snapp-incubator/crafting-table/internal/structure"
	"gopkg.in/yaml.v3"

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
	ymlPath     string
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
	generateCMD.Flags().StringVarP(&update, "update", "u", "", "Get variables for UPDATE functions in repository.  ex: -u [ [(byPar1,byPar2,...), (field1, field2)], ... ]")
	generateCMD.Flags().StringVarP(&ymlPath, "yml-path", "y", "", "generate automatically repositories from yml file")
	generateCMD.Flags().StringVarP(&structName, "struct-name", "n", "", "find struct with struct name in source file")
	generateCMD.Flags().BoolVarP(&create, "create", "c", false, "Set to create CREATE function in repository")
	generateCMD.Flags().BoolVarP(&test, "test", "t", false, "generate automatically tests for created repository")
}

func generate(_ *cobra.Command, _ []string) {
	var repositories repository.Repositories

	if ymlPath != "" {
		file, err := os.Open(ymlPath)
		if err != nil {
			panic(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)

		d := yaml.NewDecoder(file)

		if err := d.Decode(&repositories); err != nil {
			panic(err)
		}
	} else {
		if packageName == "" {
			packageName = "repository"
		} else {
			packageName = strings.Replace(packageName, " ", "", -1)
		}

		repositories = repository.Repositories{
			Repositories: []repository.Repository{
				{
					Source:      source,
					Destination: destination,
					PackageName: packageName,
					StructName:  structName,
					Get:         get,
					Update:      update,
					Create:      create,
					Test:        test,
				},
			},
		}
	}

	for _, params := range repositories.Repositories {
		generateRepository(params)
	}
}

func generateRepository(params repository.Repository) {
	if params.PackageName == "" {
		packageName = "repository"
	}

	source = strings.Replace(params.Source, " ", "", -1)
	destination = strings.Replace(params.Destination, " ", "", -1)

	if params.Get == "" && params.Update == "" && !params.Create {
		log.Fatal("you must set at least one flag for get, update or create")
	}

	var getVars *[]structure.Variables
	if params.Get != "" {
		for strings.Contains(params.Get, " ") {
			params.Get = strings.Replace(params.Get, " ", "", -1)
		}
		params.Get = strings.Replace(params.Get, " ", "", -1)
		if err := parser.ValidateGetFlag(params.Get); err != nil {
			log.Fatal(err)
		}
		getVars = parser.ExtractGetVariables(params.Get)
	}

	var updateVars *[]structure.UpdateVariables
	if params.Update != "" {
		for strings.Contains(params.Update, " ") {
			params.Update = strings.Replace(params.Update, " ", "", -1)
		}
		params.Update = strings.Replace(params.Update, " ", "", -1)
		if err := parser.ValidateUpdateFlag(params.Update); err != nil {
			log.Fatal(err)
		}
		updateVars = parser.ExtractUpdateVariables(params.Update)
	}

	if err := repository.Generate(source, destination, packageName, params.StructName, getVars, updateVars, params.Create, params.Test); err != nil {
		log.Fatal(err)
	}
}
