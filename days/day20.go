package days

import (
	"errors"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
)

type Day20Solver struct {
}

func (d Day20Solver) SolvePartA(puzzleInput string) (string, fyne.CanvasObject, error) {
	// Plan:
	// Create a datastructure with a linkled list and a slice of pointers to the linked list items
	// use the slice to iterate through the items in original order. Use the linked list to maintain
	// and update the final order. slice = O(1) access to next item to move, linked list = O(n) for reordering

	m, err := buildMixEncryption(puzzleInput)
	if err != nil {
		return "", nil, err
	}
	m.mix()

	one := m.zeroNode.getFollowingNode(1000).value
	two := m.zeroNode.getFollowingNode(2000).value
	three := m.zeroNode.getFollowingNode(3000).value

	solstr := strconv.Itoa(one + two + three)
	return solstr, nil, nil
}

func (d Day20Solver) SolvePartB(puzzleInput string) (string, fyne.CanvasObject, error) {
	m, err := buildMixEncryption(puzzleInput)
	if err != nil {
		return "", nil, err
	}
	for i := range m.originalOrder {
		m.originalOrder[i].value *= 811589153
	}
	for i := 0; i < 10; i++ {
		m.mix()
	}

	one := m.zeroNode.getFollowingNode(1000).value
	two := m.zeroNode.getFollowingNode(2000).value
	three := m.zeroNode.getFollowingNode(3000).value

	solstr := strconv.Itoa(one + two + three)
	return solstr, nil, nil
}

type mixEncryptNode struct {
	value int
	prev  *mixEncryptNode
	next  *mixEncryptNode
}

type mixEncryption struct {
	originalOrder []mixEncryptNode
	linkedHead    *mixEncryptNode
	zeroNode      *mixEncryptNode
}

func buildMixEncryption(input string) (mixEncryption, error) {
	lines := strings.Split(input, "\n")
	var m mixEncryption
	m.originalOrder = make([]mixEncryptNode, len(lines))
	for i, line := range lines {
		val, err := strconv.Atoi(line)
		if err != nil {
			return m, nil
		}
		m.originalOrder[i].value = val
		if val == 0 {
			if m.zeroNode != nil {
				return m, errors.New("found multiple zero nodes, second at index: " + strconv.Itoa(i))
			}
			m.zeroNode = &m.originalOrder[i]
		}
	}
	for i := 1; i < len(m.originalOrder); i++ {
		m.originalOrder[i-1].next = &m.originalOrder[i]
		m.originalOrder[i].prev = &m.originalOrder[i-1]
	}
	m.originalOrder[0].prev = &m.originalOrder[len(m.originalOrder)-1]
	m.originalOrder[len(m.originalOrder)-1].next = &m.originalOrder[0]
	m.linkedHead = &m.originalOrder[0]
	return m, nil
}

func (m *mixEncryption) mix() {
	// fmt.Println(m.String())
	for i := range m.originalOrder {
		src := &m.originalOrder[i]
		travel := src.value % (len(m.originalOrder) - 1)
		// travel := src.value
		// if travel > 0 {
		// 	travel += 20 * (len(m.originalOrder) - 1)
		// } else {
		// 	travel -= 20 * (len(m.originalOrder) - 1)
		// }
		if travel == 0 {
			continue
		}

		// remove source
		src.prev.next = src.next
		src.next.prev = src.prev
		if m.linkedHead == src {
			m.linkedHead = src.next
		}

		dest := src.getFollowingNode(travel)
		if src == dest {
			continue
		}

		// "add" src
		src.next = dest.next
		src.prev = dest
		dest.next.prev = src
		dest.next = src
	}
}

func (m mixEncryption) String() string {
	var sb strings.Builder
	curr := m.linkedHead
	for {
		sb.WriteString(strconv.Itoa(curr.value))
		curr = curr.next
		if curr == m.linkedHead {
			break
		}
		sb.WriteString(", ")
	}

	return sb.String()
}

func (m *mixEncryptNode) getFollowingNode(offset int) *mixEncryptNode {
	res := m
	if offset > 0 {
		for i := 0; i < offset; i++ {
			res = res.next
		}
	} else if offset < 0 {
		// need <= because movement is relative to the "prev" node so negative needs a +1
		for i := 0; i <= -offset; i++ {
			res = res.prev
		}
	}
	return res
}
