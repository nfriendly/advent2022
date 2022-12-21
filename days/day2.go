package days

import (
	"errors"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type Day2Solver struct {
}

func (d Day2Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	return calculateRockPaperScissorsScore(puzzleInput, true)
}

func (d Day2Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	return calculateRockPaperScissorsScore(puzzleInput, false)
}

func rockPaperScissorsStringToInt(moves []string) (int, int, error) {
	if len(moves) != 2 {
		return -1, -1, errors.New("Unexpected number of moves found (" + strconv.Itoa(len(moves)) + ")")
	}

	first := 0
	second := 0

	switch moves[0] {
	case "A":
		first = 1
	case "B":
		first = 2
	case "C":
		first = 3
	default:
		return -1, -1, errors.New("Unexpected first column value: " + moves[0])
	}

	switch moves[1] {
	case "X":
		second = 1
	case "Y":
		second = 2
	case "Z":
		second = 3
	default:
		return -1, -1, errors.New("Unexpected second column value: " + moves[1])
	}

	return first, second, nil
}

func calculateRockPaperScissorsScore(puzzleInput string, partA bool) (string, fyne.CanvasObject, error) {
	lines := strings.Split(strings.TrimSpace(puzzleInput), "\n")
	scores := make(plotter.XYs, len(lines)+1)
	scores = append(scores, plotter.XY{X: 0, Y: 0})
	for i, line := range lines {
		moves := strings.Split(strings.TrimSpace(line), " ")
		move0, move1, err := rockPaperScissorsStringToInt(moves)
		if err != nil {
			return "", nil, errors.New(err.Error() + " in line " + strconv.Itoa(i) + ": " + line)
		}
		score := 0
		if partA {
			// Points for your move
			score += move1
			// Points for the result
			switch {
			case move0 == move1:
				score += 3
			case move0-move1 == -1:
				score += 6
			case move0-move1 == 2:
				score += 6
			}
		} else {
			switch move1 {
			case 1: // lose
				switch move0 {
				case 1:
					score += 3
				case 2:
					score += 1
				case 3:
					score += 2
				}
			case 2: // tie
				score += 3 + move0
			case 3: // win
				score += 6 + (move0 % 3) + 1
			}
		}
		scores = append(scores, plotter.XY{
			X: float64(i + 1),
			Y: scores[len(scores)-1].Y + float64(score)})
	}

	plt := plot.New()
	plt.Title.Text = "Cumulative Rock Paper Sciessors Score"
	plt.X.Label.Text = "Round"
	plt.Y.Label.Text = "Score"

	l, err := plotter.NewLine(scores)

	if err != nil {
		return strconv.Itoa(int(scores[len(scores)-1].Y)), nil, err
	}
	plt.Add(l)

	name := "day2partB.png"
	if partA {
		name = "day2partA.png"
	}
	img, err := plotToImage(plt, name)
	return strconv.Itoa(int(scores[len(scores)-1].Y)), img, err
}
