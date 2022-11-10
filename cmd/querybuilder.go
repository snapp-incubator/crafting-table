package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/snapp-incubator/crafting-table/internal/querybuilder"
)

var dialect string
var table string
var querybuilderCmd = &cobra.Command{
	Use: "qb",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalln("needs a filename and dialect")
		}
		filename := args[0]

		inputFilePath, err := filepath.Abs(filename)
		if err != nil {
			panic(err)
		}
		pathList := filepath.SplitList(inputFilePath)
		pathList = pathList[:len(pathList)-1]
		dir := filepath.Join(pathList...)
		fset := token.NewFileSet()
		fast, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)

		if err != nil {
			panic(err)
		}
		actualName := strings.TrimSuffix(filename, filepath.Ext(filename))
		outputFilePath := filepath.Join(dir, fmt.Sprintf("%s_sqlgen_gen.go", actualName))
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()
		for _, decl := range fast.Decls {
			if _, ok := decl.(*ast.GenDecl); ok {
				declComment := decl.(*ast.GenDecl).Doc.Text()
				if len(declComment) > 0 && declComment[:len(querybuilder.ModelAnnotation)] == querybuilder.ModelAnnotation {
					args := make(map[string]string)
					output := querybuilder.Generate(dialect, fast.Name.String(), decl.(*ast.GenDecl), args)
					fmt.Fprint(outputFile, output)
				}
			}

		}
	},
}

func init() {
	querybuilderCmd.Flags().StringVarP(&dialect, "dialect", "d", "mysql", "dialect you want to generate sql for")
	querybuilderCmd.Flags().StringVarP(&table, "table", "t", "", "table name of the type if not specified defaults to snakeCase(plural(typeName))")
}
