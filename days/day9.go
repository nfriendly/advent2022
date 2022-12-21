package days

import (
	"errors"
	"image/color"
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type Day9Solver struct {
}

func (d Day9Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	commands := strings.Split(strings.TrimSpace(puzzleInput), "\n")

	r := ropeState{}
	headHistory := make(plotter.XYs, 0, len(commands)*2)
	tailHistory := make(plotter.XYs, 0, len(commands)*2)
	headHistory = append(headHistory, plotter.XY{X: 0, Y: 0})
	tailHistory = append(tailHistory, plotter.XY{X: 0, Y: 0})

	tailPosSet := make(map[point]bool)

	for _, c := range commands {
		parts := strings.Split(c, " ")
		if len(parts) != 2 {
			return "", nil, errors.New("failed to parse line: " + c)
		}
		steps, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", nil, err
		}
		for s := 0; s < steps; s++ {
			if err := r.MoveHead(parts[0]); err != nil {
				return "", nil, err
			}
			headHistory = append(headHistory, plotter.XY{X: float64(r.head.x), Y: float64(r.head.y)})
			tailHistory = append(tailHistory, plotter.XY{X: float64(r.tail.x), Y: float64(r.tail.y)})
			tailPosSet[point{x: r.tail.x, y: r.tail.y}] = true
		}
	}

	plt := plot.New()
	plt.Title.Text = "Rope Position"
	plt.X.Label.Text = "X"
	plt.Y.Label.Text = "Y"

	l0, err := plotter.NewLine(headHistory)
	if err != nil {
		return strconv.Itoa(len(tailPosSet)), nil, err
	}
	l0.Color = color.RGBA{R: 255, A: 255}
	l1, err := plotter.NewLine(tailHistory)
	if err != nil {
		return strconv.Itoa(len(tailPosSet)), nil, err
	}
	l1.Color = color.RGBA{B: 255, A: 255}

	plt.Add(l0)
	plt.Add(l1)
	plt.Legend.Add("Head", l0)
	plt.Legend.Add("Tail", l1)

	img, err := plotToImage(plt, "day9partA.png")
	return strconv.Itoa(len(tailPosSet)), img, err
}

func (d Day9Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	commands := strings.Split(strings.TrimSpace(puzzleInput), "\n")

	r := longRopeState{}

	tailPosSet := make(map[point]bool)
	tailHistory := make(plotter.XYs, 0, len(commands)*2)
	tailHistory = append(tailHistory, plotter.XY{X: 0, Y: 0})

	for _, c := range commands {
		parts := strings.Split(c, " ")
		if len(parts) != 2 {
			return "", nil, errors.New("failed to parse line: " + c)
		}
		steps, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", nil, err
		}
		for s := 0; s < steps; s++ {
			if err := r.MoveHead(parts[0]); err != nil {
				return "", nil, err
			}
			tail := r.knots[len(r.knots)-1]
			tailPosSet[tail] = true
			tailHistory = append(tailHistory, plotter.XY{X: float64(tail.x), Y: float64(tail.y)})
		}
	}

	finalPosition := make(plotter.XYs, len(r.knots))
	for i, k := range r.knots {
		finalPosition[i].X = float64(k.x)
		finalPosition[i].Y = float64(k.y)
	}
	l0, err := plotter.NewLine(finalPosition)
	if err != nil {
		return strconv.Itoa(len(tailPosSet)), nil, err
	}

	l1, err := plotter.NewLine(tailHistory)
	if err != nil {
		return strconv.Itoa(len(tailPosSet)), nil, err
	}
	l1.Color = color.RGBA{B: 255, A: 255}

	plt := plot.New()
	plt.Title.Text = "Final Rope Position"
	plt.X.Label.Text = "X"
	plt.Y.Label.Text = "Y"

	plt.Add(l0)
	plt.Add(l1)
	plt.Legend.Add("Final Position", l0)
	plt.Legend.Add("Tail Path", l1)
	img, err := plotToImage(plt, "day9partB.png")
	return strconv.Itoa(len(tailPosSet)), img, err
}

// Up, right = positive, down, left = negative
type ropeState struct {
	// headX, headY, tailX, tailY int
	head, tail point
}

type longRopeState struct {
	knots [10]point
}

func (r *ropeState) MoveHead(direction string) error {
	switch direction {
	case "U":
		r.head.y++
	case "D":
		r.head.y--
	case "L":
		r.head.x--
	case "R":
		r.head.x++
	default:
		return errors.New("invalid direction command: " + direction)
	}
	if !areAdjacent(r.head, r.tail) {
		updateTrailing(r.head, &r.tail)
	}

	return nil
}

func (r *longRopeState) MoveHead(direction string) error {
	switch direction {
	case "U":
		r.knots[0].y++
	case "D":
		r.knots[0].y--
	case "L":
		r.knots[0].x--
	case "R":
		r.knots[0].x++
	default:
		return errors.New("invalid direction command: " + direction)
	}

	for i := 0; i < len(r.knots)-1; i++ {
		if areAdjacent(r.knots[i], r.knots[i+1]) {
			break
		}
		updateTrailing(r.knots[i], &r.knots[i+1])
	}

	return nil
}

func areAdjacent(p0 point, p1 point) bool {
	if math.Abs(float64(p0.x-p1.x)) > 1 || math.Abs(float64(p0.y-p1.y)) > 1 {
		return false
	}
	return true
}

func updateTrailing(p0 point, p1 *point) {
	if p0.x > p1.x {
		p1.x++
	} else if p0.x < p1.x {
		p1.x--
	}
	if p0.y > p1.y {
		p1.y++
	} else if p0.y < p1.y {
		p1.y--
	}
}
