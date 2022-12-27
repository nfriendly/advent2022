package days

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Day19Solver struct {
}

func (d Day19Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	lines := strings.Split(puzzleInput, "\n")
	maxGeodes := make(plotter.Values, 0, len(lines))

	for _, line := range lines {
		bp, err := parseBlueprint(line)
		bp.maxMinutes = 24
		if err != nil {
			return "", nil, err
		}
		maxGeode := calculateMaxGeodes(bp)
		maxGeodes = append(maxGeodes, float64(maxGeode))

	}
	// todo calculate quality level sum
	totalQuality := 0
	for i, g := range maxGeodes {
		totalQuality += (i + 1) * int(g)
	}
	totalQualityLevel := strconv.Itoa(totalQuality)

	plt := plot.New()
	plt.Title.Text = "Maximum Geodes Produced"
	plt.X.Label.Text = "Blueprint Number"
	plt.Y.Label.Text = "Maxiumum Geodes"

	// TODO: adjust this value to make bar chart look nice
	w := vg.Points(10)
	bars, err := plotter.NewBarChart(maxGeodes, w)
	if err != nil {
		return totalQualityLevel, nil, err
	}

	plt.Add(bars)

	img, err := plotToImage(plt, "day19partA.png")

	return totalQualityLevel, img, err
}

func (d Day19Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	return "", nil, nil
}

type geodeBlueprint struct {
	oreRobotOre        int
	clayRobotOre       int
	obsidianRobotOre   int
	obsidianRobotClay  int
	geodeRobotOre      int
	geodeRobotObsidian int
	maxMinutes         int
}

type geodeState struct {
	ore, clay, obsidian, geodes                        int
	oreRobots, clayRobots, obsidianRobots, geodeRobots int
	minutesPassed                                      int
	actionHistory                                      []string
}

func parseBlueprint(input string) (geodeBlueprint, error) {
	bp := geodeBlueprint{}
	var err error
	re := regexp.MustCompile(
		`costs ([0-9]+) ore. [A-z ]+([0-9]+) ore. [A-z ]+([0-9]+) ore and ([0-9]+) clay. [A-z ]+ ([0-9]+) ore and ([0-9]+) obsidian`)

	parts := re.FindStringSubmatch(input)
	if len(parts) != 7 {
		return bp, errors.New("failed to parse line: " + input)
	}
	bp.oreRobotOre, err = strconv.Atoi(parts[1])
	if err != nil {
		return bp, err
	}
	bp.clayRobotOre, err = strconv.Atoi(parts[2])
	if err != nil {
		return bp, err
	}
	bp.obsidianRobotOre, err = strconv.Atoi(parts[3])
	if err != nil {
		return bp, err
	}
	bp.obsidianRobotClay, err = strconv.Atoi(parts[4])
	if err != nil {
		return bp, err
	}
	bp.geodeRobotOre, err = strconv.Atoi(parts[5])
	if err != nil {
		return bp, err
	}
	bp.geodeRobotObsidian, err = strconv.Atoi(parts[6])
	if err != nil {
		return bp, err
	}

	return bp, nil
}

func calculateMaxGeodes(bp geodeBlueprint) int {
	s := geodeState{oreRobots: 1}
	s.actionHistory = make([]string, 0, bp.maxMinutes)
	calculateMaxGeodesInner(&s, &bp)
	return s.geodes
}

func calculateMaxGeodesInner(state *geodeState, bp *geodeBlueprint) {
	//  TODO check for off by one
	for state.minutesPassed < bp.maxMinutes {
		// Try all 4 next builds and see what is best
		var stateCopies [4]geodeState
		for i := 0; i < len(stateCopies); i++ {
			stateCopies[i] = *state
		}
		if state.waitTimeNeededForOreRobot(bp) < bp.maxMinutes-state.minutesPassed {
			stateCopies[0].waitThenBuildOreRobot(bp)
			calculateMaxGeodesInner(&stateCopies[0], bp)
		} else {
			stateCopies[0].passTime(bp.maxMinutes-state.minutesPassed, "End wait")
		}

		if state.waitTimeNeededForClayRobot(bp) < bp.maxMinutes-state.minutesPassed {
			stateCopies[1].waitThenBuildClayRobot(bp)
			calculateMaxGeodesInner(&stateCopies[1], bp)
		} else {
			stateCopies[1].passTime(bp.maxMinutes-state.minutesPassed, "End wait")
		}
		if state.clayRobots > 0 && state.waitTimeNeededForObsidianRobot(bp) < bp.maxMinutes-state.minutesPassed {
			stateCopies[2].waitThenBuildObsidianRobot(bp)
			calculateMaxGeodesInner(&stateCopies[2], bp)
		} else {
			stateCopies[2].passTime(bp.maxMinutes-state.minutesPassed, "End wait")
		}
		if state.obsidianRobots > 0 && state.waitTimeNeededForGeodeRobot(bp) < bp.maxMinutes-state.minutesPassed {
			stateCopies[3].waitThenBuildGeodeRobot(bp)
			calculateMaxGeodesInner(&stateCopies[3], bp)
		} else {
			stateCopies[3].passTime(bp.maxMinutes-state.minutesPassed, "End wait")
		}

		// Chose best one
		maxIndex := 0
		maxResult := 0
		for i := range stateCopies {
			if stateCopies[i].geodes > maxResult {
				maxIndex = i
				maxResult = stateCopies[i].geodes
			}
		}
		*state = stateCopies[maxIndex]
	}
}

