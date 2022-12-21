package days

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

func makeSelectScreen() fyne.CanvasObject {
	day := newDayEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{{Text: "Day", Widget: day}},
		OnSubmit: func() {
			fmt.Println("Form submitted with day " + day.Text)
			// TODO: select the specified day in the menu
		},
	}

	return container.NewCenter(form)
}

type dayEntry struct {
	widget.Entry
}

func newDayEntry() *dayEntry {
	e := &dayEntry{}
	e.ExtendBaseWidget(e)
	e.Validator = validation.NewRegexp(`\d`, "Must contain a number")
	return e
}
