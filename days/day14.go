package days

import (
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Day14Solver struct {
}

func (d Day14Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	cm, err := buildCaveMap(puzzleInput)
	if err != nil {
		return "", nil, err
	}

	sandDropped := cm.fillCaveMapWithSand(500, 0)
	img, err := visualizeCaveMap(cm)

	return strconv.Itoa(sandDropped), img, err
}

func (d Day14Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	cm, err := buildCaveMap(puzzleInput)
	if err != nil {
		return "", nil, err
	}

	cm.obstructions[cm.maxY+2] = append(cm.obstructions[cm.maxY+2], caveObstruction{startX: -10000, endX: 10000})

	sandDropped := cm.fillCaveMapWithSand(500, 0)
	img, err := visualizeCaveMap(cm)

	return strconv.Itoa(sandDropped), img, err
}

type caveObstruction struct {
	startX int
	endX   int
	isSand bool
}

type caveObstructions []caveObstruction

type caveMap struct {
	obstructions           map[int]caveObstructions
	minX, maxX, minY, maxY int
}

func (obs caveObstructions) Len() int           { return len(obs) }
func (obs caveObstructions) Swap(i, j int)      { obs[i], obs[j] = obs[j], obs[i] }
func (obs caveObstructions) Less(i, j int) bool { return obs[i].startX < obs[j].startX }

type fallResult int

const (
	free fallResult = iota
	leftFree
	rightFree
	blocked
)

func buildCaveMap(input string) (*caveMap, error) {
	cm := caveMap{
		obstructions: make(map[int]caveObstructions),
		minX:         10000,
		maxX:         -10000,
		minY:         10000,
		maxY:         -10000,
	}

	lines := strings.Split(input, "\n")

	re := regexp.MustCompile(`([0-9]+),([0-9]+)`)
	for _, line := range lines {
		points := make([]point, 0)
		res := re.FindAllStringSubmatch(line, -1)
		for _, r := range res {
			if len(r) != 3 {
				return nil, errors.New("failed to parse line: " + line + ". got regexp result: " + strings.Join(r, "\\"))
			}
			x, err := strconv.Atoi(r[1])
			if err != nil {
				return nil, err
			}
			y, err := strconv.Atoi(r[2])
			if err != nil {
				return nil, err
			}
			if x > cm.maxX {
				cm.maxX = x
			}
			if x < cm.minX {
				cm.minX = x
			}
			if y > cm.maxY {
				cm.maxY = y
			}
			if y < cm.minY {
				cm.minY = y
			}
			points = append(points, point{x: x, y: y})
		}

		for i := 1; i < len(points); i++ {
			if points[i].x == points[i-1].x {
				// horizontal line
				var startY, endY int
				if points[i].y < points[i-1].y {
					startY = points[i].y
					endY = points[i-1].y
				} else {
					startY = points[i-1].y
					endY = points[i].y
				}
				for y := startY; y <= endY; y++ {
					cm.obstructions[y] = append(cm.obstructions[y],
						caveObstruction{startX: points[i].x, endX: points[i].x})
				}
			} else if points[i].y == points[i-1].y {
				// vertical line
				var startX, endX int
				if points[i].x < points[i-1].x {
					startX = points[i].x
					endX = points[i-1].x
				} else {
					startX = points[i-1].x
					endX = points[i].x
				}
				cm.obstructions[points[i].y] = append(cm.obstructions[points[i].y],
					caveObstruction{startX: startX, endX: endX})
			} else {
				// single point (shouldn't happen)
				cm.obstructions[points[i].y] = append(cm.obstructions[points[i].y],
					caveObstruction{startX: points[i].x, endX: points[i].x})
			}
		}
	}

	return &cm, nil
}

func (cm *caveMap) fillCaveMapWithSand(startingX int, startingY int) int {
	restingSand := 0
	// loop over sand units
	finished := false
	for !finished {
		x := startingX
		y := startingY

		// loop over attempts until blocked
		for {
			res := attemptSandFall(x, y, cm)
			if res == blocked {
				// Add sand
				cm.obstructions[y] = append(cm.obstructions[y],
					caveObstruction{startX: x, endX: x, isSand: true})
				restingSand++
				if x == startingX && y == startingY {
					finished = true
				}
				break
			} else if res == leftFree {
				// move down left
				x--
				y++
			} else if res == rightFree {
				// move down right
				x++
				y++
			} else if res == free {
				// move down
				y++
				if y > cm.maxY+2 {
					finished = true
					break
				}
			}
		}
	}
	return restingSand
}

func attemptSandFall(x int, y int, cm *caveMap) fallResult {
	obs := cm.obstructions[y+1]
	underBlocked := false
	leftBlocked := false
	rightBlocked := false
	for _, ob := range obs {
		if ob.startX <= x && ob.endX >= x {
			underBlocked = true
		}
		if ob.startX <= (x-1) && ob.endX >= (x-1) {
			leftBlocked = true
		}
		if ob.startX <= (x+1) && ob.endX >= (x+1) {
			rightBlocked = true
		}

		if underBlocked && leftBlocked && rightBlocked {
			return blocked
		}
	}
	if !underBlocked {
		return free
	} else if !leftBlocked {
		return leftFree
	} else if !rightBlocked {
		return rightFree
	}
	return blocked
}

func visualizeCaveMap(cm *caveMap) (fyne.CanvasObject, error) {
	var sb strings.Builder
	for y := cm.minY; y <= cm.maxY+2; y++ {
		obs := cm.obstructions[y]
		sort.Sort(obs)
		currObIdx := 0
		for x := cm.minX; x <= cm.maxX; x++ {
			if len(obs) > 0 {
				for currObIdx < len(obs)-1 && obs[currObIdx].endX < x {
					currObIdx++
				}

				if obs[currObIdx].startX <= x && x <= obs[currObIdx].endX {
					if obs[currObIdx].isSand {
						sb.WriteString("o")
					} else {
						sb.WriteString("#")
					}
				} else {
					sb.WriteString(".")
				}
			} else {
				sb.WriteString(".")
			}
		}
		sb.WriteString("\n")
	}

	label := widget.NewLabel(sb.String())
	label.TextStyle.Monospace = true

	return label, nil
}
