package ui

import (
	"apgo/modules/intercept"
	"apgo/modules/proxyserv"
	"apgo/system"
	"apgo/util"
	"log"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (a *Apgoui) StartUi() {
	// Log Redirect
	logPanel := widget.NewMultiLineEntry()
	logPanel.Disable()
	// log.SetOutput(&util.LogWriter{Panel: logPanel})
	log.Default().SetOutput(&util.LogWriter{Panel: logPanel})

	// Initialize
	log.Println("Loading config...")

	settings := system.LoadGlobalSettings("./")
	session := system.NewSession(settings)

	// Create Interceptor
	interceptor := intercept.NewInterceptor(session)

	// Initialize Proxy
	proxy := proxyserv.NewProxy(session)
	proxyimpl := proxyserv.NewProxyImpl(proxy, session, interceptor)

	// Create main tab
	mainTabs := container.NewAppTabs()

	// Create repeater request and response panels
	repeaterTabs := container.NewDocTabs()

	// Generate ui template
	interceptorContainer := interceptor.GenerateUI(mainTabs)
	logSplitContainer := proxyimpl.GenerateUI(mainTabs, repeaterTabs)

	// Create tabs with table in the first tab
	mainTabs.Append(container.NewTabItem("Intercept", interceptorContainer))
	mainTabs.Append(container.NewTabItem("Log Entry", logSplitContainer))
	mainTabs.Append(container.NewTabItem("Repeater", repeaterTabs))
	mainTabs.Append(container.NewTabItem("Application Logs", logPanel))

	// Create status label
	status := widget.NewLabel("Status: OK")

	// Combine all components in a BorderLayout
	content := container.NewBorder(nil, status, nil, nil, mainTabs)

	a.Win.SetMainMenu(a.makeMenu(settings))
	a.Win.SetContent(content)

	go func() {
		log.Println("Starting proxy server...")
		err := proxyimpl.Module.Start()
		if err != nil {
			log.Fatalf("Failed to start proxy: %v", err)
		}
	}()

	a.Win.ShowAndRun()
}
