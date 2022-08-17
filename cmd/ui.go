package cmd

import (
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var uiCMD = &cobra.Command{
	Use:   "ui",
	Short: "Shows the terminal UI",
	Run:   ui,
}

func ui(cmd *cobra.Command, args []string) {
	var uiApp = tview.NewApplication()
	var pages = tview.NewPages()
	var form = tview.NewForm()

	pages.AddPage("Repository information", form, true, true)

	form.AddInputField("Source path", "", 500, nil, stringAssigner(&source))
	form.AddInputField("Destination path", "", 500, nil, stringAssigner(&destination))
	form.AddInputField("Package name", "", 50, nil, stringAssigner(&packageName))
	form.AddInputField("Struct name", "", 50, nil, stringAssigner(&structName))
	form.AddInputField("Field names for Get Seperated By Commas", "", 200, nil, stringAssigner(&get))
	form.AddInputField("Field names for Update Seperated By Commas", "", 200, nil, stringAssigner(&update))
	form.AddCheckbox("Generate Tests", false, boolAssigner(&test))
	form.AddCheckbox("Have Create Method", false, boolAssigner(&create))

	form.AddButton("Done", func() {
		generate(cmd, args)
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
