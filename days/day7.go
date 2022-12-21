package days

import (
	"errors"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Day7Solver struct {
}

func (d Day7Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	root, err := buildDirectoryTree(puzzleInput)
	if err != nil {
		return "", nil, err
	}
	fixDirectorySizes(root)
	totalSize := sumDirectorySizeIf(root, func(d *directoryTreeItem) bool {
		return d.isDirectory && d.cumulativeSize < 100000
	})
	// TODO: make tree view bigger, probably requires changing the layout in the main app :(
	ft := directoryTreeToFyneTree(root)
	return strconv.Itoa(totalSize), ft, nil
}

func (d Day7Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	root, err := buildDirectoryTree(puzzleInput)
	if err != nil {
		return "", nil, err
	}
	fixDirectorySizes(root)

	freeSpace := 70000000 - root.cumulativeSize
	requiredSize := 30000000 - freeSpace

	directoryName, directorySize := findSmallestDirectoryBiggerThan(root, requiredSize)

	return strconv.Itoa(directorySize), widget.NewLabel("Smallest directory that's big enough: " + directoryName), nil
}

type directoryTreeItem struct {
	name           string
	isDirectory    bool
	cumulativeSize int
	children       map[string]*directoryTreeItem
}

func (d *directoryTreeItem) fyneName() string {
	var sb strings.Builder
	sb.WriteString(d.name)
	if d.isDirectory {
		sb.WriteString(" (dir, size=")
	} else {
		sb.WriteString(" (file, size=")
	}
	sb.WriteString(strconv.Itoa(d.cumulativeSize))
	sb.WriteString(")")
	return sb.String()
}

// Builds a tree and returns the root item
func buildDirectoryTree(input string) (*directoryTreeItem, error) {
	// Parse commnad
	commands := strings.Split(strings.TrimSpace(input), "\n")
	cmds := make([][]string, 0, len(commands))
	for _, c := range commands {
		cmds = append(cmds, strings.Split(strings.TrimSpace(c), " "))
	}

	breadcrumbs := make([]*directoryTreeItem, 0, 20)

	// If cd NAME:
	//   append name to currentDirectory
	// If cd ..
	//   update directory size
	//   remove back of currentDirectory
	// If ls
	//   create directoryTreeItem for each listed and attach to children of current dir

	for i := 0; i < len(cmds); i++ {
		if cmds[i][0] != "$" {
			return breadcrumbs[0], errors.New("Got data line when command was expected, line = " + commands[i])
		}
		switch cmds[i][1] {
		case "cd":
			if cmds[i][2] == ".." {
				if len(breadcrumbs) <= 1 {
					return breadcrumbs[0], errors.New("attempted to change directory up past root")
				}
				breadcrumbs = breadcrumbs[:len(breadcrumbs)-1]
			} else {
				if len(breadcrumbs) == 0 {
					breadcrumbs = append(breadcrumbs, &directoryTreeItem{
						name:           cmds[i][2],
						isDirectory:    true,
						cumulativeSize: 0,
						children:       make(map[string]*directoryTreeItem),
					})
					break
				}
				child, ok := breadcrumbs[len(breadcrumbs)-1].children[cmds[i][2]]
				if !ok {
					var errorMsg strings.Builder
					errorMsg.WriteString("Couldn't find child ")
					errorMsg.WriteString(cmds[i][2])
					errorMsg.WriteString(" in expected children: ")
					for k := range breadcrumbs[len(breadcrumbs)-1].children {
						errorMsg.WriteString(k + ", ")
					}
					return breadcrumbs[0], errors.New(errorMsg.String())
				}

				breadcrumbs = append(breadcrumbs, child)
			}
		case "ls":
			for i < len(commands)-1 {
				i++
				if cmds[i][0] == "$" {
					// Exits inner for loop, not switch
					i--
					break
				} else if cmds[i][0] == "dir" {
					breadcrumbs[len(breadcrumbs)-1].children[cmds[i][1]] = &directoryTreeItem{
						name:           cmds[i][1],
						isDirectory:    true,
						cumulativeSize: 0,
						children:       make(map[string]*directoryTreeItem),
					}
				} else {
					size, err := strconv.Atoi(cmds[i][0])
					if err != nil {
						return breadcrumbs[0], errors.New("Failed to convert size to int in line " + commands[i])
					}
					breadcrumbs[len(breadcrumbs)-1].children[cmds[i][1]] = &directoryTreeItem{
						name:           cmds[i][1],
						isDirectory:    false,
						cumulativeSize: size,
					}
				}
			}

		default:
			return breadcrumbs[0], errors.New("Failed to parse command in line: " + commands[i])
		}
	}
	return breadcrumbs[0], nil
}

func fixDirectorySizes(currentNode *directoryTreeItem) {
	if !currentNode.isDirectory {
		return
	}
	total := 0
	for _, c := range currentNode.children {
		fixDirectorySizes(c)
		total += c.cumulativeSize
	}
	currentNode.cumulativeSize = total
}

func directoryTreeToFyneTree(dt *directoryTreeItem) *widget.Tree {
	ft := make(map[string][]string, 100)
	root := make([]string, 1)
	root[0] = dt.fyneName()
	ft[""] = root
	descendTree(dt, ft)
	return widget.NewTreeWithStrings(ft)
}

func descendTree(currentNode *directoryTreeItem, mapSoFar map[string][]string) {
	// mapSoFar[]
	if !currentNode.isDirectory {
		return
	}
	kids := make([]string, 0)
	for _, c := range currentNode.children {
		kids = append(kids, c.fyneName())
		descendTree(c, mapSoFar)
	}
	mapSoFar[currentNode.fyneName()] = kids
}

func sumDirectorySizeIf(currentNode *directoryTreeItem, shouldCount func(*directoryTreeItem) bool) int {
	totalSize := 0
	if shouldCount(currentNode) {
		totalSize += currentNode.cumulativeSize
	}
	if currentNode.isDirectory {
		for _, c := range currentNode.children {
			totalSize += sumDirectorySizeIf(c, shouldCount)
		}
	}

	return totalSize
}

func findSmallestDirectoryBiggerThan(root *directoryTreeItem, minSize int) (string, int) {
	validDirectories := make([]*directoryTreeItem, 0, 20)
	findSmallestDirectoryBiggerThanInner(root, minSize, &validDirectories)

	minDirName := validDirectories[0].name
	minDirSize := validDirectories[0].cumulativeSize
	for _, d := range validDirectories {
		if d.cumulativeSize < minDirSize {
			minDirName = d.name
			minDirSize = d.cumulativeSize
		}
	}
	return minDirName, minDirSize
}

func findSmallestDirectoryBiggerThanInner(currentNode *directoryTreeItem, minSize int, validDirectories *[]*directoryTreeItem) {
	if currentNode.isDirectory {
		if currentNode.cumulativeSize > minSize {
			*validDirectories = append(*validDirectories, currentNode)
		}
		for _, c := range currentNode.children {
			findSmallestDirectoryBiggerThanInner(c, minSize, validDirectories)
		}
	}
}
