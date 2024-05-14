package desktop

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/cloyop/locker-desktop/storage"
)

type Locker struct {
	A         fyne.App
	W         fyne.Window
	M         *storage.Metadata
	Modal     *widget.PopUp
	ModalFunc func()
	View      func() fyne.CanvasObject
}

func NewLocker() *Locker {
	a := app.New()
	w := a.NewWindow("Locker")
	w.Resize(fyne.NewSize(420, 360))
	locker := &Locker{
		A: a,
		M: storage.NewMetaData(),
		W: w,
	}
	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if ke.Name == fyne.KeyEscape {
			locker.DropModal()
		}
	})
	return locker
}

func (locker *Locker) ShowItem(name string) {
	item, is := locker.M.Data[name]
	if !is {
		locker.ShowModal(centerLabel("Is not "+name), nil, true)
	}
	locker.ShowModal(ItemBox(locker, name, item), nil, false)
}

func (locker *Locker) ShowModal(content fyne.CanvasObject, f func(), withButton bool) {
	locker.DropModal()
	contentBox := container.NewVBox(content)
	mod := widget.NewModalPopUp(container.NewStack(contentBox), locker.W.Canvas())
	if withButton {
		button := widget.NewButtonWithIcon("", theme.CancelIcon(), locker.DropModal)
		contentBox.Add(container.NewCenter(button))
	}
	if f != nil {
		locker.ModalFunc = f
	}
	mod.Show()
	mod.Resize(fyne.Size{
		Width:  340,
		Height: mod.MinSize().Height,
	})
	locker.Modal = mod
}
func (locker *Locker) DropModal() {
	if locker.Modal != nil {
		locker.Modal.Hide()
		locker.Modal = nil
		if locker.ModalFunc != nil {
			locker.ModalFunc()
			locker.ModalFunc = nil
		}
	}
}
func (locker *Locker) Reload() {
	locker.W.SetContent(locker.View())
}
func (locker *Locker) ResetPinTimer() {
	if locker.M.LastAction != nil {
		locker.M.LastAction.Stop()
	}
	locker.M.LastAction = time.AfterFunc(time.Minute*10, func() {
		lastbox := locker.W.Content()
		if locker.Modal != nil {
			locker.Modal.Hide()
		}
		locker.W.SetContent(PinRequest(locker, lastbox))
	})
}
