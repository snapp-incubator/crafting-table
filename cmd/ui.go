package cmd

import (
	"github.com/rivo/tview"
	"github.com/spf13/cobra"

	"github.com/snapp-incubator/crafting-table/internal/app"
)

var uiCMD = &cobra.Command{
	Use:   "ui",
	Short: "Shows the terminal UI",
	Run:   ui,
}

func ui(_ *cobra.Command, _ []string) {
	var uiApp = tview.NewApplication()
	var pages = tview.NewPages()
	var form = tview.NewForm()

	pages.AddPage("Repository information", form, true, true)

	repo := app.Repository{}

	form.AddInputField("Source path", "", 500, nil, stringAssigner(&repo.Source))
	form.AddInputField("Destination path", "", 500, nil, stringAssigner(&repo.Destination))
	form.AddInputField("Package name", "", 50, nil, stringAssigner(&repo.PackageName))
	form.AddInputField("Field names for Get Seperated By Commas", "", 200, nil, stringAssigner(&repo.Get))
	form.AddInputField("Field names for Update Seperated By Commas", "", 200, nil, stringAssigner(&repo.Update))
	form.AddCheckbox("Generate Tests", false, boolAssigner(&repo.Test))
	form.AddCheckbox("Have Create Method", false, boolAssigner(&repo.Create))

	form.AddButton("Done", func() {
		generateRepository(repo)
		uiApp.Stop()
	})

	if err := uiApp.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func stringAssigner(lvalue *string) func(value string) {
	return func(value string) {
		*lvalue = value
	}
}

func boolAssigner(lvalue *bool) func(value bool) {
	return func(value bool) {
		*lvalue = value
	}
}
