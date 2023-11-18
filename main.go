package main

import (
	"apgo/system"
	"apgo/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&system.MyTheme{})

	myWindow := myApp.NewWindow("Attack Proxy Go")
	myWindow.Resize(fyne.NewSize(1024, 800))

	app := &ui.Apgoui{Win: myWindow}
	app.StartUi()
}
