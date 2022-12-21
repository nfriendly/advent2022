package days

import (
	"errors"
	"fmt"
	"image/color"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func calculateCalorieCounts(input string) (plotter.Values, error) {
	lines := strings.Split(input, "\n")
	// TODO: precalculate size
	calories := make(plotter.Values, 10)
	current_cal := 0.0
	for _, line := range lines {
		if line == "" {
			calories = append(calories, current_cal)
			current_cal = 0.0
		} else {
			new_cal, err := strconv.ParseFloat(line, 64)
			if err != nil {
				failMsg := fmt.Sprint("Failed to parse line into float: ", line, ", err: ", err)
				fmt.Println(failMsg)
				return nil, errors.New(failMsg)
			}
			current_cal += new_cal
		}
	}
	if current_cal != 0.0 {
		calories = append(calories, current_cal)
	}

	return calories, nil
}

func findHighestCalorieCounts(input string, elves int) (int, *plot.Plot, error) {
	calories, err := calculateCalorieCounts(input)
	if err != nil {
		return -1, nil, err
	}

	sort.Float64s(calories)
	highestCalorieCount := 0
	for i := 0; i < elves; i++ {
		highestCalorieCount += int(calories[len(calories)-i-1])
	}

	w := vg.Points(1)
	bars, err := plotter.NewBarChart(calories, w)
	if err != nil {
		fmt.Println("Failed to create plot, err: ", err)
		return highestCalorieCount, nil, err
	}

	plt := plot.New()
	plt.Add(bars)

	for i := 0; i < elves; i++ {
		pts := make(plotter.XYs, 1)
		pts[0].X = float64(len(calories) - i - 1)
		pts[0].Y = calories[len(calories)-i-1]
		s, err := plotter.NewScatter(pts)
		if err != nil {
			fmt.Println("Failed to plot point ", i)
			continue
		}
		s.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
		plt.Add(s)
	}
	plt.Title.Text = "Total Calories Per Elf"
	plt.X.Label.Text = "Elf"
	plt.Y.Label.Text = "Total Calories"

	return highestCalorieCount, plt, nil
}

type Day1Solver struct {
}

func (d Day1Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	solution, sol_plt, err := findHighestCalorieCounts(puzzleInput, 1)
	if err != nil {
		return strconv.Itoa(solution), nil, err
	}
	img, err := plotToImage(sol_plt, "day1partA.png")
	return strconv.Itoa(solution), img, err
}

func (d Day1Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	solution, sol_plt, err := findHighestCalorieCounts(puzzleInput, 3)
	if err != nil {
		return strconv.Itoa(solution), nil, err
	}
	img, err := plotToImage(sol_plt, "day1partB.png")
	return strconv.Itoa(solution), img, err
}
