package days

import (
	"errors"
	"strconv"

	"fyne.io/fyne/v2"
)

type Day17Solver struct {
}

func (d Day17Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	jets, err := buildJetDirections(puzzleInput)
	if err != nil {
		return "", nil, err
	}

	chamber := make([]chamberLevel, 1, 2022*2)
	chamber[0] = chamberLevel{'-', '-', '-', '-', '-', '-', '-'}
	chamber = append(chamber, chamberLevel{'.', '.', '.', '.', '.', '.', '.'})

	// Drop 2022 rocks
	currentHeight := 0
	jetIdx := 0
	for rockIdx := 0; rockIdx < 2022; rockIdx++ {
		// Create rock
		rock := buildFallingRock(rockIdx)
		rock.position.x = 3
		rock.position.y = currentHeight + 4

		// Drop it until it hits
		for {
			// Push
			if jets[jetIdx] {
				// push right
				if rock.position.x+rock.width <= 7 {
					rock.position.x++
				}
			} else {
				// push left
				if rock.position.x > 1 {
					rock.position.x--
				}
			}

			// Drop (break if can't)
			if rockCanFall(rock, chamber[rock.position.y-1]) {
				rock.position.y--
			} else {
				break
			}
		}

		// Update chamber currentHeight
	}

	// Find height

	// Make image

	return "", nil, nil
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

type chamberLevel [7]rune

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
			{false, false, false, false},
			{false, true, false, false},
			{true, true, true, false},
			{false, true, false, false}},
			width: 3}
	case 2:
		return fallingRock{shape: [4][4]bool{
			{false, false, false, false},
			{false, false, true, false},
			{false, false, true, false},
			{true, true, true, false}},
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
			{false, false, false, false},
			{false, false, false, false},
			{true, true, false, false},
			{true, true, false, false}},
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
		} else if dir != '>' {
			return nil, errors.New("got unexpected rune at position " + strconv.Itoa(i) + ": " + string(dir))
		}
	}
	return res, nil
}

func rockCanFall(rock fallingRock, levelBelow chamberLevel) bool {
	for i := 0; i < 4; i++ {
		if rock.shape[3][i] && levelBelow[rock.position.x+i] != '.' {
			return false
		}
	}
	return true
}
