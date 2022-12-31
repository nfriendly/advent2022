package days

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Day22Solver struct {
}

func (d Day22Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	mm, err := buildMonkeyMap(puzzleInput)
	if err != nil {
		return "", nil, err
	}
	err = mm.follow()
	if err != nil {
		return "", nil, err
	}
	sol := 1000*(mm.currentPosition.y+1) + 4*(mm.currentPosition.x+1) + int(mm.currentOrientation)

	img := mm.Visualize()
	return strconv.Itoa(sol), img, err
}

func (d Day22Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	return "", nil, nil
}

type monkeyMapOrientation int

const (
	monkeyRight monkeyMapOrientation = 0
	monkeyDown  monkeyMapOrientation = 1
	monkeyLeft  monkeyMapOrientation = 2
	monkeyUp    monkeyMapOrientation = 3
)

type monkeyMapRow struct {
	minValid, maxValid int
	// slice from [0, maxValid-minValid+1], use x-minValid to index in
	occuiped        []bool
	lastOrientation []monkeyMapOrientation
}

type monkeyMapDirection interface {
	isMonkeyMapDirection()
}

// number of spaces to move forward
type monkeyMapDirectionMove int

// true for R/CW, false for left/CCW
type monkeyMapDirectionTurn bool

func (monkeyMapDirectionMove) isMonkeyMapDirection() {}
func (monkeyMapDirectionTurn) isMonkeyMapDirection() {}

type monkeyMap struct {
	rows               []monkeyMapRow
	currentPosition    point
	currentOrientation monkeyMapOrientation
	directions         []monkeyMapDirection
}

func buildMonkeyMap(input string) (monkeyMap, error) {
	lines := strings.Split(input, "\n")
	mm := monkeyMap{
		rows:               make([]monkeyMapRow, 0, len(lines)-2),
		currentOrientation: monkeyRight,
		directions:         make([]monkeyMapDirection, 0),
	}

	mapComplete := false
	re := regexp.MustCompile("[0-9]+|[RL]")
	for _, line := range lines {
		if mapComplete {
			parts := re.FindAllString(line, -1)
			if len(parts) == 0 {
				return mm, errors.New("failed to parse directions")
			}
			for _, part := range parts {
				if part == "R" {
					mm.directions = append(mm.directions, monkeyMapDirectionTurn(true))
				} else if part == "L" {
					mm.directions = append(mm.directions, monkeyMapDirectionTurn(false))
				} else {
					val, err := strconv.Atoi(part)
					if err != nil {
						return mm, err
					}
					mm.directions = append(mm.directions, monkeyMapDirectionMove(val))
				}
			}
			break
		}
		if line == "" {
			mapComplete = true
			continue
		}
		mm.rows = append(mm.rows, buildMonkeyMapRow(line))
	}

	mm.currentPosition.x = mm.rows[0].minValid
	mm.currentPosition.y = 0

	return mm, nil
}

func buildMonkeyMapRow(input string) monkeyMapRow {
	r := monkeyMapRow{}
	r.minValid = strings.Count(input, " ")
	r.maxValid = len(input) - 1
	r.occuiped = make([]bool, r.maxValid-r.minValid+1)
	r.lastOrientation = make([]monkeyMapOrientation, r.maxValid-r.minValid+1)
	for i := r.minValid; i <= r.maxValid; i++ {
		if input[i] == '.' {
			r.occuiped[i-r.minValid] = false
		} else if input[i] == '#' {
			r.occuiped[i-r.minValid] = true
		} else {
			panic("unexpected character " + string(input[i]) +
				" at position " + strconv.Itoa(i) +
				" in line " + input)
		}
		r.lastOrientation[i-r.minValid] = -1
	}
	return r
}

func (m *monkeyMap) follow() error {
	for i := range m.directions {
		switch t := m.directions[i].(type) {
		case monkeyMapDirectionMove:
			m.MoveForward(t)
		case monkeyMapDirectionTurn:
			m.Rotate(t)
		default:
			return errors.New("unexpected direction type at index " + strconv.Itoa(i))
		}
	}
	return nil
}

func (m *monkeyMap) Visualize() fyne.CanvasObject {
	label := widget.NewLabel(m.String())
	label.TextStyle.Monospace = true

	return label
}

