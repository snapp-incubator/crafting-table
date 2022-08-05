package app

type Repository struct {
	Source      string `yml:"source"`
	Destination string `yml:"destination"`
	PackageName string `yml:"packageName"`
	Get         string `yml:"get"`
	Update      string `yml:"update"`
	Create      bool   `yml:"create"`
	Test        bool   `yml:"test"`
}

type Repositories struct {
	Repositories []Repository `yml:"repositories"`
}
