package cmd

import "github.com/spf13/cobra"

type variables struct {
	Name []string
}

var (
	source      string
	destination string
	packageName string
	getVars     []variables
	get         string
	updateVars  []variables
	update      string
	create      bool
)

var genarateCMD = &cobra.Command{
	Use:   "genarate",
	Short: "Start generating reposiory",
	Run:   genarate,
}

func init() {
	genarateCMD.Flags().StringVarP(&source, "source", "s", "", "Path of struct file")
	genarateCMD.MarkFlagRequired("source")
	genarateCMD.Flags().StringVarP(&destination, "destination", "d", "", "Path of destination to save repository file")
	genarateCMD.MarkFlagRequired("destination")

	// TODO: add flag for table name

	genarateCMD.Flags().StringVarP(&packageName, "package", "p", "", "Name of repository package. default is 'repository'")
	genarateCMD.Flags().StringVarP(&get, "get", "g", "", "Get variables for GET functions in repository. ex: -g [ (var1,var2), (var2,var4), var3 ]")
	genarateCMD.Flags().StringVarP(&update, "update", "u", "", "Get variables for UPDATE functions in repository.  ex: -g [ (var1,var2), (var2,var4), var3 ]")
	genarateCMD.Flags().BoolVarP(&create, "create", "c", false, "Set to create CREATE function in repository")
}

func genarate(cmd *cobra.Command, args []string) {
	if packageName == "" {
		packageName = "repository"
	}

	if get == "" && update == "" && !create {
		panic("You must set at least one flag for get, update or create")
	}

	if get != "" {
		if string(get[0]) != "[" && string(get[len(get)-1]) != "]" {
			panic("You must set get variables in format of [ (var1,var2), (var2,var4), var3 ]")
		}
		getVars = parseVariables(get)
	}

	if update != "" {
		if string(update[0]) != "[" && string(update[len(update)-1]) != "]" {
			panic("You must set get variables in format of [ (var1,var2), (var2,var4), var3 ]")
		}
		updateVars = parseVariables(update)
	}

}

func parseVariables(vars string) []variables {
	// TODO : create parser for variables
}
