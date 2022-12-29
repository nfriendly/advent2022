package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"example.com/advent2022/days"
)

func TestDay1(t *testing.T) {
	res, _, err := days.Days[1].Solver.SolvePartA("1000")
	if err != nil {
		t.Error("Returned an error! err: " + err.Error())
	}

	if res != "1000" {
		t.Error("Returned: " + res + ", expected 1000")
	}
}

func TestDay3(t *testing.T) {
	puzzleInput := `vJrwpWtwJgWrhcsFMMfFFhFp
jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL
PmmdzqPrVvPwwTWBwg`
	lines := strings.Split(strings.TrimSpace(puzzleInput), "\n")
	for i := range lines {
		sort.Slice([]rune(lines[i]), func(j, k int) bool { return lines[i][j] < lines[i][k] })
		fmt.Println(lines[i])
	}
}

func TestDay5(t *testing.T) {
	_, _, err := days.Days[5].Solver.SolvePartA(days.Days[5].PartATests[0].Input)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDay5Puzzle(t *testing.T) {
	_, _, err := days.Days[5].Solver.SolvePartA(days.Days[5].PuzzleInput)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDay7(t *testing.T) {
	_, _, err := days.Days[7].Solver.SolvePartA(days.Days[7].PartATests[0].Input)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDay10(t *testing.T) {
	_, _, err := days.Days[10].Solver.SolvePartB(days.Days[10].PartBTests[0].Input)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDay11(t *testing.T) {
	_, _, err := days.Days[11].Solver.SolvePartA(days.Days[11].PartATests[0].Input)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDay13(t *testing.T) {
	_, _, err := days.Days[13].Solver.SolvePartB(days.Days[13].PartBTests[0].Input)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDay14(t *testing.T) {
	_, _, err := days.Days[14].Solver.SolvePartA(days.Days[14].PartBTests[0].Input)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDay15(t *testing.T) {
	_, _, err := days.Days[15].Solver.SolvePartA(days.Days[15].PartATests[0].Input)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDay16(t *testing.T) {
	_, _, err := days.Days[16].Solver.SolvePartA(days.Days[16].PartATests[0].Input)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDay17(t *testing.T) {
	_, _, err := days.Days[17].Solver.SolvePartA(days.Days[17].PartATests[0].Input)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDay19(t *testing.T) {
	_, _, err := days.Days[19].Solver.SolvePartA(days.Days[19].PartATests[0].Input)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDay20(t *testing.T) {
	// _, _, err := days.Days[20].Solver.SolvePartA(days.Days[20].PartATests[0].Input)
	_, _, err := days.Days[20].Solver.SolvePartA(days.Days[20].PuzzleInput)
	if err != nil {
		t.Error(err.Error())
	}
}
