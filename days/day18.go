package days

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
)

type Day18Solver struct {
}

func (d Day18Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	droplet, err := buildLavaDroplet(puzzleInput)
	if err != nil {
		return "", nil, err
	}
	return strconv.Itoa(droplet.surfaceArea), nil, nil
}

func (d Day18Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	droplet, err := buildLavaDroplet(puzzleInput)
	if err != nil {
		return "", nil, err
	}
	intEdges := droplet.countInteriorEdges()
	res := droplet.surfaceArea - intEdges
	return strconv.Itoa(res), nil, nil
}

type threePoint struct {
	x, y, z int
}

type lavaDroplet struct {
	scannedPoints map[threePoint]bool
	surfaceArea   int
}

func buildLavaDroplet(input string) (lavaDroplet, error) {
	lines := strings.Split(input, "\n")
	ld := lavaDroplet{
		scannedPoints: make(map[threePoint]bool, len(lines)),
		surfaceArea:   0,
	}
	re := regexp.MustCompile(`([0-9]+),([0-9]+),([0-9]+)`)

	for _, line := range lines {
		parts := re.FindStringSubmatch(line)
		if len(parts) != 4 {
			return ld, errors.New("failed to parse line: " + line)
		}
		var p threePoint
		var err error
		p.x, err = strconv.Atoi(parts[1])
		if err != nil {
			return ld, err
		}
		p.y, err = strconv.Atoi(parts[2])
		if err != nil {
			return ld, err
		}
		p.z, err = strconv.Atoi(parts[3])
		if err != nil {
			return ld, err
		}

		neighbors := ld.neighborsPresent(p)
		ld.scannedPoints[p] = true
		ld.surfaceArea += 6 - 2*neighbors
	}
	return ld, nil
}

func (ld *lavaDroplet) neighborsPresent(p threePoint) int {
	neighbors := 0
	for _, pn := range p.neighbors() {
		if ld.scannedPoints[pn] {
			neighbors++
		}
	}

	return neighbors
}

func (ld *lavaDroplet) countInteriorEdges() int {
	// Not sure how to do this
	return 0
}

func (p threePoint) neighbors() []threePoint {
	return []threePoint{
		{x: p.x + 1, y: p.y, z: p.z},
		{x: p.x - 1, y: p.y, z: p.z},
		{x: p.x, y: p.y + 1, z: p.z},
		{x: p.x, y: p.y - 1, z: p.z},
		{x: p.x, y: p.y, z: p.z + 1},
		{x: p.x, y: p.y, z: p.z - 1},
	}
}
