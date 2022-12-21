package days

import (
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type Day8Solver struct {
}

func (d Day8Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	tg := buildTreeGrid(puzzleInput)
	visible := buildTreeVisibilityGrid(tg)
	visibleCount := 0
	for _, row := range visible {
		for _, b := range row {
			if b {
				visibleCount++
			}
		}
	}

	content, err := visualizeTreeVisibility(tg, visible)

	return strconv.Itoa(visibleCount), content, err
}

func (d Day8Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	tg := buildTreeGrid(puzzleInput)
	x, y, score := findBestScenicScore(tg)

	highlight := make([][]bool, len(tg))
	for y := 0; y < len(tg); y++ {
		highlight[y] = make([]bool, len(tg[y]))
	}

	highlight[y][x] = true

	content, err := visualizeTreeVisibility(tg, highlight)

	return strconv.Itoa(score), content, err
}

type treeGrid [][]int

func buildTreeGrid(input string) treeGrid {
	lines := strings.Split(input, "\n")
	grid := make(treeGrid, len(lines))

	for y, l := range lines {
		grid[y] = make([]int, len(l))
		for x, r := range l {
			grid[y][x] = int(r - '0')
		}
	}
	return grid
}

func buildTreeVisibilityGrid(tg treeGrid) [][]bool {
	maxY := len(tg)
	maxX := len(tg[0])
	visible := make([][]bool, maxY)
	for y := 0; y < maxY; y++ {
		visible[y] = make([]bool, maxX)
	}

	// Explore from left
	for y := 0; y < maxY; y++ {
		maxHeight := -1
		for x := 0; x < maxX; x++ {
			if tg[y][x] > maxHeight {
				visible[y][x] = true
				maxHeight = tg[y][x]
			}
			// Can't get any higher so exit early
			if maxHeight == 9 {
				break
			}
		}
	}

	// Explore from right
	for y := 0; y < maxY; y++ {
		maxHeight := -1
		for x := maxX - 1; x >= 0; x-- {
			if tg[y][x] > maxHeight {
				visible[y][x] = true
				maxHeight = tg[y][x]
			}
			// Can't get any higher so exit early
			if maxHeight == 9 {
				break
			}
		}
	}

	// Explore from top
	for x := 0; x < maxX; x++ {
		maxHeight := -1
		for y := 0; y < maxY; y++ {
			if tg[y][x] > maxHeight {
				visible[y][x] = true
				maxHeight = tg[y][x]
			}
			// Can't get any higher so exit early
			if maxHeight == 9 {
				break
			}
		}
	}

	// Explore from bottom
	for x := 0; x < maxX; x++ {
		maxHeight := -1
		for y := maxY - 1; y >= 0; y-- {
			if tg[y][x] > maxHeight {
				visible[y][x] = true
				maxHeight = tg[y][x]
			}
			// Can't get any higher so exit early
			if maxHeight == 9 {
				break
			}
		}
	}

	return visible
}

func visualizeTreeVisibility(tg treeGrid, visible [][]bool) (fyne.CanvasObject, error) {
	grid := container.New(layout.NewGridLayout(len(tg[0])))
	for y := 0; y < len(tg); y++ {
		for x := 0; x < len(tg[0]); x++ {
			if visible[y][x] {
				grid.Add(canvas.NewText(strconv.Itoa(tg[y][x]), color.RGBA{255, 0, 0, 255}))
			} else {
				grid.Add(canvas.NewText(strconv.Itoa(tg[y][x]), color.Black))
			}

		}
	}

	return container.NewHScroll(grid), nil
}

func findBestScenicScore(tg treeGrid) (int, int, int) {
	bestScore := 0
	bestX := 0
	bestY := 0
	for y := 1; y < len(tg)-1; y++ {
		for x := 1; x < len(tg[y])-1; x++ {
			score := calculateScenicScore(tg, x, y)
			if score > bestScore {
				bestScore = score
				bestX = x
				bestY = y
			}
		}
	}
	return bestX, bestY, bestScore
}

func calculateScenicScore(tg treeGrid, x int, y int) int {
	var up, down, left, right int

	// Calculate left
	for xOffset := 1; x-xOffset >= 0; xOffset++ {
		left = xOffset
		if tg[y][x-xOffset] >= tg[y][x] {
			break
		}
	}
	// Calculate right
	for xOffset := 1; x+xOffset < len(tg[y]); xOffset++ {
		right = xOffset
		if tg[y][x+xOffset] >= tg[y][x] {
			break
		}
	}
	// Calculate up
	for yOffset := 1; y-yOffset >= 0; yOffset++ {
		up = yOffset
		if tg[y-yOffset][x] >= tg[y][x] {
			break
		}
	}
	// Calculate down
	for yOffset := 1; y+yOffset < len(tg); yOffset++ {
		down = yOffset
		if tg[y+yOffset][x] >= tg[y][x] {
			break
		}
	}

	return up * down * left * right
}
