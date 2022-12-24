package days

import (
	"errors"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Day17Solver struct {
}

func (d Day17Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	jets, err := buildJetDirections(puzzleInput)
	if err != nil {
		return "", nil, err
	}

	chamber := make(chamberMap, 1, 2022*2)
	chamber[0] = chamberLevel{'|', '-', '-', '-', '-', '-', '-', '-', '|'}
	chamber = append(chamber, chamberLevel{'|', '.', '.', '.', '.', '.', '.', '.', '|'})
	chamber = append(chamber, chamberLevel{'|', '.', '.', '.', '.', '.', '.', '.', '|'})
	chamber = append(chamber, chamberLevel{'|', '.', '.', '.', '.', '.', '.', '.', '|'})

	// Drop 2022 rocks
	rocksToDrop := 2022
	// rocksToDrop := 10
	currentHeight := 0
	jetIdx := 0
	for rockIdx := 0; rockIdx < rocksToDrop; rockIdx++ {
		// Create rock
		rock := buildFallingRock(rockIdx)
		rock.position.x = 3
		rock.position.y = currentHeight + 4

		// Extend chamber if needed
		for len(chamber) < rock.position.y+4 {
			chamber = append(chamber, chamberLevel{'|', '.', '.', '.', '.', '.', '.', '.', '|'})
		}

		// Drop it until it hits
		for {
			// Push
			if jets[jetIdx] {
				// push right
				if isRightFree(rock, chamber) {
					rock.position.x++
				}
			} else {
				// push left
				if isLeftFree(rock, chamber) {
					rock.position.x--
				}
			}
			jetIdx++
			jetIdx %= len(jets)

			// Drop (break if can't)
			// if rockCanFall(rock, chamber[rock.position.y-1]) {
			if rockCanFall(rock, chamber) {
				rock.position.y--
			} else {
				break
			}
		}

		// Update chamber and currentHeight
		for yOffset, row := range rock.shape {
			for xOffset, pos := range row {
				if pos {
					chamber[rock.position.y+yOffset][rock.position.x+xOffset] = '#'
					if rock.position.y+yOffset > currentHeight {
						currentHeight = rock.position.y + yOffset
					}
				}
			}
		}
	}

	// Find height
	maxHeight := 0
	for h := len(chamber) - 1; h > 0; h-- {
		for _, pos := range chamber[h] {
			if pos == '#' {
				maxHeight = h
				break
			}
		}
		if maxHeight > 0 {
			break
		}
	}

	// Make image
	img := visualizeChamber(chamber)

	return strconv.Itoa(maxHeight), img, nil
}

func (d Day17Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	return "", nil, nil
}

type fallingRock struct {
	// origin is at spawn location (aka bottom left)
	// y pos up, x pos right
	shape [4][4]bool
	// convience value to avoid scanning shape
	width int
	// Floor is y=0, right wall is at x=0
	position point
}

type chamberLevel [9]rune
type chamberMap []chamberLevel

func buildFallingRock(index int) fallingRock {
	switch index % 5 {
	case 0:
		return fallingRock{shape: [4][4]bool{
			{true, true, true, true},
			{false, false, false, false},
			{false, false, false, false},
			{false, false, false, false}},
			width: 4}
	case 1:
		return fallingRock{shape: [4][4]bool{
			{false, true, false, false},
			{true, true, true, false},
			{false, true, false, false},
			{false, false, false, false}},
			width: 3}
	case 2:
		return fallingRock{shape: [4][4]bool{
			{true, true, true, false},
			{false, false, true, false},
			{false, false, true, false},
			{false, false, false, false}},
			width: 3}
	case 3:
		return fallingRock{shape: [4][4]bool{
			{true, false, false, false},
			{true, false, false, false},
			{true, false, false, false},
			{true, false, false, false}},
			width: 1}
	case 4:
		return fallingRock{shape: [4][4]bool{
			{true, true, false, false},
			{true, true, false, false},
			{false, false, false, false},
			{false, false, false, false}},
			width: 2}
	default:
		panic("got an unexpected return from the mod operator")
	}
}

// Returns true for right, false for left
func buildJetDirections(input string) ([]bool, error) {
	res := make([]bool, len(input))

	for i, dir := range input {
		if dir == '>' {
			res[i] = true
		} else if dir != '<' {
			return nil, errors.New("got unexpected rune at position " + strconv.Itoa(i) + ": " + string(dir))
		}
	}
	return res, nil
}

// func rockCanFall(rock fallingRock, levelBelow chamberLevel) bool {
func rockCanFall(rock fallingRock, chamber chamberMap) bool {
	for xOffset := 3; xOffset >= 0; xOffset-- {
		// Find the lowest rock point in each column
		for yOffset := 0; yOffset < 4; yOffset++ {
			if rock.shape[yOffset][xOffset] {
				if chamber[rock.position.y+yOffset-1][rock.position.x+xOffset] == '.' {
					continue
				} else {
					return false
				}
			}
		}
	}
	// for i := 0; i < 4; i++ {
	// 	if rock.shape[0][i] && levelBelow[rock.position.x+i] != '.' {
	// 		return false
	// 	}
	// }
	return true
}

func isRightFree(rock fallingRock, chamber chamberMap) bool {
	for yOffset := 0; yOffset < 4; yOffset++ {
		for xOffset := 3; xOffset >= 0; xOffset-- {
			// Find farthest right rock point in row and then
			if rock.shape[yOffset][xOffset] {
				if chamber[rock.position.y+yOffset][rock.position.x+xOffset+1] == '.' {
					// Done with this row, move on to next on
					break
				} else {
					// Collision on this row
					return false
				}
			}
		}
	}
	return true
}

func isLeftFree(rock fallingRock, chamber chamberMap) bool {
	for yOffset := 0; yOffset < 4; yOffset++ {
		for xOffset := 0; xOffset < 4; xOffset++ {
			// Find farthest left rock point in row and then
			if rock.shape[yOffset][xOffset] {
				if chamber[rock.position.y+yOffset][rock.position.x+xOffset-1] == '.' {
					// Done with this row, move on to next on
					break
				} else {
					// Collision on this row
					return false
				}
			}
		}
	}
	return true
}

func visualizeChamber(chamber chamberMap) fyne.CanvasObject {
	label := widget.NewLabel(chamber.String())
	label.TextStyle.Monospace = true
	return container.NewHScroll(label)
}

func (c chamberMap) String() string {
	var sb strings.Builder
	for i := len(c) - 1; i >= 0; i-- {
		sb.WriteString(strconv.Itoa(i))
		for _, r := range c[i] {
			sb.WriteRune(r)
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}