func (s *geodeState) passTime(minutes int, action string) {
	s.minutesPassed += minutes
	s.ore += s.oreRobots * minutes
	s.clay += s.clayRobots * minutes
	s.obsidian += s.obsidianRobots * minutes
	s.geodes += s.geodeRobots * minutes
	for i := 0; i < minutes; i++ {
		s.actionHistory = append(s.actionHistory, action)
	}
}

func (s *geodeState) waitTimeNeededForOreRobot(bp *geodeBlueprint) int {
	return int(math.Ceil(float64(bp.oreRobotOre-s.ore) / float64(s.oreRobots)))
}

func (s *geodeState) waitThenBuildOreRobot(bp *geodeBlueprint) {
	waitTime := s.waitTimeNeededForOreRobot(bp)
	if waitTime > 0 {
		s.passTime(waitTime, "wait")
	}
	s.buildOreRobot(bp)
}

func (s *geodeState) waitTimeNeededForClayRobot(bp *geodeBlueprint) int {
	return int(math.Ceil(float64(bp.clayRobotOre-s.ore) / float64(s.oreRobots)))
}

func (s *geodeState) waitThenBuildClayRobot(bp *geodeBlueprint) {
	waitTime := s.waitTimeNeededForClayRobot(bp)
	if waitTime > 0 {
		s.passTime(waitTime, "wait")
	}
	s.buildClayRobot(bp)
}

func (s *geodeState) waitTimeNeededForObsidianRobot(bp *geodeBlueprint) int {
	if s.clayRobots == 0 {
		panic("no clay robots so we'll never get enough clay")
	}
	oreWaitTime := math.Ceil(float64(bp.obsidianRobotOre-s.ore) / float64(s.oreRobots))
	clayWaitTime := math.Ceil(float64(bp.obsidianRobotClay-s.clay) / float64(s.clayRobots))
	return int(math.Max(oreWaitTime, clayWaitTime))
}

func (s *geodeState) waitThenBuildObsidianRobot(bp *geodeBlueprint) {
	waitTime := s.waitTimeNeededForObsidianRobot(bp)
	if waitTime > 0 {
		s.passTime(waitTime, "wait")
	}
	s.buildObsidianRobot(bp)
}

func (s *geodeState) waitTimeNeededForGeodeRobot(bp *geodeBlueprint) int {
	if s.obsidianRobots == 0 {
		panic("no obsidian robots so we'll never get enough obsidian")
	}
	oreWaitTime := math.Ceil(float64(bp.geodeRobotOre-s.ore) / float64(s.oreRobots))
	obsidianWaitTime := math.Ceil(float64(bp.geodeRobotObsidian-s.obsidian) / float64(s.obsidianRobots))
	return int(math.Max(oreWaitTime, obsidianWaitTime))
}

func (s *geodeState) waitThenBuildGeodeRobot(bp *geodeBlueprint) {
	waitTime := s.waitTimeNeededForGeodeRobot(bp)
	if waitTime > 0 {
		s.passTime(waitTime, "wait")
	}
	s.buildGeodeRobot(bp)
}

func (s *geodeState) buildOreRobot(bp *geodeBlueprint) {
	if s.ore < bp.oreRobotOre {
		panic("not enough ore to build ore robot")
	}
	s.ore -= bp.oreRobotOre
	s.passTime(1, "Built ore robot")
	s.oreRobots++
}

func (s *geodeState) buildClayRobot(bp *geodeBlueprint) {
	if s.ore < bp.clayRobotOre {
		panic("not enough ore to build clay robot")
	}
	s.ore -= bp.clayRobotOre
	s.passTime(1, "Built clay robot")
	s.clayRobots++
}

func (s *geodeState) buildObsidianRobot(bp *geodeBlueprint) {
	if s.ore < bp.obsidianRobotOre {
		panic("not enough ore to build obsidian robot")
	}
	if s.clay < bp.obsidianRobotClay {
		panic("not enough ore to build obsidian robot")
	}
	s.clay -= bp.obsidianRobotClay
	s.ore -= bp.obsidianRobotOre
	s.passTime(1, "Built obsidian robot")
	s.obsidianRobots++
}

func (s *geodeState) buildGeodeRobot(bp *geodeBlueprint) {
	if s.ore < bp.geodeRobotOre {
		panic("not enough ore to build geode robot")
	}
	if s.obsidian < bp.geodeRobotObsidian {
		panic("not enough ore to build geode robot")
	}
	s.obsidian -= bp.geodeRobotObsidian
	s.ore -= bp.geodeRobotOre
	s.passTime(1, "Built geode robot")
	s.geodeRobots++
}
