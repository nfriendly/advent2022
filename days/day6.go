package days

import (
	"errors"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/gammazero/deque"
)

type Day6Solver struct {
}

func (d Day6Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	requiredUniqueA := 4
	res, err := findStartOfPacket(puzzleInput, requiredUniqueA)
	if err != nil {
		return "", nil, err
	}
	canvas, err := buildWordDisplay(puzzleInput, res, requiredUniqueA)
	return strconv.Itoa(res), canvas, err
}

func (d Day6Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	requiredUniqueB := 14
	res, err := findStartOfPacket(puzzleInput, requiredUniqueB)
	if err != nil {
		return "", nil, err
	}
	canvas, err := buildWordDisplay(puzzleInput, res, requiredUniqueB)
	return strconv.Itoa(res), canvas, err
}

func findStartOfPacket(packet string, requiredUnique int) (int, error) {
	q := deque.New[rune](requiredUnique)

	for i, b := range packet {
		if q.Len() < requiredUnique {
			q.PushBack(b)
			continue
		}
		q.PopFront()
		q.PushBack(b)
		if uniqueContents(q) {
			return i + 1, nil
		}
	}

	return 0, errors.New("Failed to find a unique sequence of length " +
		strconv.Itoa(requiredUnique) + " in datastream: " + packet)
}

func uniqueContents(d *deque.Deque[rune]) bool {
	for i := 0; i < d.Len()-1; i++ {
		for j := i + 1; j < d.Len(); j++ {
			if d.At(i) == d.At(j) {
				return false
			}
		}
	}
	return true
}

func buildWordDisplay(word string, solutionIndex int, requiredUnique int) (fyne.CanvasObject, error) {
	t0 := canvas.NewText(word[:solutionIndex-requiredUnique], color.Black)
	t1 := canvas.NewText(word[solutionIndex-requiredUnique:solutionIndex], color.RGBA{255, 0, 0, 255})
	t2 := canvas.NewText(word[:solutionIndex], color.Black)
	hbox := container.NewHBox(t0, t1, t2)

	return container.NewHScroll(hbox), nil
}
