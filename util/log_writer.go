package util

import (
	"fmt"

	"fyne.io/fyne/v2/widget"
)

type LogWriter struct {
	Panel *widget.Entry
}

func (w *LogWriter) Write(p []byte) (n int, err error) {
	w.Panel.SetText(w.Panel.Text + string(p))
	fmt.Printf(string(p))
	return len(p), nil
}
