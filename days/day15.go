package days

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Day15Solver struct {
}

func (d Day15Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	bs, err := buildBeacons(puzzleInput)
	if err != nil {
		return "", nil, err
	}

	// Hack to determine if it's a test or the puzzle
	if len(strings.Split(puzzleInput, "\n")) > 20 {
		sol, img, err := countImpossibleColumns(bs, 2000000, false)
		if err != nil {
			return "", img, err
		}
		return strconv.Itoa(sol), img, nil
	} else {
		sol, img, err := countImpossibleColumns(bs, 10, true)
		if err != nil {
			return "", img, err
		}
		return strconv.Itoa(sol), img, nil
	}
}

func (d Day15Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	// Assume the location is on the edge of a beacon's closest circle (aka it's manhattan distance + 1)
	// this should be valid because otherwise there would be multiple possible locations

	// for each beacon
	// for each point along it's border
	// if point is inside valid range x=[0,4000000] and y=[0,4000000] (or [0,20] for test)
	// check if too close to each beacon
	minX := 0
	maxX := 20
	minY := 0
	maxY := 20
	// Hack to determine if it's a test or the puzzle
	if len(strings.Split(puzzleInput, "\n")) > 20 {
		maxX = 4000000
		maxY = 4000000
	}

	bs, err := buildBeacons(puzzleInput)
	if err != nil {
		return "", nil, err
	}

	solX := -1
	solY := -1

	// Either I got lucky, or more likely:
	// Looks like you only need to examine one side of the border since we're assuming the point is pefectly surrounded
	// by 4 beacons. That means it would be found 4 times if we searched all 4 edges of each beacon boundary. neat!
	for _, b := range bs {
		// num/edge = manhattan+2
		// checl top left edge
		for i := 0; i < b.closestManhattanDistance+2; i++ {
			x := b.x - b.closestManhattanDistance - 1 + i
			y := b.y + i
			if x < minX || x > maxX || y < minY || y > maxY {
				continue
			}
			possible := true
			for _, b2 := range bs {
				mh := manhattanDistance(x, y, b2.x, b2.y)
				if mh < b2.closestManhattanDistance {
					possible = false
					break
				}
			}
			if possible {
				solX = x
				solY = y
				break
			}
		}
	}
	sol := solX*4000000 + solY

	return strconv.Itoa(sol), nil, nil
}

type beacon struct {
	x                        int
	y                        int
	closestX                 int
	closestY                 int
	closestManhattanDistance int
}

type beacons []beacon

func countImpossibleColumns(bs beacons, y int, test bool) (int, fyne.CanvasObject, error) {
	// Find min and max x possible using distances from beacon to closest
	minX := 100000000
	maxX := -100000000
	for _, b := range bs {
		if b.closestManhattanDistance < int(math.Abs(float64(y-b.y))) {
			// This one can't have any affect
			continue
		}
		x0 := b.x - b.closestManhattanDistance
		if x0 < minX {
			minX = x0
		}
		x1 := b.x + b.closestManhattanDistance
		if x1 > maxX {
			maxX = x1
		}
	}

	// Now count impossible locations in that x range for the specified y

	impossiblePositions := 0

	var sb strings.Builder
	if test {
		yStr := strconv.Itoa(y)
		sb.WriteString(strings.Repeat(" ", len(yStr)-minX))
		sb.WriteString("0")
		sb.WriteString(strings.Repeat(" ", maxX-minX-1))
		sb.WriteString("\n")
	}

	for x := minX; x < maxX; x++ {
		impossible := false
		isBeacon := false
		for _, b := range bs {
			if x == b.closestX && y == b.closestY {
				isBeacon = true
				break
			}
			if manhattanDistance(b.x, b.y, x, y) <= b.closestManhattanDistance {
				impossible = true
			}
		}
		if isBeacon {
			if test {
				sb.WriteString("B")
			}
		} else if impossible {
			impossiblePositions++
			if test {
				sb.WriteString("#")
			}
		} else {
			if test {
				sb.WriteString(".")
			}
		}
	}

	label := widget.NewLabel(sb.String())
	label.TextStyle.Monospace = true

	return impossiblePositions, label, nil
}

func manhattanDistance(x0, y0, x1, y1 int) int {
	return int(math.Abs(float64(x0-x1)) + math.Abs(float64(y0-y1)))
}

func buildBeacons(input string) (beacons, error) {
	re := regexp.MustCompile(
		`Sensor at x=(-?[0-9]+), y=(-?[0-9]+): closest beacon is at x=(-?[0-9]+), y=(-?[0-9]+)`)

	lines := strings.Split(input, "\n")
	bs := make(beacons, 0, len(lines))
	for _, line := range lines {
		parts := re.FindStringSubmatch(line)
		if len(parts) != 5 {
			return nil, errors.New("failed to parse line: " + line)
		}
		bX, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		bY, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, err
		}
		sX, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, err
		}
		sY, err := strconv.Atoi(parts[4])
		if err != nil {
			return nil, err
		}
		bs = append(bs, beacon{x: bX, y: bY, closestX: sX, closestY: sY, closestManhattanDistance: manhattanDistance(bX, bY, sX, sY)})
	}

	return bs, nil
}
