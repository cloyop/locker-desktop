package desktop

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/cloyop/locker-desktop/pkg"
)

func PinRequest(locker *Locker, lastBox fyne.CanvasObject) fyne.CanvasObject {
	label := centerLabel("You has been afk, pin required")
	entry := newEntry("insert pin")
	entry.Resize(fyne.NewSize(180, entry.MinSize().Height))
	button := widget.NewButton("enter", func() {
		if entry.Text == "" {
			return
		}
		if entry.Text != locker.M.Pin {
			m := widget.NewModalPopUp(centerLabel("Invalid pin"), locker.W.Canvas())
			m.Show()
			time.Sleep(time.Second * 1)
			m.Hide()
			return
		}
		locker.W.SetContent(lastBox)
		if locker.Modal != nil {
			locker.Modal.Show()
		}
		locker.ResetPinTimer()
	})
	entry.OnSubmitted = func(s string) { button.OnTapped() }
	box := container.NewVBox(label, entry, button)
	return container.NewGridWithRows(3, base(), container.NewPadded(box))
}
func MakeLocker(locker *Locker) fyne.CanvasObject {
	pwdInfo := "The password is used to decrypt the first part of your data.\n even if it leaks. your data would remain safe by your pin.\nMinimum 8 length."
	pinInfo := "After 5 minutes of inactivity will mark you as AFK\nWill be asked your pin to allow you read or write\nMinimun 4 length"

	pwdLabel := container.NewBorder(base(), base(), centerLabel("Set a password"), widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		locker.ShowModal(centerLabel(pwdInfo), nil, true)
	}))
	pinLabel := container.NewBorder(base(), base(), centerLabel("Set a pin"), widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		locker.ShowModal(centerLabel(pinInfo), nil, true)
	}))
	pwdEntry := newEntry("password")
	pinEntry := newEntry("pin")
	button := widget.NewButton("Enter", func() {
		invalid := []fyne.CanvasObject{}
		if len(pwdEntry.Text) < 8 {
			invalid = append(invalid, centerLabel("Password minimun 8 length"))
		}
		if len(pinEntry.Text) < 4 {
			invalid = append(invalid, centerLabel("Pin minimun 4 length"))
		}
		if len(invalid) != 0 {
			locker.ShowModal(container.NewVBox(invalid...), nil, true)
			return
		}
		locker.M.Password = pwdEntry.Text
		locker.M.Pin = pinEntry.Text
		locker.M.Save()
		locker.ResetPinTimer()
		locker.View = func() fyne.CanvasObject { return HomeView(locker) }
		locker.Reload()
	})
	pinEntry.OnSubmitted = func(s string) { button.OnTapped() }
	pwdEntry.OnSubmitted = func(s string) { button.OnTapped() }
	return container.NewPadded(container.NewGridWithRows(3,
		container.NewVBox(pwdLabel, pwdEntry),
		container.NewVBox(pinLabel, pinEntry),
		container.NewCenter(button),
	))
}
func UnLockView(locker *Locker) fyne.CanvasObject {
	data, err := os.ReadFile(os.Getenv("LOCKER_PATH") + "/locker.txt")
	if err != nil {
		log.Fatal(err)
	}
	pwd := newEntry("Locker password")
	pwdB := widget.NewButton("Unlock", nil)
	box := container.NewVBox(pwd, container.NewCenter(pwdB))
	pwdB.OnTapped = func() { Unlock(pwd.Text, locker, &data, box) }
	pwd.OnSubmitted = func(s string) { pwdB.OnTapped() }
	return container.NewPadded(container.NewGridWithRows(3, base(), box))
}
func Unlock(password string, locker *Locker, firstLayerData *[]byte, box *fyne.Container) {
	if password == "" {
		return
	}
	unC, err := pkg.UnCipher([]byte(password), *firstLayerData)
	if err != nil {
		locker.ShowModal(centerLabel("Password Incorrect \n"+err.Error()), nil, true)
		return
	}
	locker.M.Password = password
	pin := newEntry("Write && Read Pin")
	pinB := widget.NewButton("Enter", nil)
	pinB.OnTapped = func() {
		if pin.Text == "" {
			return
		}
		StoreDataBytes, err := pkg.UnCipher([]byte(locker.M.Password+pin.Text), *unC)
		if err != nil {
			locker.ShowModal(centerLabel("Pin Incorrect \n"+err.Error()), nil, true)
			return
		}
		buffer := new(bytes.Buffer)
		_, err = buffer.Write(*StoreDataBytes)
		if err != nil {
			locker.ShowModal(centerLabel(err.Error()), nil, true)
			return
		}
		if err := gob.NewDecoder(buffer).Decode(&locker.M.Data); err != nil {
			locker.ShowModal(centerLabel(err.Error()), nil, true)
			return
		}
		locker.M.Pin = pin.Text
		locker.View = func() fyne.CanvasObject { return HomeView(locker) }
		locker.Reload()
		locker.ResetPinTimer()
	}
	pin.OnSubmitted = func(s string) { pinB.OnTapped() }
	box.RemoveAll()
	box.Add(pin)
	box.Add(container.NewCenter(pinB))
}
