package days

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gammazero/deque"
)

type Day16Solver struct {
}

func (d Day16Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	vd, err := buildValveData(puzzleInput)
	if err != nil {
		return "", nil, err
	}

	// start with lazy BFS where each node has the edges: move to each neighbor, open valve (if current is closed)
	score, path, err := valveDepthBFS(&vd)
	if err != nil {
		return "", nil, err
	}

	img, err := visualizeValvePath(path)

	return strconv.Itoa(score), img, err
}

func (d Day16Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	return "", nil, nil
}

type valveData struct {
	rates     map[string]int
	neighbors map[string][]string
}

func valveDepthBFS(vd *valveData) (int, []valveStateNode, error) {
	// Find the max rate as an easy way to stop searching early
	maxRate := 0
	for _, r := range vd.rates {
		if r > maxRate {
			maxRate = r
		}
	}
	frontier := deque.New[*valveStateNode](30 * len(vd.rates))
	path := make([]valveStateNode, 0)
	stateMap := make(map[string]bool)
	for key := range vd.rates {
		stateMap[key] = false
	}
	frontier.PushBack(&valveStateNode{
		timePassed: 0,
		score:      0,
		valveState: stateMap,
		location:   "AA",
		parent:     nil})

	visited := make(map[string][]visistedValveNode)

	maxScore := 0
	var maxNode *valveStateNode

	i := 0

	for frontier.Len() > 0 {
		current := frontier.PopFront()
		if containsOrBetter(current.score, makeOnValvesString(current.valveState), visited[current.location]) {
			continue
		}

		v := visited[current.location]
		v = append(v, visistedValveNode{
			score:    current.score,
			onValves: makeOnValvesString(current.valveState),
		})
		visited[current.location] = v

		// v = append(v, )
		// visited = append(visited, current)

		// todo if maxrate can't beat best score, continue

		if current.timePassed < 30 {
			// Add open visit valve if possible and worthwhile
			if !current.valveState[current.location] && vd.rates[current.location] > 0 {
				// newStateMap := current.valveState
				newStateMap := make(map[string]bool)
				for id, v := range current.valveState {
					newStateMap[id] = v
				}
				newStateMap[current.location] = true
				frontier.PushBack(&valveStateNode{
					timePassed: current.timePassed + 1,
					score:      current.score + vd.rates[current.location]*(30-current.timePassed-1),
					valveState: newStateMap,
					location:   current.location,
					parent:     current})
				if frontier.Back().score > maxScore {
					maxScore = frontier.Back().score
					maxNode = frontier.Back()
				}
			}
			// Add move to neighbors
			for _, neighbor := range vd.neighbors[current.location] {
				frontier.PushBack(&valveStateNode{
					timePassed: current.timePassed + 1,
					score:      current.score,
					valveState: current.valveState,
					location:   neighbor,
					parent:     current})
			}
		}
		i++

		if i%10000 == 0 {
			fmt.Println("Nodes explored: ", i, ", size of frontier: ", frontier.Len())
		}
	}

	pNode := maxNode
	for pNode != nil {
		path = append(path, *pNode)
		pNode = pNode.parent
	}

	return maxScore, path, nil
}

func visualizeValvePath(path []valveStateNode) (fyne.CanvasObject, error) {
	prevLoc := "AA"
	var sb strings.Builder
	sb.WriteString(prevLoc)
	for i := len(path) - 2; i >= 0; i-- {
		sb.WriteString("->")
		if path[i].location == prevLoc {
			sb.WriteString("[")
			sb.WriteString(path[i].location)
			sb.WriteString("]")
		} else {
			sb.WriteString(path[i].location)
		}
		prevLoc = path[i].location
	}

	label := widget.NewLabel(sb.String())
	label.TextStyle.Monospace = true
	return container.NewHScroll(label), nil
}

type valveStateNode struct {
	timePassed int
	score      int
	valveState map[string]bool
	location   string
	parent     *valveStateNode
}

type visistedValveNode struct {
	score    int
	onValves string
}

func containsOrBetter(score int, onValues string, nodes []visistedValveNode) bool {
	for _, n := range nodes {
		if onValues == n.onValves && score <= n.score {
			return true
		}
	}
	return false
}

func makeOnValvesString(valveState map[string]bool) string {
	onValves := make([]string, 0)
	for k, v := range valveState {
		if v {
			onValves = append(onValves, k)
		}
	}
	sort.Strings(onValves)
	return strings.Join(onValves, ",")
}

// type valvePriorityQueue []valveStateNode

// func (q valvePriorityQueue) Len() int { return len(q) }
// func (q valvePriorityQueue) Less(i, j int) bool {
// 	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
// 	return q[i].score-q[i].timePassed > q[j].score-q[j].timePassed
// }
// func (q valvePriorityQueue) Swap(i, j int)       { q[i], q[j] = q[j], q[i] }
// func (q *valvePriorityQueue) Push(n interface{}) { *q = append(*q, n.(valveStateNode)) }
// func (q *valvePriorityQueue) Pop() interface{} {
// 	t := *q
// 	var n interface{}
// 	n, *q = t[len(t)-1], t[:len(t)-1]
// 	return n
// }

func buildValveData(input string) (valveData, error) {
	vd := valveData{
		rates:     make(map[string]int, 0),
		neighbors: make(map[string][]string),
	}
	lines := strings.Split(input, "\n")

	re0 := regexp.MustCompile(`Valve ([A-Z][A-Z]) has flow rate=([0-9]+)`)
	re1 := regexp.MustCompile(`([A-Z][A-Z])`)
	for _, line := range lines {
		innerParts := strings.Split(line, ";")
		if len(innerParts) != 2 {
			return vd, errors.New("failed to parse ; in line: " + line)
		}
		firstParts := re0.FindStringSubmatch(innerParts[0])
		if len(firstParts) != 3 {
			return vd, errors.New("failed to parse first half of line: " + line)
		}
		name := firstParts[1]
		rate, err := strconv.Atoi(firstParts[2])
		if err != nil {
			return vd, nil
		}
		vd.rates[name] = rate
		vd.neighbors[name] = re1.FindAllString(innerParts[1], -1)
	}
	return vd, nil
}

func (vd valveData) String() string {
	var sb strings.Builder
	for name := range vd.rates {
		sb.WriteString("Valve ")
		sb.WriteString(name)
		sb.WriteString(" has flow rate=")
		sb.WriteString(strconv.Itoa(vd.rates[name]))
		sb.WriteString("; tunnels lead to valves ")
		sb.WriteString(strings.Join(vd.neighbors[name], ", "))
		sb.WriteString("\n")
	}

	return sb.String()
}
