package days

import (
	"errors"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
)

type Day21Solver struct {
}

func (d Day21Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	monkeyYellMap, err := buildMonkeyYellMap(puzzleInput)
	if err != nil {
		return "", nil, err
	}
	res, err := monkeyYellMap.solveFor("root")
	return strconv.Itoa(res), nil, err
}

func (d Day21Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	monkeyYellMap, err := buildMonkeyYellMap(puzzleInput)
	if err != nil {
		return "", nil, err
	}
	res, err := monkeyYellMap.findEquality()
	return strconv.Itoa(res), nil, err
}

type monkeyYell interface {
	isMonkeyYell()
}

type monkeyYellValue int
type monkeyYellOperation struct {
	left     string
	right    string
	operator string
}

func (monkeyYellValue) isMonkeyYell()     {}
func (monkeyYellOperation) isMonkeyYell() {}

type monkeyYellMap map[string]monkeyYell

func buildMonkeyYellMap(input string) (monkeyYellMap, error) {
	lines := strings.Split(input, "\n")
	m := make(monkeyYellMap, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			return m, errors.New("failed to parse parts from line: " + line)
		}
		val, err := strconv.Atoi(parts[1])
		if err == nil {
			m[parts[0]] = monkeyYellValue(val)
		} else {
			pieces := strings.Split(parts[1], " ")
			if len(pieces) != 3 {
				return m, errors.New("failed to parse pieces from line: " + line)
			}
			m[parts[0]] = monkeyYellOperation{
				left:     pieces[0],
				operator: pieces[1],
				right:    pieces[2],
			}
		}
	}
	return m, nil
}

func (m *monkeyYellMap) solveFor(target string) (int, error) {
	yell, present := (*m)[target]
	if !present {
		return 0, errors.New("failed to find key: " + target)
	}

	switch v := yell.(type) {
	case monkeyYellValue:
		return int(v), nil
	case monkeyYellOperation:
		result := 0
		left, err := m.solveFor(v.left)
		if err != nil {
			return 0, err
		}
		right, err := m.solveFor(v.right)
		if err != nil {
			return 0, err
		}
		switch v.operator {
		case "+":
			result = left + right
		case "-":
			result = left - right
		case "*":
			result = left * right
		case "/":
			result = left / right
		default:
			return 0, errors.New("Invalid operator: " + v.operator)
		}
		(*m)[target] = monkeyYellValue(result)
		return result, nil
	default:
		return 0, errors.New("unhandled type")
	}
}

func (m *monkeyYellMap) findEquality() (int, error) {
	// Figure out which side contains human
	topLevel, present := (*m)["root"]
	if !present {
		return 0, errors.New("map doesn't contain root node")
	}
	root, ok := topLevel.(monkeyYellOperation)
	if !ok {
		return 0, errors.New("root was a value")
	}
	leftReady := !m.containsHuman(root.left)
	rightReady := !m.containsHuman(root.right)
	if leftReady && rightReady {
		return 0, errors.New("illegal formation, both sides contain humman")
	}
	if !leftReady && !rightReady {
		return 0, errors.New("illegal formation, neither side contains humman")
	}
	// Calculate the value of the side withouth human
	var solvedSide int
	var err error
	var humanSide string
	if leftReady {
		solvedSide, err = m.solveFor(root.left)
		humanSide = root.right
	} else {
		solvedSide, err = m.solveFor(root.right)
		humanSide = root.left
	}
	if err != nil {
		return 0, err
	}

	// solve resulting equation by inverting non-human side of each branch onto solvide side
	return m.invertNonHumanHalf(humanSide, solvedSide)
}

func (m *monkeyYellMap) invertNonHumanHalf(target string, solvedSide int) (int, error) {
	if target == "humn" {
		return solvedSide, nil
	}

	topLevel, present := (*m)[target]
	if !present {
		return 0, errors.New("can't find target " + target + " in map for invertNonHumanHalf()")
	}
	root, ok := topLevel.(monkeyYellOperation)
	if !ok {
		return 0, errors.New("can't invert non-operation mokey yell")
	}
	leftReady := !m.containsHuman(root.left)
	rightReady := !m.containsHuman(root.right)

	if leftReady && rightReady {
		return 0, errors.New("illegal formation in invertNonHumanHalf(), both sides contain humman")
	}
	if !leftReady && !rightReady {
		return 0, errors.New("illegal formation in invertNonHumanHalf(), neither side contains humman")
	}
	var solvedHalf int
	var err error
	var humanSide string
	if leftReady {
		solvedHalf, err = m.solveFor(root.left)
		humanSide = root.right
	} else {
		solvedHalf, err = m.solveFor(root.right)
		humanSide = root.left
	}
	if err != nil {
		return 0, err
	}
	switch root.operator {
	case "+":
		// left + right = solvedSide
		// => right = solvedSide - left
		// => left = solvedSide - right
		solvedSide -= solvedHalf
	case "-":
		// left - right = solvedSide
		// => right = -(solvedSide - left)
		// => left = solvedSide + right
		if leftReady {
			solvedSide = -(solvedSide - solvedHalf)
		} else {
			solvedSide += solvedHalf
		}
	case "*":
		// left * right = solvedSide
		// => right = solvedSide / left
		// => left = solvedSide / right
		if solvedSide%solvedHalf != 0 {
			panic("unable to evenly divide (" +
				strconv.Itoa(solvedSide) + "/" +
				strconv.Itoa(solvedHalf) + ") while inverting " +
				target)
		}
		solvedSide /= solvedHalf
	case "/":
		// left / right = solvedSide
		// => right = left / solved side
		// => left = solvedSide * right
		if leftReady {
			solvedSide = solvedHalf / solvedSide
		} else {
			solvedSide *= solvedHalf
		}
	default:
		panic("Invalid operator: " + root.operator + " for invertNonHumanHalf()")
	}

	return m.invertNonHumanHalf(humanSide, solvedSide)
}

func (m *monkeyYellMap) containsHuman(target string) bool {
	if target == "humn" {
		return true
	}
	topLevel, present := (*m)[target]
	if !present {
		panic("can't find target " + target + " in map for containsHuman()")
	}
	root, ok := topLevel.(monkeyYellOperation)
	if !ok {
		return false
	}
	return m.containsHuman(root.left) || m.containsHuman(root.right)
}
