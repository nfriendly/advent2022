package days

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Day17Solver struct {
}

func (d Day17Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	chamber, err := buildChamber(puzzleInput)
	if err != nil {
		return "", nil, err
	}

	// Drop 2022 rocks
	rocksToDrop := 2022
	// rocksToDrop := 10
	for rockIdx := 0; rockIdx < rocksToDrop; rockIdx++ {
		// Create rock
		rock := buildFallingRock(rockIdx)
		chamber.dropRock(&rock)
	}

	// Make image
	img := visualizeChamber(chamber)

	return strconv.Itoa(chamber.maxHeight), img, nil
}

func (d Day17Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	// TODO: find when pattern of rocks and movement loops
	// calculate height of first pass (might be different?)
	// calculate height of second pass
	// multiply by number of repeats needed
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

type chamberMap struct {
	// Represents the chamber full of rocks
	m []chamberLevel
	// Directions to push the rock
	jets []bool
	// Current index into jets
	jetIdx int
	// convenience variable to track heighest point (heightOffset+len(occupied m))
	maxHeight int
	// Height from base (0) to the first element in m
	heightOffset int
}

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

func buildChamber(input string) (*chamberMap, error) {
	jets, err := buildJetDirections(input)
	if err != nil {
		return nil, err
	}

	chamber := chamberMap{
		m:    make([]chamberLevel, 1, 2022*2),
		jets: jets,
	}
	chamber.m[0] = chamberLevel{'|', '-', '-', '-', '-', '-', '-', '-', '|'}
	chamber.addLevel()
	chamber.addLevel()
	chamber.addLevel()

	return &chamber, nil
}

func (c *chamberMap) addLevel() {
	c.m = append(c.m, chamberLevel{'|', '.', '.', '.', '.', '.', '.', '.', '|'})
}

func (chamber *chamberMap) dropRock(rock *fallingRock) {
	// Init rock position
	rock.position.x = 3
	rock.position.y = chamber.maxHeight + 4

	// Extend chamber if needed
	for len(chamber.m) < rock.position.y+4-chamber.heightOffset {
		chamber.addLevel()
	}

	// Drop it until it hits
	for {
		// Push
		if chamber.jets[chamber.jetIdx] {
			// push right
			if chamber.isRightFree(*rock) {
				rock.position.x++
			}
		} else {
			// push left
			if chamber.isLeftFree(*rock) {
				rock.position.x--
			}
		}
		chamber.jetIdx++
		chamber.jetIdx %= len(chamber.jets)

		// Drop (break if can't)
		if chamber.rockCanFall(*rock) {
			rock.position.y--
		} else {
			break
		}
	}

	// Update chamber and currentHeight
	for yOffset, row := range rock.shape {
		for xOffset, pos := range row {
			if pos {
				chamber.m[rock.position.y+yOffset-chamber.heightOffset][rock.position.x+xOffset] = '#'
				if rock.position.y+yOffset > chamber.maxHeight {
					chamber.maxHeight = rock.position.y + yOffset
				}
			}
		}
	}

	// Check for tetris in applicable lines
	for yOffset := 3; yOffset >= 0; yOffset-- {
		blocked := true
		for x := 1; x < 8; x++ {
			if chamber.m[rock.position.y+yOffset-chamber.heightOffset][x] == '.' {
				blocked = false
				break
			}
		}
		if blocked {
			// Found a full blocked row. Delete everything below and update offsets
			fmt.Println("Found tetris at line ", rock.position.y+yOffset)
			chamber.m = chamber.m[rock.position.y+yOffset-chamber.heightOffset:]
			chamber.heightOffset = rock.position.y + yOffset
			break
		}
	}
}

func (chamber chamberMap) rockCanFall(rock fallingRock) bool {
	for xOffset := 3; xOffset >= 0; xOffset-- {
		// Find the lowest rock point in each column
		for yOffset := 0; yOffset < 4; yOffset++ {
			if rock.shape[yOffset][xOffset] {
				if chamber.m[rock.position.y+yOffset-1-chamber.heightOffset][rock.position.x+xOffset] == '.' {
					continue
				} else {
					return false
				}
			}
		}
	}
	return true
}

func (chamber chamberMap) isRightFree(rock fallingRock) bool {
	// func isRightFree(rock fallingRock, chamber chamberMap) bool {
	for yOffset := 0; yOffset < 4; yOffset++ {
		for xOffset := 3; xOffset >= 0; xOffset-- {
			// Find farthest right rock point in row and then
			if rock.shape[yOffset][xOffset] {
				if chamber.m[rock.position.y+yOffset-chamber.heightOffset][rock.position.x+xOffset+1] == '.' {
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

func (chamber chamberMap) isLeftFree(rock fallingRock) bool {
	// func isLeftFree(rock fallingRock, chamber chamberMap) bool {
	for yOffset := 0; yOffset < 4; yOffset++ {
		for xOffset := 0; xOffset < 4; xOffset++ {
			// Find farthest left rock point in row and then
			if rock.shape[yOffset][xOffset] {
				if chamber.m[rock.position.y+yOffset-chamber.heightOffset][rock.position.x+xOffset-1] == '.' {
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

func visualizeChamber(chamber *chamberMap) fyne.CanvasObject {
	label := widget.NewLabel(chamber.String())
	label.TextStyle.Monospace = true
	return container.NewHScroll(label)
}

func (c chamberMap) String() string {
	var sb strings.Builder
	for i := len(c.m) - 1; i >= 0; i-- {
		sb.WriteString(strconv.Itoa(i))
		for _, r := range c.m[i] {
			sb.WriteRune(r)
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}
