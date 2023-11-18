package proxyserv

import (
	"apgo/modules/repeater"
	"apgo/util"
	"apgo/util/decoder"
	"fmt"
	"log"
	"net/http/httputil"
	"sort"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var isDecending = true

func (p *ProxyImpl) GenerateUI(mainTabs *container.AppTabs, repeaterTabs *container.DocTabs) fyne.CanvasObject {
	headers := []string{"Method", "URL", "Content Length", "Time"}
	table := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(p.LogIDs), len(headers)
		},
		func() fyne.CanvasObject {
			l := widget.NewLabel("placeholder")
			l.Truncation = fyne.TextTruncateEllipsis
			return l
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			ids := p.LogIDs
			if isDecending {
				sort.Slice(ids, func(i, j int) bool { return ids[i] > ids[j] }) // Sort the keys in descending order
			}
			logId := ids[i.Row]
			logEntry := p.LogEntry[logId] // Get the log entry
			var text string
			switch i.Col {
			case 0:
				text = logEntry.Method
			case 1:
				text = logEntry.URL
			case 2:
				text = fmt.Sprintf("%d", logEntry.ContentLength)
			case 3:
				text = logEntry.Time.Format(time.RFC3339)
			}
			label.SetText(text) // Set the cell content to the corresponding column of the log entry
		},
	)
	table.SetColumnWidth(0, 100)
	table.SetColumnWidth(1, 300)
	table.SetColumnWidth(2, 120)
	table.SetColumnWidth(3, 200)

	table.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("000")
	}

	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		l := o.(*widget.Label)
		if id.Col == -1 {
			l.SetText(fmt.Sprintf("%d", p.LogIDs[id.Row]))
		} else {
			l.SetText(headers[id.Col])
		}
	}

	// Create request and response panels
	requestPanel := widget.NewMultiLineEntry()
	requestPanel.Wrapping = fyne.TextWrapBreak
	responsePanel := widget.NewMultiLineEntry()
	responsePanel.Wrapping = fyne.TextWrapBreak

	selectedTarget := ""
	table.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 {
			selectedRow := id.Row
			go func() {
				ids := p.LogIDs
				if isDecending {
					sort.Slice(ids, func(i, j int) bool { return ids[i] > ids[j] }) // Sort the keys in descending order
				}
				logId := ids[id.Row]
				if selectedRow == id.Row {
					logEntry := p.LogEntry[logId] // Get the selected log entry

					printToEntry(requestPanel, logEntry.RequestMessage, logEntry.RequestBody)
					printToEntry(responsePanel, logEntry.ResponseMessage, logEntry.ResponseBody)

					selectedTarget = logEntry.Target
				}
			}()
		}
	}

	repeaterCount := 1

	// Create buttons to copy the text from the requestPanel to the repeaterPanel
	sendToRepeater := widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), func() {
		// requestBinding := binding.BindString(&requestPanel.Text)
		requestBinding := binding.NewString()
		requestBinding.Set(requestPanel.Text)
		targetBinding := binding.NewString()

		targetEntry := widget.NewEntry()
		targetEntry.SetPlaceHolder("Target")
		targetEntry.Bind(targetBinding)
		targetBinding.Set(selectedTarget)

		repeaterRequestPanel := widget.NewMultiLineEntry()
		repeaterRequestPanel.Wrapping = fyne.TextWrapBreak
		repeaterRequestPanel.Validator = nil
		repeaterResponsePanel := widget.NewMultiLineEntry()
		repeaterResponsePanel.Wrapping = fyne.TextWrapBreak
		repeaterSplitContainer := container.NewVSplit(repeaterRequestPanel, repeaterResponsePanel)

		// Send the request
		sendButton := widget.NewButton("Send", func() {
			go func() {
				repeaterResponsePanel.SetText("Sending...")
				repeaterRequestPanel.Disable()

				requestString, _ := requestBinding.Get()
				targetString, _ := targetBinding.Get()

				resp, err := repeater.SendRawRequest(targetString, requestString)
				if err != nil {
					log.Println(targetString, requestString)
					log.Panicf("Failed to send request: %v", err)
				}

				if resp != nil {
					bodyBytes, err := decoder.DecodeBodyResponse(resp)
					if err != nil {
						log.Panicln("Error decoding response:", err)
					}

					dump, err := httputil.DumpResponse(resp, false)
					if err != nil {
						log.Panicln("Error dumping response:", err)
						return
					}

					printToEntry(repeaterResponsePanel, dump, bodyBytes)
				} else {
					repeaterResponsePanel.SetText("Error on receiving response")
				}
				repeaterRequestPanel.Enable()
			}()
		})
		toolPanel := container.NewHSplit(sendButton, targetEntry)

		content := container.NewBorder(toolPanel, nil, nil, nil, repeaterSplitContainer)

		repeaterTabs.Append(container.NewTabItem(strconv.Itoa(repeaterCount), content))
		repeaterRequestPanel.Bind(requestBinding)

		repeaterCount += 1

		mainTabs.SelectIndex(2)
		lastTabIndex := len(repeaterTabs.Items) - 1
		repeaterTabs.SelectIndex(lastTabIndex)

	})
	toolBox := container.NewVBox(sendToRepeater)

	logEntryNestedTabs := container.NewAppTabs(
		container.NewTabItem("Request", requestPanel),
		container.NewTabItem("Response", responsePanel),
	)

	customLayout := container.NewBorder(nil, nil, nil, toolBox, logEntryNestedTabs)

	// Create a split container with the table on the left and the request and response panels on the right
	logSplitContainer := container.NewVSplit(table, customLayout)

	return logSplitContainer
}

func printToEntry(entry *widget.Entry, message []byte, body []byte) {
	entry.SetText(util.ConvertToReadableString(message, body))
}
