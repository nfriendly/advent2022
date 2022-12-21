package days

import (
	"errors"
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/gammazero/deque"
)

type Day5Solver struct {
}

// TODO: determine this programatically?
const maxElevation = 10

// TODO: fix bug when animate gets called multiple times (ex test and then solve)
// TODO: smarter approach = don't animate, instead have forward/backward buttons
// and each click executes an instruction. Don't animate, just display

func (d Day5Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	config, err := parseCargoCraneConfiguration(puzzleInput)

	animations := make([]*fyne.Animation, 0, len(config.instructions)*5)

	container := container.NewWithoutLayout()

	for x, stack := range config.stacks {
		for y := 0; y < stack.Len(); y++ {
			container.Add(stack.At(y))
			stack.At(y).Move(point{x: x, y: y}.toFynePos())
		}
	}
	currentAnimationIdx := 0

	fmt.Println("--(Start)------------")
	config.Visualize()

	failed := false
	for i, instruction := range config.instructions {
		for m := 0; m < instruction.numCrates; m++ {
			start := point{x: instruction.source - 1, y: config.stacks[instruction.source-1].Len()}.toFynePos()
			top0 := point{x: instruction.source - 1, y: maxElevation}.toFynePos()
			top1 := point{x: instruction.destination - 1, y: maxElevation}.toFynePos()
			end := point{x: instruction.destination - 1, y: config.stacks[instruction.destination-1].Len() + 1}.toFynePos()
			if config.stacks[instruction.source-1].Len() == 0 {
				fmt.Println("Failed to process instruction #", i, ", move number ", m)
				failed = true
				break
			}
			mover := config.stacks[instruction.source-1].PopBack()

			animations = append(animations, canvas.NewPositionAnimation(start, top0, time.Second*1, func(p fyne.Position) {
				mover.Move(p)
				canvas.Refresh(mover)
				if p == top0 {
					currentAnimationIdx++
					if currentAnimationIdx < len(animations) {
						animations[currentAnimationIdx].Start()
					}

				}
			}))
			animations = append(animations, canvas.NewPositionAnimation(top0, top1, time.Second*1, func(p fyne.Position) {
				mover.Move(p)
				canvas.Refresh(mover)
				if p == top1 {
					currentAnimationIdx++
					if currentAnimationIdx < len(animations) {
						animations[currentAnimationIdx].Start()
					}
				}
			}))
			animations = append(animations, canvas.NewPositionAnimation(top1, end, time.Second*1, func(p fyne.Position) {
				mover.Move(p)
				canvas.Refresh(mover)
				if p == end {
					currentAnimationIdx++
					if currentAnimationIdx < len(animations) {
						animations[currentAnimationIdx].Start()
					}
				}
			}))

			config.stacks[instruction.destination-1].PushBack(mover)
		}

		if failed {
			break
		}
	}

	animations[0].Start()
	fmt.Println("--(End)--------------")
	config.Visualize()
	if failed {
		return "", container, err
	}

	return buildPuzzleAnswer(config.stacks), container, err
}

func (d Day5Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	config, err := parseCargoCraneConfiguration(puzzleInput)

	fmt.Println("--(Start)------------")
	config.Visualize()

	failed := false
	for i, instruction := range config.instructions {
		var temp deque.Deque[*canvas.Text]
		if config.stacks[instruction.source-1].Len() < instruction.numCrates {
			fmt.Println("Failed to process instruction #", i)
			failed = true
			break
		}
		for i = 0; i < instruction.numCrates; i++ {
			mover := config.stacks[instruction.source-1].PopBack()
			temp.PushBack(mover)
		}
		for i = 0; i < instruction.numCrates; i++ {
			mover := temp.PopBack()
			config.stacks[instruction.destination-1].PushBack(mover)

		}

		if failed {
			break
		}
	}

	fmt.Println("--(End)--------------")
	config.Visualize()
	if failed {
		return "", nil, err
	}

	return buildPuzzleAnswer(config.stacks), nil, err

}

type point struct {
	x, y int
}

func (p point) toFynePos() fyne.Position {
	return fyne.Position{X: float32(20 + 20*p.x), Y: float32(20 + 20*p.y)}
}

func buildPuzzleAnswer(stacks []deque.Deque[*canvas.Text]) string {
	var sb strings.Builder

	for _, s := range stacks {
		sb.WriteString(s.Back().Text)
	}
	return sb.String()
}

type craneInstruction struct {
	numCrates, source, destination int
}

type craneConfiguration struct {
	stacks       []deque.Deque[*canvas.Text]
	instructions []craneInstruction
}

func parseCargoCraneConfiguration(input string) (craneConfiguration, error) {
	var config craneConfiguration
	lines := strings.Split(input, "\n")
	// Find split between initial conditions and directions
	split := 0
	for i, l := range lines {
		if l == "" {
			split = i
			break
		}
	}
	if split == 0 {
		return config, errors.New("couldn't find expected empty line")
	}

	// determine number of stacks
	stackNumString := strings.Split(strings.TrimSpace(lines[split-1]), " ")
	numStacks, err := strconv.Atoi(strings.TrimSpace(stackNumString[len(stackNumString)-1]))
	if err != nil {
		return config, err
	}
	config.stacks = make([]deque.Deque[*canvas.Text], numStacks)

	// Fill stacks
	for i := numStacks - 1; i >= 0; i-- {
		for j := 1; j < len(lines[i]); j += 4 {
			if lines[i][j] != ' ' {
				config.stacks[j/4].PushBack(canvas.NewText(string(lines[i][j]), color.Black))
			}
		}
	}

	// Build directions
	config.instructions = make([]craneInstruction, len(lines)-split-1)
	re := regexp.MustCompile(`move ([0-9]+) from ([0-9]+) to ([0-9]+)`)
	for i := split + 1; i < len(lines); i++ {
		res := re.FindSubmatch([]byte(lines[i]))
		if len(res) != 4 {
			return config, errors.New("Failed to parse line:" + lines[i])
		}
		var err error
		config.instructions[i-split-1].numCrates, err = strconv.Atoi(string(res[1]))
		if err != nil {
			return config, err
		}
		config.instructions[i-split-1].source, err = strconv.Atoi(string(res[2]))
		if err != nil {
			return config, err
		}
		config.instructions[i-split-1].destination, err = strconv.Atoi(string(res[3]))
		if err != nil {
			return config, err
		}
	}

	return config, nil
}

func (c craneConfiguration) Visualize() {
	for y := 0; y < 200; y++ {
		var sb strings.Builder
		for _, s := range c.stacks {
			if y < s.Len() {
				sb.WriteString(s.At(y).Text)
			} else {
				sb.WriteString(" ")
			}
		}
		if len(strings.TrimSpace(sb.String())) == 0 {
			return
		}
		fmt.Println(sb.String())
	}

}
