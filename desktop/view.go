package desktop

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func HomeView(locker *Locker) fyne.CanvasObject {
	box := container.NewVBox()
	var itemsBox *fyne.Container
	if len(locker.M.Data) == 0 {
		itemsBox = container.NewCenter(centerLabel("No items"))
	} else {
		itemsBox = container.NewAdaptiveGrid(4)
		for k := range locker.M.Data {
			var name = k
			itemsBox.Add(widget.NewButton(name, func() { locker.ShowItem(name) }))
		}
	}
	box.Add(Toolbar(locker))
	if len(locker.M.Data) > 0 {
		box.Add(searchBar(locker, itemsBox))
	}

	box.Add(itemsBox)
	return container.NewPadded(box)
}
func searchBar(locker *Locker, box *fyne.Container) *widget.Entry {
	searchEntry := newEntry("Search")
	var filtered bool
	searchEntry.OnChanged = func(s string) {
		if s == "" && filtered {
			box.RemoveAll()
			for k := range locker.M.Data {
				var name = k
				box.Add(widget.NewButton(name, func() { locker.ShowItem(name) }))
			}
			filtered = false
			return
		}
		filtered = true
		box.RemoveAll()
		for k := range locker.M.Data {
			var name = k
			if strings.Contains(strings.ToLower(name), s) {
				box.Add(widget.NewButton(name, func() { locker.ShowItem(name) }))
			}
		}
	}
	return searchEntry
}
func Toolbar(locker *Locker) *fyne.Container {
	addItem := AddItemButton(locker)
	save := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
		if locker.M.ChangesMade {
			if err := locker.M.Save(); err != nil {
				locker.ShowModal(centerLabel(err.Error()), nil, true)
				return
			}
			locker.M.ChangesMade = false
			locker.ShowModal(centerLabel("changes saved succesfully"), nil, true)
		} else {
			locker.ShowModal(centerLabel("nothing not save"), nil, true)
		}
	})
	menu := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		locker.View = func() fyne.CanvasObject { return SettingsView(locker) }
		locker.Reload()
	})

	cols := container.NewAdaptiveGrid(3, addItem, save, menu)
	return cols
}
