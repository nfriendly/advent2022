package days

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

type Day4Solver struct {
}

type sectorAssignmentPair struct {
	start0, end0, start1, end1 int
}

func createOverlappingPlot(pairs []sectorAssignmentPair, overlapping []bool, numOverlapping int) (*plot.Plot, error) {
	// Build plot using x-error bars
	plt := plot.New()
	plt.Title.Text = "Overlapping Sector Assignment Pairs"
	plt.X.Label.Text = "Assigned sectors"
	plt.Y.Label.Text = "Pair Index"

	// build an XYs (implements XYer)
	// x holds pair index (0 to numPairs-1)
	// y holds start0
	// build two XErrors (implments XErrorer)
	// for first: Low holds 0, High holds end0-start0
	// for second: Low holds start1-start0, High holds end1-start0

	// actually have to hold min(start0, start1) because plot error bars takes abs value of the errors...

	// actually this is two sets of ErrorPoints (contains a XYs, an XErrors, and a YErrors)
	// then pass into plotter.NewXErrorsBars(), collect results into slice then pass to plt.Add(ps...)

	overlappingPoints0 := plotutil.ErrorPoints{
		XYs:     make(plotter.XYs, 0, numOverlapping),
		XErrors: make(plotter.XErrors, numOverlapping),
	}
	overlappingPoints1 := plotutil.ErrorPoints{
		XYs:     make(plotter.XYs, 0, numOverlapping),
		XErrors: make(plotter.XErrors, numOverlapping),
	}

	nonOverlappingPoints0 := plotutil.ErrorPoints{
		XYs:     make(plotter.XYs, 0, len(pairs)-numOverlapping),
		XErrors: make(plotter.XErrors, len(pairs)-numOverlapping),
	}

	nonOverlappingPoints1 := plotutil.ErrorPoints{
		XYs:     make(plotter.XYs, 0, len(pairs)-numOverlapping),
		XErrors: make(plotter.XErrors, len(pairs)-numOverlapping),
	}

	overlappingIdx := 0
	nonOverlappingIdx := 0
	for i, pair := range pairs {
		xy0 := plotter.XY{
			X: float64(pair.start0),
			Y: float64(i),
		}
		xy1 := plotter.XY{
			X: float64(pair.start1),
			Y: float64(i),
		}
		if overlapping[i] {
			overlappingPoints0.XYs = append(overlappingPoints0.XYs, xy0)
			overlappingPoints1.XYs = append(overlappingPoints1.XYs, xy1)

			overlappingPoints0.XErrors[overlappingIdx].Low = 0
			overlappingPoints0.XErrors[overlappingIdx].High = float64(pair.end0 - pair.start0)

			overlappingPoints1.XErrors[overlappingIdx].Low = 0
			overlappingPoints1.XErrors[overlappingIdx].High = float64(pair.end1 - pair.start1)

			overlappingIdx++
		} else {
			nonOverlappingPoints0.XYs = append(nonOverlappingPoints0.XYs, xy0)
			nonOverlappingPoints1.XYs = append(nonOverlappingPoints1.XYs, xy1)

			nonOverlappingPoints0.XErrors[nonOverlappingIdx].Low = 0
			nonOverlappingPoints0.XErrors[nonOverlappingIdx].High = float64(pair.end0 - pair.start0)

			nonOverlappingPoints1.XErrors[nonOverlappingIdx].Low = 0
			nonOverlappingPoints1.XErrors[nonOverlappingIdx].High = float64(pair.end1 - pair.start1)

			nonOverlappingIdx++
		}
	}
	e0, err := plotter.NewXErrorBars(overlappingPoints0)
	if err != nil {
		return nil, err
	}
	e0.Color = plotutil.Color(0)
	e1, err := plotter.NewXErrorBars(overlappingPoints1)
	if err != nil {
		return nil, err
	}
	e1.Color = plotutil.Color(1)
	e2, err := plotter.NewXErrorBars(nonOverlappingPoints0)
	if err != nil {
		return nil, err
	}
	e2.Color = plotutil.Color(2)
	e3, err := plotter.NewXErrorBars(nonOverlappingPoints1)
	if err != nil {
		return nil, err
	}
	e3.Color = plotutil.Color(3)

	plt.Add(e0, e1, e2, e3)
	return plt, nil
}

func (d Day4Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	pairs, err := parseSectorAssignmentList(puzzleInput)
	if err != nil {
		return "", nil, err
	}

	fullyOverlappingPairs := 0
	overlapping := make([]bool, len(pairs))
	for i, p := range pairs {
		if (p.start0 <= p.start1 && p.end0 >= p.end1) || (p.start1 <= p.start0 && p.end1 >= p.end0) {
			fullyOverlappingPairs++
			overlapping[i] = true
		}
	}

	plt, err := createOverlappingPlot(pairs, overlapping, fullyOverlappingPairs)
	if err != nil {
		return strconv.Itoa(fullyOverlappingPairs), nil, err
	}

	img, err := plotToImage(plt, "day4partA.png")

	return strconv.Itoa(fullyOverlappingPairs), img, err
}

func (d Day4Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	pairs, err := parseSectorAssignmentList(puzzleInput)
	if err != nil {
		return "", nil, err
	}

	overlappingPairs := 0
	overlapping := make([]bool, len(pairs))
	for i, p := range pairs {
		if !((p.start0 > p.end1) || (p.start1 > p.end0)) {
			overlappingPairs++
			overlapping[i] = true
		}
	}

	plt, err := createOverlappingPlot(pairs, overlapping, overlappingPairs)
	if err != nil {
		return strconv.Itoa(overlappingPairs), nil, err
	}

	img, err := plotToImage(plt, "day4partA.png")

	return strconv.Itoa(overlappingPairs), img, err
}

func parseSectorAssignmentList(puzzleInput string) ([]sectorAssignmentPair, error) {
	lines := strings.Split(strings.TrimSpace(puzzleInput), "\n")
	result := make([]sectorAssignmentPair, len(lines))

	re := regexp.MustCompile(`([0-9]+)-([0-9]+),([0-9]+)-([0-9]+)`)
	for i, line := range lines {
		res := re.FindSubmatch([]byte(line))
		if len(res) != 5 {
			return result, errors.New("Failed to parse line:" + line)
		}
		var err error = nil
		result[i].start0, err = strconv.Atoi(string(res[1]))
		if err != nil {
			return result, errors.New("Failed to convert " + string(res[1]) + " to an interger from line:" + line)
		}
		result[i].end0, err = strconv.Atoi(string(res[2]))
		if err != nil {
			return result, errors.New("Failed to convert " + string(res[2]) + " to an interger from line:" + line)
		}
		result[i].start1, err = strconv.Atoi(string(res[3]))
		if err != nil {
			return result, errors.New("Failed to convert " + string(res[3]) + " to an interger from line:" + line)
		}
		result[i].end1, err = strconv.Atoi(string(res[4]))
		if err != nil {
			return result, errors.New("Failed to convert " + string(res[4]) + " to an interger from line:" + line)
		}
		if result[i].start0 > result[i].end0 || result[i].start1 > result[i].end1 {
			return result, errors.New("Invalid ranges for line:" + line)
		}
	}

	return result, nil
}
