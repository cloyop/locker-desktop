package desktop

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func SettingsView(locker *Locker) fyne.CanvasObject {
	back := widget.NewButtonWithIcon("",
		theme.NavigateBackIcon(),
		func() {
			locker.View = func() fyne.CanvasObject { return HomeView(locker) }
			locker.Reload()
		},
	)

	changePassword := widget.NewButton("Change password", func() {
		currentPwd := newEntry("Current password")
		passwordEntry := newEntry("Password")
		box := container.NewVBox(currentPwd, passwordEntry)
		button := widget.NewButtonWithIcon("", theme.ConfirmIcon(),
			func() {
				if currentPwd.Text == "" || passwordEntry.Text == "" {
					return
				}
				if currentPwd.Text != locker.M.Password {
					locker.ShowModal(centerLabel("Current password invalid"), func() { locker.ShowModal(box, nil, false) }, true)
					return
				}
				if len(passwordEntry.Text) < 8 {
					locker.ShowModal(centerLabel("Password must be min 8 length"), func() { locker.ShowModal(box, nil, false) }, true)
					return
				}
				if passwordEntry.Text == locker.M.Password {
					locker.ShowModal(centerLabel("Cannot use the same password"), func() { locker.ShowModal(box, nil, false) }, true)
					return
				}
				locker.M.Password = passwordEntry.Text
				if err := locker.M.Save(); err != nil {
					locker.ShowModal(centerLabel(err.Error()), nil, false)
					return
				}
				locker.ShowModal(centerLabel("password changed succesfully"), nil, false)
				time.Sleep(time.Second)
				locker.DropModal()
			})
		box.Add(container.NewCenter(container.NewHBox(button, dropModalButton(locker))))
		passwordEntry.OnSubmitted = func(s string) { button.OnTapped() }
		currentPwd.OnSubmitted = func(s string) { button.OnTapped() }
		locker.ShowModal(box, nil, false)
	})

	changePin := widget.NewButton("Change pin",
		func() {
			currentPin := newEntry("Current pin")
			pinEntry := newEntry("pin")
			box := container.NewVBox(currentPin, pinEntry)
			button := widget.NewButtonWithIcon("", theme.ConfirmIcon(),
				func() {
					if currentPin.Text == "" || pinEntry.Text == "" {
						return
					}
					if currentPin.Text != locker.M.Pin {
						locker.ShowModal(centerLabel("Current pin invalid"), func() { locker.ShowModal(box, nil, false) }, true)
						return
					}
					if len(pinEntry.Text) < 4 {
						locker.ShowModal(centerLabel("Pin must be min 4 length"), func() { locker.ShowModal(box, nil, false) }, true)
						return
					}
					if pinEntry.Text == locker.M.Pin {
						locker.ShowModal(centerLabel("Cannot use the same pin"), func() { locker.ShowModal(box, nil, false) }, true)
						return
					}
					locker.M.Pin = pinEntry.Text
					if err := locker.M.Save(); err != nil {
						locker.ShowModal(centerLabel(err.Error()), nil, false)
						return
					}
					locker.ShowModal(centerLabel("password changed succesfully"), nil, false)
					time.Sleep(time.Second)
					locker.DropModal()
				},
			)
			pinEntry.OnSubmitted = func(s string) { button.OnTapped() }
			currentPin.OnSubmitted = func(s string) { button.OnTapped() }
			box.Add(container.NewCenter(container.NewHBox(button, dropModalButton(locker))))
			locker.ShowModal(box, nil, false)
		})

	lfj := widget.NewButton("load from json file", func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil {
				locker.ShowModal(centerLabel(err.Error()), nil, true)
				return
			}
			if uc != nil {
				if uc.URI().Extension() != ".json" {
					locker.ShowModal(centerLabel("only .json files"), nil, true)
					return
				}
				go func() {
					y := func() {
						n, err := locker.M.SetFromJson(uc)
						if err != nil {
							locker.ShowModal(centerLabel(err.Error()), nil, true)
							return
						}
						locker.ShowModal(centerLabel(fmt.Sprintf("%d items loaded", n)), nil, false)
						time.Sleep(time.Second)
						locker.DropModal()
					}
					c := areUsure("This action will overwrite items with the same name", y, locker.DropModal)
					locker.ShowModal(c, nil, false)
				}()
			}

		}, locker.W)

	})

	return container.NewPadded(container.NewCenter(container.NewVBox(back, changePassword, changePin, lfj)))
}