func (m *monkeyMap) MoveForward(amount monkeyMapDirectionMove) {
	m.setLastOrientation(m.currentPosition.x, m.currentPosition.y, m.currentOrientation)
	switch m.currentOrientation {
	case monkeyRight:
		r := &m.rows[m.currentPosition.y]
		for i := 0; i < int(amount); i++ {
			if m.currentPosition.x+1 > r.maxValid {
				// attempt to roll over
				if r.isOccupied(r.minValid) {
					break
				}
				m.currentPosition.x = r.minValid
				m.setLastOrientation(m.currentPosition.x, m.currentPosition.y, m.currentOrientation)
			} else {
				if r.isOccupied(m.currentPosition.x + 1) {
					break
				}
				m.currentPosition.x++
				m.setLastOrientation(m.currentPosition.x, m.currentPosition.y, m.currentOrientation)
			}
		}
	case monkeyLeft:
		r := &m.rows[m.currentPosition.y]
		for i := 0; i < int(amount); i++ {
			if m.currentPosition.x-1 < r.minValid {
				// attempt to roll over
				if r.isOccupied(r.maxValid) {
					break
				}
				m.currentPosition.x = r.maxValid
				m.setLastOrientation(m.currentPosition.x, m.currentPosition.y, m.currentOrientation)
			} else {
				if r.isOccupied(m.currentPosition.x - 1) {
					break
				}
				m.currentPosition.x--
				m.setLastOrientation(m.currentPosition.x, m.currentPosition.y, m.currentOrientation)
			}
		}
	case monkeyUp:
		x := m.currentPosition.x
		for i := 0; i < int(amount); i++ {
			if m.currentPosition.y-1 < 0 || !m.rows[m.currentPosition.y-1].contains(x) {
				// find top (highest y) row containing x
				newY := len(m.rows) - 1
				for ; !m.rows[newY].contains(x); newY-- {
					if newY+1 < 0 {
						panic("invalid map, infite up rollover")
					}
				}
				if m.rows[newY].isOccupied(x) {
					break
				}
				m.currentPosition.y = newY
				m.setLastOrientation(m.currentPosition.x, m.currentPosition.y, m.currentOrientation)
			} else {
				if m.rows[m.currentPosition.y-1].isOccupied(x) {
					break
				}
				m.currentPosition.y--
				m.setLastOrientation(m.currentPosition.x, m.currentPosition.y, m.currentOrientation)
			}
		}
	case monkeyDown:
		x := m.currentPosition.x
		for i := 0; i < int(amount); i++ {
			if m.currentPosition.y+1 >= len(m.rows) || !m.rows[m.currentPosition.y+1].contains(x) {
				// find top (lowest y) row containing x
				newY := 0
				for ; !m.rows[newY].contains(x); newY++ {
					if newY+1 >= len(m.rows) {
						panic("invalid map, infite down rollover")
					}
				}
				if m.rows[newY].isOccupied(x) {
					break
				}
				m.currentPosition.y = newY
				m.setLastOrientation(m.currentPosition.x, m.currentPosition.y, m.currentOrientation)
			} else {
				if m.rows[m.currentPosition.y+1].isOccupied(x) {
					break
				}
				m.currentPosition.y++
				m.setLastOrientation(m.currentPosition.x, m.currentPosition.y, m.currentOrientation)
			}
		}
	default:
		panic("invallid monkdy orientation value: " + strconv.Itoa(int(m.currentOrientation)))
	}
}

func (m *monkeyMap) Rotate(right monkeyMapDirectionTurn) {
	if right {
		m.currentOrientation++
	} else {
		m.currentOrientation--
	}

	// Assumes we never roll over by more than 1
	if m.currentOrientation > monkeyUp {
		m.currentOrientation = monkeyRight
	}
	if m.currentOrientation < monkeyRight {
		m.currentOrientation = monkeyUp
	}
	m.setLastOrientation(m.currentPosition.x, m.currentPosition.y, m.currentOrientation)
}

func (m *monkeyMap) setLastOrientation(x int, y int, orientation monkeyMapOrientation) {
	if y >= len(m.rows) || !m.rows[y].contains(x) {
		panic("can't set orientation of invalid position (" + strconv.Itoa(x) + ", " + strconv.Itoa(y) + ")")
	}
	m.rows[y].lastOrientation[x-m.rows[y].minValid] = orientation
}

func (m *monkeyMap) String() string {
	var sb strings.Builder
	for y := range m.rows {
		if m.rows[y].minValid > 0 {
			sb.WriteString(strings.Repeat(" ", m.rows[y].minValid-1))
		}
		for i := range m.rows[y].lastOrientation {
			if m.rows[y].occuiped[i] {
				sb.WriteRune('#')
			} else {
				switch m.rows[y].lastOrientation[i] {
				case monkeyRight:
					sb.WriteRune('>')
				case monkeyDown:
					sb.WriteRune('v')
				case monkeyLeft:
					sb.WriteRune('<')
				case monkeyUp:
					sb.WriteRune('^')
				default:
					sb.WriteRune('.')
				}
			}
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (m *monkeyMapRow) contains(x int) bool {
	return m.minValid <= x && x <= m.maxValid
}

func (m *monkeyMapRow) isOccupied(x int) bool {
	if !m.contains(x) {
		panic("attemped to check  occupation of noncontained value " + strconv.Itoa(x))
	}
	return m.occuiped[x-m.minValid]
}
