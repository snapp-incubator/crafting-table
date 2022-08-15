package app

import "github.com/snapp-incubator/crafting-table/internal/structure"

type Repo struct {
	Tag         string                      `json:"tag"`
	Source      string                      `yaml:"source"`
	Destination string                      `yaml:"destination"`
	PackageName string                      `yaml:"package_name"`
	StructName  string                      `yaml:"struct_name"`
	TableName   string                      `yaml:"table_name"`
	DBLibrary   string                      `yaml:"db_library"`
	Get         []structure.GetVariable     `yaml:"get"`
	Update      []structure.UpdateVariables `yaml:"update"`
	Create      structure.CreateVariables   `yaml:"create"`
	Test        bool                        `yaml:"test"`
}

type Manifest struct {
	Repos []Repo `yml:"repositories"`
}