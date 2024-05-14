package desktop

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/cloyop/locker-desktop/pkg"
	"github.com/cloyop/locker-desktop/storage"
)

func ItemBox(locker *Locker, name string, kvs *storage.KeyValueStore) *fyne.Container {
	box := container.NewVBox(centerLabel(name))
	if kvs.Value != "" {
		box.Add(container.NewBorder(base(), base(), centerLabel(kvs.Key), copyToClipButton(kvs.Key, locker.W)))
		box.Add(container.NewBorder(base(), base(), centerLabel(kvs.Value), copyToClipButton(kvs.Value, locker.W)))
	} else {
		lbl := centerLabel(kvs.Key)
		lbl.Wrapping = fyne.TextWrapBreak
		box.Add(container.NewBorder(base(), base(), centerLabel("key"), copyToClipButton(kvs.Key, locker.W)))
		box.Add(lbl)
	}
	editButton := editItemButton(locker, kvs, name)
	deleteButton := deleteItemButton(locker, kvs, name)
	unPop := dropModalButton(locker)
	butts := container.NewCenter(container.NewHBox(editButton, deleteButton, unPop))
	box.Add(butts)
	return box
}
func editItemButton(locker *Locker, kvs *storage.KeyValueStore, name string) *widget.Button {
	return widget.NewButton("Edit", func() {
		nameLabel := centerLabel("Name")
		nameEntry := newEntry("Name input...")
		nameEntry.Text = name
		keyLabel := centerLabel("Key")
		keyEntry := newEntry("key input...")
		keyEntry.Text = kvs.Key

		valueLabel := centerLabel("Set value")
		valueEntry := newEntry("value input...")
		valueEntry.Text = kvs.Value

		button := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
			var changed bool
			if keyEntry.Text != kvs.Key {
				kvs.Key = keyEntry.Text
				locker.M.ChangesMade = true
				changed = true
			}
			if valueEntry.Text != kvs.Value {
				kvs.Value = valueEntry.Text
				locker.M.ChangesMade = true
				changed = true
			}
			if nameEntry.Text != name {
				ckvs, is := locker.M.Data[nameEntry.Text]
				if is {
					fmt.Println("here")
					lstM := locker.Modal
					fmt.Println(locker.Modal)
					locker.ShowModal(container.NewVBox(centerLabel("Item "+nameEntry.Text+" already exists"), centerLabel(fmt.Sprintf("%v\n%v", ckvs.Key, ckvs.Value))), func() {
						locker.Modal = lstM
						locker.Modal.Show()
					}, false)
					return
				}
				locker.M.Data[nameEntry.Text] = kvs
				delete(locker.M.Data, name)
				changed = true
			}
			if changed {
				locker.Reload()
				locker.ResetPinTimer()
				locker.ShowModal(ItemBox(locker, nameEntry.Text, kvs), nil, false)
			}
		})
		valueEntry.OnSubmitted = func(s string) { button.OnTapped() }
		keyEntry.OnSubmitted = func(s string) { button.OnTapped() }
		nameEntry.OnSubmitted = func(s string) { button.OnTapped() }
		cc := container.NewVBox(nameLabel, nameEntry, keyLabel, keyEntry, valueLabel, valueEntry, container.NewCenter(container.NewHBox(button)))
		locker.ShowModal(cc, nil, false)

	})
}
func deleteItemButton(locker *Locker, kvs *storage.KeyValueStore, name string) *widget.Button {
	return widget.NewButton("Delete", func() {
		go func() {
			y := func() {
				delete(locker.M.Data, name)
				locker.ShowModal(centerLabel(name+" removed"), nil, false)
				locker.M.ChangesMade = true
				locker.Reload()
				locker.ResetPinTimer()
				time.Sleep(time.Second * 1)
				locker.DropModal()
			}
			n := func() {
				locker.ShowModal(ItemBox(locker, name, kvs), nil, false)
			}
			b := areUsure("Are u sure u want to delete "+name+" ?", y, n)
			locker.ShowModal(b, nil, false)
		}()
	})
}
func AddItemButton(locker *Locker) *widget.Button {
	return widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		n := newEntry("Name")
		k := newEntry("Key")
		v := newEntry("Value")
		v.Text = pkg.RandomPassword()
		box := container.NewVBox(centerLabel("New item"), n, k, v)
		addButton := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
			var missing []string
			if n.Text == "" {
				missing = append(missing, "Missing name")
			}
			if k.Text == "" {
				missing = append(missing, "Missing Key")
			}
			if len(missing) > 0 {
				errBox := container.NewVBox()
				for _, v := range missing {
					errBox.Add(centerLabel(v))
				}
				locker.ShowModal(errBox, nil, false)
				time.Sleep(time.Second * 1)
				locker.ShowModal(box, nil, false)
				return
			}
			kvs, is := locker.M.Data[n.Text]
			if is {
				locker.ShowModal(centerLabel("Item "+n.Text+" already exists"), nil, false)
				time.Sleep(time.Second * 1)
				locker.ShowModal(ItemBox(locker, n.Text, kvs), nil, false)
				return
			}
			kvs = &storage.KeyValueStore{Key: k.Text, Value: v.Text}
			locker.M.Data[n.Text] = kvs
			locker.M.ChangesMade = true
			locker.ShowModal(ItemBox(locker, n.Text, kvs), nil, false)
			locker.ResetPinTimer()
			locker.Reload()
		})
		box.Add(container.NewCenter(container.NewHBox(addButton, dropModalButton(locker))))
		n.OnSubmitted = func(s string) { addButton.OnTapped() }
		k.OnSubmitted = func(s string) { addButton.OnTapped() }
		v.OnSubmitted = func(s string) { addButton.OnTapped() }
		locker.ShowModal(container.NewPadded(box), nil, false)
	})
}
