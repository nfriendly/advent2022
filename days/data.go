package days

import (
	"bytes"
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

type SinglePartTest struct {
	Input          string
	ExpectedOutput string
}

type DaySolver interface {
	SolvePartA(string) (string, fyne.CanvasObject, error)
	SolvePartB(string) (string, fyne.CanvasObject, error)
}

type Day struct {
	Number      int
	PartATests  []SinglePartTest
	PartBTests  []SinglePartTest
	PuzzleInput string
	PartAPrompt string
	PartBPrompt string
	Solver      DaySolver
}

func plotToImage(sol_plt *plot.Plot, imageName string) (*canvas.Image, error) {
	w, writer_err := sol_plt.WriterTo(5*vg.Inch, 3*vg.Inch, "png")
	if writer_err != nil {
		failMsg := fmt.Sprint("Failed to save plot, err: ", writer_err)
		return nil, errors.New(failMsg)
	}
	b := bytes.Buffer{}
	if _, err := w.WriteTo(&b); err != nil {
		failMsg := fmt.Sprint("Failed to save plot, err: ", err)
		return nil, errors.New(failMsg)
	}
	img := canvas.NewImageFromReader(&b, imageName)
	img.FillMode = canvas.ImageFillOriginal
	return img, nil
}

var Day1 = Day{
	Number:      1,
	PartATests:  day1TestsPartA,
	PartBTests:  day1TestsPartB,
	PuzzleInput: day1PuzzleInput,
	PartAPrompt: "Find the Elf carrying the most Calories. How many total Calories is that Elf carrying?",
	PartBPrompt: "Find the top three Elves carrying the most Calories. How many Calories are those Elves carrying in total?",
	Solver:      Day1Solver{},
}

var Day2 = Day{
	Number:      2,
	PartATests:  day2TestsPartA,
	PartBTests:  day2TestsPartB,
	PuzzleInput: day2PuzzleInput,
	PartAPrompt: "What would your total score be if everything goes exactly according to your strategy guide?",
	PartBPrompt: "Following the Elf's instructions for the second column, what would your total score be if everything goes exactly according to your strategy guide?",
	Solver:      Day2Solver{},
}
var Day3 = Day{
	Number:      3,
	PartATests:  day3TestsPartA,
	PartBTests:  day3TestsPartB,
	PuzzleInput: day3PuzzleInput,
	PartAPrompt: "Find the item type that appears in both compartments of each rucksack. What is the sum of the priorities of those item types?",
	PartBPrompt: "Find the item type that corresponds to the badges of each three-Elf group. What is the sum of the priorities of those item types?",
	Solver:      Day3Solver{},
}
var Day4 = Day{
	Number:      4,
	PartATests:  day4TestsPartA,
	PartBTests:  day4TestsPartB,
	PuzzleInput: day4PuzzleInput,
	PartAPrompt: "In how many assignment pairs does one range fully contain the other?",
	PartBPrompt: "In how many assignment pairs do the ranges overlap?",
	Solver:      Day4Solver{},
}
var Day5 = Day{
	Number:      5,
	PartATests:  day5TestsPartA,
	PartBTests:  day5TestsPartB,
	PuzzleInput: day5PuzzleInput,
	PartAPrompt: "After the rearrangement procedure completes, what crate ends up on top of each stack?",
	PartBPrompt: "After the rearrangement procedure completes, what crate ends up on top of each stack?",
	Solver:      Day5Solver{},
}
var Day6 = Day{
	Number:      6,
	PartATests:  day6TestsPartA,
	PartBTests:  day6TestsPartB,
	PuzzleInput: day6PuzzleInput,
	PartAPrompt: "How many characters need to be processed before the first start-of-packet marker is detected?",
	PartBPrompt: "How many characters need to be processed before the first start-of-message marker is detected?",
	Solver:      Day6Solver{},
}
var Day7 = Day{
	Number:      7,
	PartATests:  day7TestsPartA,
	PartBTests:  day7TestsPartB,
	PuzzleInput: day7PuzzleInput,
	PartAPrompt: "Find all of the directories with a total size of at most 100000. What is the sum of the total sizes of those directories?",
	PartBPrompt: "Find the smallest directory that, if deleted, would free up enough space. What is the total size of that directory?",
	Solver:      Day7Solver{},
}
var Day8 = Day{
	Number:      8,
	PartATests:  day8TestsPartA,
	PartBTests:  day8TestsPartB,
	PuzzleInput: day8PuzzleInput,
	PartAPrompt: "How many trees are visible from outside the grid?",
	PartBPrompt: "What is the highest scenic score possible for any tree?",
	Solver:      Day8Solver{},
}
var Day9 = Day{
	Number:      9,
	PartATests:  day9TestsPartA,
	PartBTests:  day9TestsPartB,
	PuzzleInput: day9PuzzleInput,
	PartAPrompt: "Simulate your complete hypothetical series of motions. How many positions does the tail of the rope visit at least once?",
	PartBPrompt: "Simulate your complete series of motions on a larger rope with ten knots. How many positions does the tail of the rope visit at least once?",
	Solver:      Day9Solver{},
}
var Day10 = Day{
	Number:      10,
	PartATests:  day10TestsPartA,
	PartBTests:  day10TestsPartB,
	PuzzleInput: day10PuzzleInput,
	PartAPrompt: "Find the signal strength during the 20th, 60th, 100th, 140th, 180th, and 220th cycles. What is the sum of these six signal strengths?",
	PartBPrompt: "Render the image given by your program. What eight capital letters appear on your CRT?",
	Solver:      Day10Solver{},
}
var Day11 = Day{
	Number:      11,
	PartATests:  day11TestsPartA,
	PartBTests:  day11TestsPartB,
	PuzzleInput: day11PuzzleInput,
	PartAPrompt: "What is the level of monkey business after 20 rounds of stuff-slinging simian shenanigans?",
	PartBPrompt: "What is the level of monkey business after 10000 rounds?",
	Solver:      Day11Solver{},
}
var Day12 = Day{
	Number:      12,
	PartATests:  day12TestsPartA,
	PartBTests:  day12TestsPartB,
	PuzzleInput: day12PuzzleInput,
	PartAPrompt: "What is the fewest steps required to move from your current position to the location that should get the best signal?",
	PartBPrompt: "What is the fewest steps required to move starting from any square with elevation a to the location that should get the best signal?",
	Solver:      Day12Solver{},
}
var Day13 = Day{
	Number:      13,
	PartATests:  day13TestsPartA,
	PartBTests:  day13TestsPartB,
	PuzzleInput: day13PuzzleInput,
	PartAPrompt: "Determine which pairs of packets are already in the right order. What is the sum of the indices of those pairs?",
	PartBPrompt: "Organize all of the packets into the correct order. What is the decoder key for the distress signal?",
	Solver:      Day13Solver{},
}
var Day14 = Day{
	Number:      14,
	PartATests:  day14TestsPartA,
	PartBTests:  day14TestsPartB,
	PuzzleInput: day14PuzzleInput,
	PartAPrompt: "How many units of sand come to rest before sand starts flowing into the abyss below?",
	PartBPrompt: "Using your scan, simulate the falling sand until the source of the sand becomes blocked. How many units of sand come to rest?",
	Solver:      Day14Solver{},
}
var Day15 = Day{
	Number:      15,
	PartATests:  day15TestsPartA,
	PartBTests:  day15TestsPartB,
	PuzzleInput: day15PuzzleInput,
	PartAPrompt: "Consult the report from the sensors you just deployed. In the row where y=2000000, how many positions cannot contain a beacon?",
	PartBPrompt: "Find the only possible position for the distress beacon. What is its tuning frequency?",
	Solver:      Day15Solver{},
}
var Day16 = Day{
	Number:      16,
	PartATests:  day16TestsPartA,
	PartBTests:  day16TestsPartB,
	PuzzleInput: day16PuzzleInput,
	PartAPrompt: "Work out the steps to release the most pressure in 30 minutes. What is the most pressure you can release?",
	PartBPrompt: "With you and an elephant working together for 26 minutes, what is the most pressure you could release?",
	Solver:      Day16Solver{},
}
var Day17 = Day{
	Number:      17,
	PartATests:  day17TestsPartA,
	PartBTests:  day17TestsPartB,
	PuzzleInput: day17PuzzleInput,
	PartAPrompt: "How many units tall will the tower of rocks be after 2022 rocks have stopped falling?",
	PartBPrompt: "TODO",
	Solver:      Day17Solver{},
}

var Days = []Day{{}, Day1, Day2, Day3, Day4, Day5, Day6, Day7, Day8, Day9, Day10, Day11, Day12,
	Day13, Day14, Day15, Day16, Day17}
