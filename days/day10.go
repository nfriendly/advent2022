package days

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Day10Solver struct {
}

func (d Day10Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	inst, err := buildSimpleCpuInstructions(puzzleInput)
	if err != nil {
		return "", nil, err
	}
	cpu := SimpleCpu{Instructions: inst, X: 1, Cycle: 0}
	signalStrengths := make([]int, 0, 6)
	for !cpu.Finsihed() {
		v, err := cpu.Tick()
		if err != nil {
			return "", nil, err
		}
		if cpu.Cycle == 20 || ((cpu.Cycle-20)%40 == 0) {
			signalStrengths = append(signalStrengths, cpu.Cycle*v)
		}
	}

	sumStrengths := 0
	for _, ss := range signalStrengths {
		sumStrengths += ss
	}

	return strconv.Itoa(sumStrengths), nil, nil
}

func (d Day10Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	inst, err := buildSimpleCpuInstructions(puzzleInput)
	if err != nil {
		return "", nil, err
	}
	cpu := SimpleCpu{Instructions: inst, X: 1, Cycle: 0}
	var display [240]bool
	for !cpu.Finsihed() {
		v, err := cpu.Tick()
		if err != nil {
			return "", nil, err
		}
		drawPos := (cpu.Cycle - 1) % 240
		display[drawPos] = (math.Abs(float64(v-(drawPos%40))) <= 1)
	}

	img, err := createDisplay(&display)
	return "(See image!)", img, err

}

func createDisplay(input *[240]bool) (fyne.CanvasObject, error) {
	var sb strings.Builder
	for i, b := range input {
		if i%40 == 0 {
			sb.WriteString("\n")
		}
		if b {
			sb.WriteString("#")
		} else {
			sb.WriteString(".")
		}
	}

	label := widget.NewLabel(sb.String())
	label.TextStyle.Monospace = true

	return label, nil
}

type SimpleCpuInstruction struct {
	Command  string
	Quantity int
}

type SimpleCpu struct {
	Instructions    []SimpleCpuInstruction
	X               int
	Cycle           int
	instructionIdx  int
	partialProgress int
}

// Returns the value of X during the Tick
func (c *SimpleCpu) Tick() (int, error) {
	value := c.X
	inst := c.Instructions[c.instructionIdx]
	switch inst.Command {
	case "noop":
		c.instructionIdx++
		c.partialProgress = 0
	case "addx":
		if c.partialProgress == 0 {
			// First cycle
			c.partialProgress++
		} else {
			// Second cycle
			c.X += inst.Quantity
			c.instructionIdx++
			c.partialProgress = 0
		}
	default:
		return value, errors.New("failed to process unknown command: " + inst.Command)
	}

	c.Cycle++
	return value, nil
}

func (c SimpleCpu) Finsihed() bool {
	return c.instructionIdx >= len(c.Instructions)
}

func buildSimpleCpuInstructions(input string) ([]SimpleCpuInstruction, error) {
	lines := strings.Split(input, "\n")
	insts := make([]SimpleCpuInstruction, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, " ")

		switch len(parts) {
		case 1:
			insts[i].Command = parts[0]
		case 2:
			var err error
			insts[i].Command = parts[0]
			insts[i].Quantity, err = strconv.Atoi(parts[1])
			if err != nil {
				return insts, err
			}
		default:
			return insts, errors.New("failed to parse instruction with too many parts: " + line)
		}
	}
	return insts, nil
}
