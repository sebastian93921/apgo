package ui

import (
	"apgo/system"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (a *Apgoui) makeMenu(settings *system.Settings) *fyne.MainMenu {
	menu := fyne.NewMainMenu(
		fyne.NewMenu("File"))

	menu.Items = append(menu.Items,
		a.makeHelpMenu(settings, a.Win),
	)
	return menu
}

func (a *Apgoui) makeHelpMenu(settings *system.Settings, w fyne.Window) *fyne.Menu {
	return fyne.NewMenu("Help",
		fyne.NewMenuItem("Settings", func() {
			downloadCaButton := generateDownloadViewAndButton(settings, w)
			form := &widget.Form{
				Items: []*widget.FormItem{ // we can specify items in the constructor
					{Text: "Export CA Cert", Widget: downloadCaButton}},
				OnSubmit: func() { // optional, handle form submission
					log.Println("Form submitted:")
				},
			}
			d := dialog.NewCustom("Settings", "Close", form, w)
			d.Show()
		}),
	)
}

func generateDownloadViewAndButton(settings *system.Settings, w fyne.Window) *widget.Button {
	return widget.NewButton("Save to file", func() {
		saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if writer == nil {
				return
			}

			// Write the CA certificate to the file
			_, err = writer.Write(settings.CACertificate)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			// Close the file
			err = writer.Close()
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			// Show a success message
			dialog.ShowInformation("Success", "The CA certificate has been saved successfully.", w)
		}, w)
		saveDialog.SetFileName("ca.pem")
		saveDialog.Show()
	})
}
