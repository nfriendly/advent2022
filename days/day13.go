package days

import (
	"errors"
	"image/color"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/gammazero/deque"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type Day13Solver struct {
}

func (d Day13Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	lines := strings.Split(puzzleInput, "\n")
	rightIndices := make([]int, 0, len(lines)/6)
	wrongIndices := make([]int, 0, len(lines)/6)
	// len+1 because there's no empty line at the end of the input
	for i := 0; i < (len(lines)+1)/3; i++ {
		packet, err := buildDistressPacket(lines[3*i], lines[3*i+1])
		if err != nil {
			return "", nil, err
		}
		res := isSorted(packet.first, packet.second)
		switch res {
		case 0:
			return "", nil, errors.New("got two identical packets: " + lines[3*i] + ", and: " + lines[3*i+1])
		case 1:
			rightIndices = append(rightIndices, i+1)
		case -1:
			wrongIndices = append(wrongIndices, i+1)
		default:
			return "", nil, errors.New("got unexpected return from isSorted: " + strconv.Itoa(res))
		}
	}

	rightSum := 0
	for _, idx := range rightIndices {
		rightSum += idx
	}

	solStr := strconv.Itoa(rightSum)

	rightSeries := make(plotter.XYs, len(rightIndices))
	wrongSeries := make(plotter.XYs, len(wrongIndices))
	for i, rdx := range rightIndices {
		rightSeries[i].X = float64(rdx)
		rightSeries[i].Y = 1.0
	}
	for i, rdx := range wrongIndices {
		wrongSeries[i].X = float64(rdx)
		wrongSeries[i].Y = 0.0
	}
	sr, err := plotter.NewScatter(rightSeries)
	sr.Color = color.RGBA{G: 255, A: 255}
	if err != nil {
		return solStr, nil, err
	}
	sw, err := plotter.NewScatter(wrongSeries)
	if err != nil {
		return solStr, nil, err
	}
	sw.Color = color.RGBA{R: 255, A: 255}

	plt := plot.New()
	plt.Title.Text = "Ordered Packet Pairs"
	plt.X.Label.Text = "Packet Pair Index"
	plt.Y.Label.Text = "Corret"
	plt.NominalY("Wrong Order", "Right Order")

	plt.Add(sr)
	plt.Add(sw)
	img, err := plotToImage(plt, "day12partA.png")
	return solStr, img, err
}

func (d Day13Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	lines := strings.Split(puzzleInput, "\n")

	packets := make(distressDataSlice, 0)
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			p, err := buildDistressData(l)
			if err != nil {
				return "", nil, err
			}
			packets = append(packets, p)
		}
	}
	// Add separators
	sep0, err := buildDistressData("[[2]]")
	if err != nil {
		return "", nil, err
	}
	sep1, err := buildDistressData("[[6]]")
	if err != nil {
		return "", nil, err
	}
	packets = append(packets, sep0, sep1)

	sort.Sort(packets)

	// Find separators
	s0idx := sort.Search(len(packets), func(i int) bool { return isSorted(sep0, packets[i]) >= 0 })
	s1idx := sort.Search(len(packets), func(i int) bool { return isSorted(sep1, packets[i]) >= 0 })

	solStr := strconv.Itoa((s0idx + 1) * (s1idx + 1))

	label := widget.NewLabel("The separators get sorted into index " + strconv.Itoa((s0idx + 1)) + " and " + strconv.Itoa((s1idx + 1)))

	return solStr, label, nil
}

type distressData interface {
	isDistressData()
}

type distressDataValue int
type distressDataList []distressData

func (distressDataValue) isDistressData() {}
func (distressDataList) isDistressData()  {}

type distressPacket struct {
	first, second distressData
}

func buildDistressPacket(input0, input1 string) (distressPacket, error) {
	data0, err := buildDistressData(input0)
	if err != nil {
		return distressPacket{}, err
	}
	data1, err := buildDistressData(input1)
	if err != nil {
		return distressPacket{}, err
	}
	return distressPacket{first: data0, second: data1}, nil

}

func buildDistressData(input string) (distressData, error) {
	var outer distressDataList
	collectedRunes := make([]rune, 0)
	listQueue := deque.New[*distressDataList]()
	listQueue.PushBack(&outer)
	for _, r := range input {
		switch r {
		case '[':
			// Create distressDataList in current location
			var newList distressDataList
			listQueue.PushBack(&newList)
		case ']':
			// Create distressDataValue in current location using collected runes,
			// then increment current location
			if len(collectedRunes) == 0 {
				// should only happen after exiting a list
				// ex: [1,[[2]],[3]]
				//            ^
				tailList := listQueue.PopBack()
				secondTailList := listQueue.PopBack()
				*secondTailList = append(*secondTailList, *tailList)
				listQueue.PushBack(secondTailList)
				continue
			}
			v, err := strconv.Atoi(string(collectedRunes))
			if err != nil {
				return outer, errors.New("failed to parse int: " + string(collectedRunes))
			}
			tailList := listQueue.PopBack()
			newValue := distressDataValue(v)
			*tailList = append((*tailList), newValue)
			secondTailList := listQueue.PopBack()
			*secondTailList = append(*secondTailList, *tailList)
			listQueue.PushBack(secondTailList)

			collectedRunes = make([]rune, 0)

		case ',':
			// Create distressDataValue in current location using collected runes
			if len(collectedRunes) == 0 {
				// should only happen after exiting a list
				// ex: [1,[2],[3]]
				//           ^
				continue
			}
			v, err := strconv.Atoi(string(collectedRunes))
			if err != nil {
				return outer, errors.New("failed to parse int: " + string(collectedRunes))
			}
			tailList := listQueue.PopBack()
			newValue := distressDataValue(v)
			*tailList = append((*tailList), newValue)
			listQueue.PushBack(tailList)

			collectedRunes = make([]rune, 0)
		default:
			// int
			collectedRunes = append(collectedRunes, r)
			// collect rune
		}
	}
	return outer, nil
}

type distressDataSlice []distressData

func (d distressDataSlice) Len() int      { return len(d) }
func (d distressDataSlice) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

// return true if d[i] < d[j]
func (d distressDataSlice) Less(i, j int) bool { return isSorted(d[i], d[j]) > 0 }

// Returns 1 if in correct order, -1 if in wrong order, or 0 if they are equal
func isSorted(p0, p1 distressData) int {
	switch l := p0.(type) {
	case distressDataValue:
		switch r := p1.(type) {
		case distressDataValue:
			if l < r {
				return 1
			} else if l > r {
				return -1
			} else {
				return 0
			}

		case distressDataList:
			res := isSorted(distressDataList{l}, r)
			if res != 0 {
				return res
			}

		default:
			panic("unknown type received")
		}

	case distressDataList:
		switch r := p1.(type) {
		case distressDataValue:
			res := isSorted(l, distressDataList{r})
			if res != 0 {
				return res
			}

		case distressDataList:
			leftLen := len(l)
			rightLen := len(r)
			for i := 0; i < leftLen && i < rightLen; i++ {
				res := isSorted(l[i], r[i])
				if res != 0 {
					return res
				}
			}
			if leftLen < rightLen {
				return 1
			} else if leftLen > rightLen {
				return -1
			} else {
				return 0
			}

		default:
			panic("unknown type received")
		}

	default:
		panic("unknown type received")
	}

	return 0
}
