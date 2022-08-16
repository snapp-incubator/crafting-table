package cmd

import (
	"log"
	"os"

	"github.com/snapp-incubator/crafting-table/internal/app"
	"github.com/snapp-incubator/crafting-table/internal/repository"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	manifestPath string
)

var manifestCMD = &cobra.Command{
	Use:   "manifest",
	Short: "Manifest commands",
}

var applyCMD = &cobra.Command{
	Use:   "apply",
	Short: "Create repositories from manifest file",
	Run:   apply,
}

func init() {
	applyCMD.Flags().StringVarP(&manifestPath, "manifest-path", "p", "", "generate automatically repositories from ct-manifest file")
}

func apply(_ *cobra.Command, _ []string) {
	if manifestPath == "" {
		panic("manifest path is not set")
	}

	var manifest app.Manifest
	file, err := os.Open(manifestPath)
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
	if err := d.Decode(&manifest); err != nil {
		panic(err)
	}

	for _, repo := range manifest.Repos {
		if err := repository.Generate(repo.Source, repo.Destination, repo.PackageName, repo.StructName, &repo.Get, &repo.Update, repo.Create.Enable, repo.Test); err != nil {
			log.Fatal(err)
		}
	}
}