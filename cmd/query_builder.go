package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/snapp-incubator/crafting-table/internal/querybuilder"
)

var (
	dialect  string
	filePath string
	table    string
)

var queryBuilderCmd = &cobra.Command{
	Use:     "query-builder",
	Aliases: []string{"qb"},
	Short:   `Generates a query builder based on annotations in your code`,
	Run: func(cmd *cobra.Command, args []string) {
		generateQueryBuilder(cmd, filePath)
	},
}

func init() {
	queryBuilderCmd.Flags().StringVarP(&filePath, "file-path", "f", "", "which file you want to parse and generate query builder for")
	queryBuilderCmd.Flags().StringVarP(&dialect, "dialect", "d", "mysql", "dialect you want to generate sql for")
	queryBuilderCmd.Flags().StringVarP(&table, "table", "t", "", "table name of the type if not specified defaults to snakeCase(plural(typeName))")
}

func generateQueryBuilder(cmd *cobra.Command, filePath string) {
	if filePath == "" {
		fmt.Println("You need to fill --file-path|-f flag")
		_ = cmd.Help()
		return
	}

	inputFilePath, err := filepath.Abs(filePath)
	if err != nil {
		panic(err)
	}

	pathList := filepath.SplitList(inputFilePath)
	pathList = pathList[:len(pathList)-1]
	fileDir := filepath.Join(pathList...)
	fileSet := token.NewFileSet()
	fileAst, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	actualName := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	outputFilePath := filepath.Join(fileDir, fmt.Sprintf("%s_ct_gen.go", actualName))
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(outputFile)

	for _, decl := range fileAst.Decls {
		if _, ok := decl.(*ast.GenDecl); ok {
			declComment := decl.(*ast.GenDecl).Doc.Text()
			if len(declComment) > 0 && declComment[:len(querybuilder.ModelAnnotation)] == querybuilder.ModelAnnotation {
				output := querybuilder.Generate(dialect, fileAst.Name.String(), decl.(*ast.GenDecl))
				_, err := fmt.Fprint(outputFile, output)
				if err != nil {
					panic(err)
				}
			}
		}

	}
}
