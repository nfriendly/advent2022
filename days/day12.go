package days

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type Day12Solver struct {
}

func (d Day12Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	g := buildElevationGraph(puzzleInput)

	// Can seem to create a new node object with the same id as the start
	pth := path.DijkstraFrom(g.Node(g.startId), g)
	sol, _ := pth.To(g.endId)
	solStr := strconv.Itoa(len(sol) - 1)

	fmt.Println(sol)

	path := make(plotter.XYs, 0, len(sol))
	for i := range sol {
		if n, ok := sol[i].(elevationNode); ok {
			// negate Y so it visually looks like the prompts
			path = append(path, plotter.XY{X: float64(n.location.x), Y: float64(-n.location.y)})
		} else {
			return solStr, nil, errors.New("got a bad type back in path solution")
		}
	}
	l, err := plotter.NewLine(path)
	if err != nil {
		return solStr, nil, err
	}

	plt := plot.New()
	plt.Title.Text = "Solution Path"
	plt.X.Label.Text = "X"
	plt.Y.Label.Text = "Y"

	plt.Add(l)
	img, err := plotToImage(plt, "day12partA.png")
	return solStr, img, err
}

func (d Day12Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	g := buildElevationGraph(puzzleInput)

	// Find shortest path from end to every lowest elevation, save the shortest of these
	var shortestPath []graph.Node
	for _, lowId := range g.lowestElevationIds {
		pth := path.DijkstraFrom(g.Node(lowId), g)
		sol, _ := pth.To(g.endId)
		if len(sol) == 0 {
			// No solution found
			continue
		}
		if len(shortestPath) == 0 || len(sol) < len(shortestPath) {
			shortestPath = sol
		}
	}

	solStr := strconv.Itoa(len(shortestPath) - 1)

	fmt.Println(shortestPath)

	path := make(plotter.XYs, 0, len(shortestPath))
	for i := range shortestPath {
		if n, ok := shortestPath[i].(elevationNode); ok {
			// negate Y so it visually looks like the prompts
			path = append(path, plotter.XY{X: float64(n.location.x), Y: float64(-n.location.y)})
		} else {
			return solStr, nil, errors.New("got a bad type back in path solution")
		}
	}
	l, err := plotter.NewLine(path)
	if err != nil {
		return solStr, nil, err
	}

	plt := plot.New()
	plt.Title.Text = "Solution Path"
	plt.X.Label.Text = "X"
	plt.Y.Label.Text = "Y"

	plt.Add(l)
	img, err := plotToImage(plt, "day12partB.png")
	return solStr, img, err
}

type elevationGraph struct {
	startId            int64
	endId              int64
	lowestElevationIds []int64
	*simple.DirectedGraph
}

type elevationNode struct {
	elevation rune
	location  point
	id        int64
}

func (n elevationNode) ID() int64      { return n.id }
func (n elevationNode) String() string { return string(n.elevation) }

func buildElevationGraph(input string) elevationGraph {
	g := elevationGraph{lowestElevationIds: make([]int64, 0), DirectedGraph: simple.NewDirectedGraph()}
	lines := strings.Split(input, "\n")
	height := len(lines)
	width := len(lines[0])
	nodes := make([]elevationNode, height*width)
	for y, line := range lines {
		for x, e := range line {
			id := y*width + x
			nodes[id].id = int64(id)
			nodes[id].location = point{x: x, y: y}
			switch e {
			case 'S':
				g.startId = int64(id)
				g.lowestElevationIds = append(g.lowestElevationIds, int64(id))
				nodes[id].elevation = 'a'
			case 'E':
				g.endId = int64(id)
				nodes[id].elevation = 'z'
			case 'a':
				g.lowestElevationIds = append(g.lowestElevationIds, int64(id))
				nodes[id].elevation = e
			default:
				nodes[id].elevation = e
			}
			g.AddNode(nodes[id])
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			id := y*width + x
			if y > 0 {
				upId := (y-1)*width + x
				addEdge(&g, &nodes, id, upId)
			}
			if y < height-1 {
				downId := (y+1)*width + x
				addEdge(&g, &nodes, id, downId)
			}
			if x > 0 {
				leftId := y*width + x - 1
				addEdge(&g, &nodes, id, leftId)
			}
			if x < width-1 {
				rightId := y*width + x + 1
				addEdge(&g, &nodes, id, rightId)
			}
		}
	}

	return g
}

func addEdge(g *elevationGraph, nodes *[]elevationNode, fromId int, toId int) {
	from := g.Node(int64(fromId))
	to := g.Node(int64(toId))
	if from == nil || to == nil {
		panic("Couldn't find nodes to add")
	}
	if (*nodes)[toId].elevation <= (*nodes)[fromId].elevation+1 {
		g.SetEdge(simple.Edge{F: from, T: to})
	}
}
