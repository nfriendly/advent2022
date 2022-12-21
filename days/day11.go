package days

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/gammazero/deque"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

type Day11Solver struct {
}

func (d Day11Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	return calculateMonkeyBusiness(puzzleInput, true, 20)
}

func (d Day11Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	return calculateMonkeyBusiness(puzzleInput, false, 10000)
}

func calculateMonkeyBusiness(puzzleInput string, partA bool, numRounds int) (string, fyne.CanvasObject, error) {
	ms, err := buildMonkeys(puzzleInput)
	if err != nil {
		return "", nil, err
	}

	inspectionHistory := make([]plotter.XYs, len(ms))
	for i := 0; i < len(ms); i++ {
		inspectionHistory[i] = make(plotter.XYs, numRounds+1)
		inspectionHistory[i][0].X = 0
		inspectionHistory[i][0].Y = 0
	}

	modulator := 1
	if !partA {
		for i := 0; i < len(ms); i++ {
			modulator *= ms[i].testDivisibleBy
		}
	}

	for i := 0; i < numRounds; i++ {
		ms.performRound(partA, modulator)
		for m := 0; m < len(ms); m++ {
			inspectionHistory[m][i+1].X = float64(i + 1)
			inspectionHistory[m][i+1].Y = float64(ms[m].totalInspections)
		}
	}

	finalValues := make([]int, len(ms))
	for i, m := range ms {
		finalValues[i] = m.totalInspections
	}

	sort.Ints(finalValues)
	result := finalValues[len(finalValues)-1] * finalValues[len(finalValues)-2]

	plt := plot.New()
	plt.Title.Text = "Monkey Inspection Counts"
	plt.X.Label.Text = "Round"
	plt.Y.Label.Text = "Total Items Inspected"

	// Mostly a reimplementation of plotutil.AddLinePoints() because I don't know how to alternate string and data
	for i, hist := range inspectionHistory {
		if partA {
			line, s, err := plotter.NewLinePoints(hist)
			if err != nil {
				return strconv.Itoa(result), nil, err
			}
			line.Color = plotutil.Color(i)
			line.Dashes = plotutil.Dashes(i)
			s.Color = plotutil.Color(i)
			s.Shape = plotutil.Shape(i)
			plt.Add(line)
			plt.Add(s)
			plt.Legend.Add("Monkey "+strconv.Itoa(i), line, s)
		} else {
			line, err := plotter.NewLine(hist)
			if err != nil {
				return strconv.Itoa(result), nil, err
			}
			line.Color = plotutil.Color(i)
			line.Dashes = plotutil.Dashes(i)
			plt.Add(line)
			plt.Legend.Add("Monkey "+strconv.Itoa(i), line)
		}
	}

	var name string
	if partA {
		name = "day11partA.png"
	} else {
		name = "day11partB.png"
	}
	img, err := plotToImage(plt, name)

	return strconv.Itoa(result), img, err
}

type monkey struct {
	items            deque.Deque[int]
	operator         string
	leftOperand      string
	rightOperand     string
	testDivisibleBy  int
	trueTarget       int
	falseTarget      int
	totalInspections int
}

type monkeys []monkey

func (m *monkey) inspectItems(partA bool, modulator int) {
	for i := 0; i < m.items.Len(); i++ {
		m.totalInspections++

		old := m.items.At(i)
		var left, right int
		var err error
		if m.leftOperand == "old" {
			left = old
		} else {
			left, err = strconv.Atoi(m.leftOperand)
			if err != nil {
				panic(err)
			}
		}
		if m.rightOperand == "old" {
			right = old
		} else {
			right, err = strconv.Atoi(m.rightOperand)
			if err != nil {
				panic(err)
			}
		}

		new := 0
		switch m.operator {
		case "+":
			new = left + right
		case "-":
			new = left - right
		case "*":
			new = left * right
		default:
			panic(errors.New("Invalid operator: " + m.operator))
		}
		if partA {
			m.items.Set(i, new/3)
		} else {
			m.items.Set(i, new%modulator)
		}
	}
}

func (ms *monkeys) performRound(partA bool, modulator int) {
	for i := 0; i < len(*ms); i++ {
		m := &(*ms)[i]
		m.inspectItems(partA, modulator)
		for m.items.Len() > 0 {
			item := m.items.PopFront()
			if item%m.testDivisibleBy == 0 {
				(*ms)[m.trueTarget].items.PushBack(item)
			} else {
				(*ms)[m.falseTarget].items.PushBack(item)
			}
		}
	}
}

func buildMonkeys(input string) (monkeys, error) {
	lines := strings.Split(input, "\n")
	ms := make(monkeys, (len(lines)/7)+1)

	for i, line := range lines {
		parts := strings.Split(strings.TrimSpace(line), " ")
		switch i % 7 {
		case 0:
			if parts[0] != "Monkey" {
				return nil, errors.New("parser out of sync at line #" + strconv.Itoa(i) + ": " + line)
			}

		case 1:
			if parts[0] != "Starting" {
				return nil, errors.New("parser out of sync at line #" + strconv.Itoa(i) + ": " + line)
			}
			for idx := 2; idx < len(parts); idx++ {
				v, err := strconv.Atoi(strings.TrimSuffix(parts[idx], ","))
				if err != nil {
					return nil, errors.New("Failed to parse int at position " + strconv.Itoa(idx) + " in line: " + line)
				}
				ms[i/7].items.PushBack(v)
			}

		case 2:
			if parts[0] != "Operation:" || len(parts) != 6 {
				return nil, errors.New("parser out of sync at line #" + strconv.Itoa(i) + ": " + line)
			}

			_, err := strconv.Atoi(parts[3])
			if parts[3] != "old" && err != nil {
				return nil, errors.New("failed to parse left operand in line: " + line)
			}
			_, err = strconv.Atoi(parts[5])
			if parts[5] != "old" && err != nil {
				return nil, errors.New("failed to parse right operand in line: " + line + ". " + err.Error())
			}
			switch parts[4] {
			case "+":
			case "-":
			case "*":
			default:
				return nil, errors.New("failed to parse operator in line: " + line)
			}

			ms[i/7].leftOperand = parts[3]
			ms[i/7].operator = parts[4]
			ms[i/7].rightOperand = parts[5]

		case 3:
			if parts[0] != "Test:" || len(parts) != 4 {
				return nil, errors.New("parser out of sync at line #" + strconv.Itoa(i) + ": " + line)
			}
			v, err := strconv.Atoi(parts[3])
			if err != nil {
				return nil, errors.New("failed to parse divisibility int in line: " + line)
			}
			ms[i/7].testDivisibleBy = v

		case 4:
			if parts[0] != "If" || parts[1] != "true:" || len(parts) != 6 {
				return nil, errors.New("parser out of sync at line #" + strconv.Itoa(i) + ": " + line)
			}
			v, err := strconv.Atoi(parts[5])
			if err != nil {
				return nil, errors.New("failed to parse throw target in line: " + line)
			}
			ms[i/7].trueTarget = v

		case 5:
			if parts[0] != "If" || parts[1] != "false:" || len(parts) != 6 {
				return nil, errors.New("parser out of sync at line #" + strconv.Itoa(i) + ": " + line)
			}
			v, err := strconv.Atoi(parts[5])
			if err != nil {
				return nil, errors.New("failed to parse throw target in line: " + line)
			}
			ms[i/7].falseTarget = v

		case 6:
			if parts[0] != "" {
				return nil, errors.New("parser out of sync at line #" + strconv.Itoa(i) + ": " + line)
			}

		}
	}

	return ms, nil
}
