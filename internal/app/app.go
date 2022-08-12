package app

type Repository struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
	PackageName string `yaml:"package_name"`
	StructName  string `yaml:"struct_name"`
	Get         string `yaml:"get"`
	Update      string `yaml:"update"`
	Create      bool   `yaml:"create"`
	Test        bool   `yaml:"test"`
}

type Repositories struct {
	Repositories []Repository `yml:"repositories"`
}
