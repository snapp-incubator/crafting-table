package cmd

import (
	"log"
	"os"

	"github.com/snapp-incubator/crafting-table/internal/build"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var (
	manifestPath string
	tags         string
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
	applyCMD.Flags().StringVarP(&tags, "tags", "t", "", "select tags from ct-manifest file for generating repositories")
}

func apply(_ *cobra.Command, _ []string) {
	if manifestPath == "" {
		panic("manifest path is not set")
	}

	var manifest build.Repo
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

	if err := build.Generate(manifest); err != nil {
		log.Fatal(err)
	}
}
