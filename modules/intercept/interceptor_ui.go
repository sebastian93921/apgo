package intercept

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (p *Interceptor) GenerateUI(mainTabs *container.AppTabs) fyne.CanvasObject {
	p.targetEntry = widget.NewEntry()
	p.targetEntry.SetPlaceHolder("Target")
	p.targetEntry.Bind(p.targetBinding)

	p.requestPanel = widget.NewMultiLineEntry()
	p.requestPanel.Wrapping = fyne.TextWrapBreak
	p.requestPanel.Validator = nil
	p.requestPanel.Bind(p.requestBinding)

	p.responsePanel = widget.NewMultiLineEntry()
	p.responsePanel.Wrapping = fyne.TextWrapBreak
	p.responsePanel.Disable()

	nestedTabs := container.NewAppTabs(
		container.NewTabItem("Request", p.requestPanel),
		container.NewTabItem("Response", p.responsePanel),
	)

	// Forward the request
	sendButton := widget.NewButton("Forward", func() {
		go func() {
			requestBody, _ := p.requestBinding.Get()
			if len(requestBody) > 0 {
				p.requestPanel.Disable()

				p.forwardChan <- true
				p.requestPanel.Enable()
				nestedTabs.SelectIndex(1)
			}
		}()
	})
	sendButton.Disable()
	dropButton := widget.NewButton("Drop", func() {
		go func() {
			requestBody, _ := p.requestBinding.Get()
			if len(requestBody) > 0 {
				p.requestPanel.Disable()

				p.dropChan <- true
				p.requestPanel.Enable()
			}
		}()
	})
	dropButton.Disable()
	var interceptToggle *widget.Button
	interceptToggle = widget.NewButton("Stop Intercept", func() {
		go func() {
			if p.isIntercepting {
				p.isIntercepting = false
				interceptToggle.SetText("Stop Intercept")
			} else {
				p.isIntercepting = true
				interceptToggle.SetText("Start Intercept")
			}
		}()
	})

	// Check if the panel has content
	p.requestPanel.OnChanged = func(content string) {
		if len(content) > 0 && content != "" {
			sendButton.Enable()
			dropButton.Enable()
		} else {
			sendButton.Disable()
			dropButton.Disable()
		}
	}

	sendDropPanel := container.NewHSplit(sendButton, dropButton)
	toggleTargetPanel := container.NewHSplit(interceptToggle, p.targetEntry)
	toolPanel := container.NewHSplit(sendDropPanel, toggleTargetPanel)

	return container.NewBorder(toolPanel, nil, nil, nil, nestedTabs)
}
