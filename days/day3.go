package days

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type Day3Solver struct {
}

func (d Day3Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	lines := strings.Split(strings.TrimSpace(puzzleInput), "\n")

	priorities := make(plotter.XYs, 0, len(lines)+1)
	priorities = append(priorities, plotter.XY{X: 0, Y: 0})

	for i, line := range lines {
		left, right, err := createRucksackCompartments(line)
		if err != nil {
			return "", nil, err
		}
		priority := 0
		for k := range left {
			if _, ok := right[k]; ok {
				priority = k
				break
			}
		}
		priorities = append(priorities, plotter.XY{
			X: float64(i + 1),
			Y: priorities[len(priorities)-1].Y + float64(priority)})
	}

	plt := plot.New()
	plt.Title.Text = "Cumulative Priority Score"
	plt.X.Label.Text = "Rucksack"
	plt.Y.Label.Text = "Total Priority"

	l, err := plotter.NewLine(priorities)

	if err != nil {
		return strconv.Itoa(int(priorities[len(priorities)-1].Y)), nil, err
	}
	plt.Add(l)

	img, err := plotToImage(plt, "day3partA.png")
	return strconv.Itoa(int(priorities[len(priorities)-1].Y)), img, err

}

func (d Day3Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	lines := strings.Split(strings.TrimSpace(puzzleInput), "\n")
	for i := range lines {
		r := []rune(lines[i])
		sort.Slice(r, func(i, j int) bool { return r[i] < r[j] })
		lines[i] = string(r)
	}

	// Build arrays of unique sorted characters
	sets := make([][]byte, len(lines))
	for i, line := range lines {
		byteArray := []byte(strings.TrimSpace(line))
		sets[i] = make([]byte, 0, len(byteArray)/2)
		sets[i] = append(sets[i], byteArray[0])
		for j := 1; j < len(byteArray); j++ {
			if byteArray[j] != byteArray[j-1] {
				sets[i] = append(sets[i], byteArray[j])
			}
		}
	}

	// Compare sets of three
	badges := make([]byte, 0, len(sets)/3)
	for trio := 0; trio < len(sets)-2; trio += 3 {
		var badge byte = 0
		for _, b0 := range sets[trio] {
			possibleMatch1 := sort.Search(len(sets[trio+1]),
				func(ii int) bool { return sets[trio+1][ii] >= b0 })
			if sets[trio+1][possibleMatch1] == b0 {
				possibleMatch2 := sort.Search(len(sets[trio+2]),
					func(ii int) bool { return sets[trio+2][ii] >= b0 })
				if sets[trio+2][possibleMatch2] == b0 {
					badge = b0
					break
				}
			}
		}
		if badge == 0 {
			return "", nil, errors.New("Failed to find badge for trio index " + strconv.Itoa(trio))
		}
		badges = append(badges, badge)
	}

	// Calculate priority score
	priorities := make(plotter.XYs, 0, len(lines)+1)
	priorities = append(priorities, plotter.XY{X: 0, Y: 0})

	for i, badge := range badges {
		priority := 0
		switch {
		case 65 <= badge && badge < 91:
			priority = int(badge) - 65 + 27 // gets 27 - 52
		case 97 <= badge && badge < 123:
			priority = int(badge) - 97 + 1 // gets 1 - 26
		default:
			return "", nil, errors.New("Invalid character " + string(badge) + " in trio index " + strconv.Itoa(i))
		}
		priorities = append(priorities, plotter.XY{
			X: float64(i + 1),
			Y: priorities[len(priorities)-1].Y + float64(priority)})
	}

	plt := plot.New()
	plt.Title.Text = "Cumulative Priority Score"
	plt.X.Label.Text = "Trio"
	plt.Y.Label.Text = "Total Priority"

	l, err := plotter.NewLine(priorities)

	if err != nil {
		return strconv.Itoa(int(priorities[len(priorities)-1].Y)), nil, err
	}
	plt.Add(l)

	img, err := plotToImage(plt, "day3partB.png")
	return strconv.Itoa(int(priorities[len(priorities)-1].Y)), img, err
}

func createRucksackCompartments(rucksackContents string) (map[int]bool, map[int]bool, error) {
	left := map[int]bool{}
	right := map[int]bool{}

	byteArray := []byte(strings.TrimSpace(rucksackContents))
	totalItems := len(byteArray)

	for i, b := range byteArray {
		val := 0
		switch {
		case 65 <= b && b < 91:
			val = int(b) - 65 + 27 // gets 27 - 52
		case 97 <= b && b < 123:
			val = int(b) - 97 + 1 // gets 1 - 26
		default:
			return nil, nil, errors.New("Invalid character " + string(b) + " in line: " + rucksackContents)
		}
		if i >= totalItems/2 {
			right[val] = true
		} else {
			left[val] = true
		}
	}

	return left, right, nil
}
