package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func dropModalButton(locker *Locker) *widget.Button {
	return widget.NewButtonWithIcon("", theme.CancelIcon(), func() { locker.DropModal() })
}
func copyToClipButton(s string, w fyne.Window) *widget.Button {
	return widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() { w.Clipboard().SetContent(s) })
}

func areUsure(adv string, yes, no func()) *fyne.Container {
	l := centerLabel(adv)
	y := widget.NewButton("Yes", yes)
	n := widget.NewButton("No", no)
	box := container.NewVBox(l, container.NewCenter(container.NewHBox(y, n)))
	return box
}
func centerLabel(msg string) *widget.Label {
	label := widget.NewLabel(msg)
	label.Alignment = fyne.TextAlignCenter
	return label
}
func newEntry(placeholder string) *widget.Entry {
	entry := widget.NewEntry()
	entry.PlaceHolder = placeholder
	return entry
}
func base() fyne.CanvasObject {
	return &widget.BaseWidget{}
}
