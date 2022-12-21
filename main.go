package main

import (
	"fmt"
	"strconv"

	"example.com/advent2022/days"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Friendly's Advent of code 2022")
	w.SetMaster()

	title := widget.NewLabel("Day Name")
	content := container.NewMax()

	setDay := func(d days.Day) {
		title.SetText("Day " + strconv.Itoa(d.Number))
		vbox := container.NewVBox()
		vbox.Add(widget.NewLabel("Part A:"))
		vbox.Add(widget.NewLabel(d.PartAPrompt))
		solveButtonA := widget.Button{Text: "Solve part A"}
		solveButtonA.OnTapped = func() {
			solveButtonA.Disable()

			// TODO: use progress bar
			vbox.Add(widget.NewLabel("Testing part A..."))
			allTestsPassed := true
			for i, test := range d.PartATests {
				actualOutput, testImg, err := d.Solver.SolvePartA(test.Input)
				if err != nil {
					fail_str := fmt.Sprint("Day ", d.Number, " part A: test ", i, " returned err: ", err.Error())
					fmt.Println(fail_str)
					vbox.Add(widget.NewLabel(fail_str))
					allTestsPassed = false
				} else if actualOutput != test.ExpectedOutput {
					fail_str := fmt.Sprint("Day ", d.Number, " part A: test ", i, " failed. Got: ", actualOutput, ", expected: ", test.ExpectedOutput)
					fmt.Println(fail_str)
					vbox.Add(widget.NewLabel(fail_str))
					if testImg != nil {
						vbox.Add(widget.NewLabel("Part A failed test image:"))
						vbox.Add(testImg)
					}
					allTestsPassed = false
				}

				if i == 0 {
					vbox.Add(widget.NewLabel("Part A test 0 image:"))
					vbox.Add(testImg)
				}
			}
			if !allTestsPassed {
				return
			}

			vbox.Add(widget.NewLabel("All tests passed"))
			vbox.Add(widget.NewLabel("Solving part A..."))
			solution, solutionImg, err := d.Solver.SolvePartA(d.PuzzleInput)
			if err != nil {
				fail_str := fmt.Sprint("Day ", d.Number, " solution returned err: ", err.Error())
				fmt.Println(fail_str)
				vbox.Add(widget.NewLabel(fail_str))
				return
			}
			result := widget.NewEntry()
			result.SetText(solution)
			result.Disable()
			vbox.Add(widget.NewForm(&widget.FormItem{Text: "Puzzle solution is: ", Widget: result}))
			fmt.Println("Part A solution is: ", solution)

			vbox.Add(solutionImg)
		}

		vbox.Add(&solveButtonA)

		vbox.Add(widget.NewLabel("Part B:"))
		vbox.Add(widget.NewLabel(d.PartBPrompt))

		solveButtonB := widget.Button{Text: "Solve part B"}
		solveButtonB.OnTapped = func() {
			solveButtonB.Disable()

			// TODO: use progress bar
			vbox.Add(widget.NewLabel("Testing part B..."))
			allTestsPassed := true
			for i, test := range d.PartBTests {
				actualOutput, testImg, err := d.Solver.SolvePartB(test.Input)
				if err != nil {
					fail_str := fmt.Sprint("Day ", d.Number, " part B: test ", i, " returned err: ", err.Error())
					fmt.Println(fail_str)
					vbox.Add(widget.NewLabel(fail_str))
					allTestsPassed = false
				} else if actualOutput != test.ExpectedOutput {
					fail_str := fmt.Sprint("Day ", d.Number, " part B: test ", i, " failed. Got: ", actualOutput, ", expected: ", test.ExpectedOutput)
					fmt.Println(fail_str)
					vbox.Add(widget.NewLabel(fail_str))
					if testImg != nil {
						vbox.Add(widget.NewLabel("Part B failed test image:"))
						vbox.Add(testImg)
					}
					allTestsPassed = false
				}

				if i == 0 {
					vbox.Add(widget.NewLabel("Part B test 0 image:"))
					vbox.Add(testImg)
				}
			}
			if !allTestsPassed {
				return
			}

			vbox.Add(widget.NewLabel("All tests passed"))
			vbox.Add(widget.NewLabel("Solving part B..."))
			solution, sol_img, err := d.Solver.SolvePartB(d.PuzzleInput)
			if err != nil {
				fail_str := fmt.Sprint("Day ", d.Number, " Part B solution returned err: ", err.Error())
				fmt.Println(fail_str)
				vbox.Add(widget.NewLabel(fail_str))
				return
			}
			result := widget.NewEntry()
			result.SetText(solution)
			result.Disable()
			vbox.Add(widget.NewForm(&widget.FormItem{Text: "Puzzle solution is: ", Widget: result}))
			fmt.Println("Part B solution is: ", solution)

			vbox.Add(sol_img)
		}

		vbox.Add(&solveButtonB)

		content.Objects = []fyne.CanvasObject{container.NewVScroll(vbox)}
		content.Refresh()
	}

	day := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, content)
	w.SetContent(container.NewHSplit(makeNav(setDay), day))
	w.Resize(fyne.NewSize(1500, 1000))
	w.ShowAndRun()
}

func makeNav(setDay func(days.Day)) fyne.CanvasObject {
	list := &widget.List{
		Length: func() int {
			return len(days.Days)
		},
		CreateItem: func() fyne.CanvasObject {
			return widget.NewLabel("Day TEMPLATE")
		},
		UpdateItem: func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText("Day " + strconv.Itoa(days.Days[id].Number))
		},
		OnSelected: func(id widget.ListItemID) {
			setDay(days.Days[id])
		},
	}

	return list
}
