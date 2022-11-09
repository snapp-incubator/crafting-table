package repository

import (
	"github.com/snapp-incubator/crafting-table/internal/structure"
)

type Repo struct {
	Tag         string                      `yaml:"tag" json:"tag"`
	Source      string                      `yaml:"source" json:"source"`
	Destination string                      `yaml:"destination" json:"destination"`
	PackageName string                      `yaml:"package_name" json:"package_name"`
	StructName  string                      `yaml:"struct_name" json:"struct_name"`
	TableName   string                      `yaml:"table_name" json:"table_name"`
	DBLibrary   string                      `yaml:"db_library" json:"db_library"`
	Join        []structure.JoinVariables   `yaml:"join" json:"join"`
	Get         []structure.GetVariable     `yaml:"get" json:"get"`
	Update      []structure.UpdateVariables `yaml:"update" json:"update"`
	Create      structure.CreateVariables   `yaml:"create" json:"create"`
	Test        bool                        `yaml:"test" json:"create_test"`
}

type Manifest struct {
	Repos []Repo `yaml:"repositories"`
}

func (m *Repo) EqualTag(tags []string) bool {
	for _, tag := range tags {
		if m.Tag == tag {
			return true
		}
	}
	return false
}
