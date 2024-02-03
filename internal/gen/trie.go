package main

import (
	"fmt"
)

type Trie struct {
	Root *Node

	Values   []uint64
	Index    [][64]byte
	CurIndex int
}

type Node struct {
	Children [64]*Node
	Values   [64]byte
	ValIndex int
}

func (t *Trie) Put(r rune) {
	s := string(r)

	node := t.Root
	var s0 byte

	for ; len(s) > 1; s = s[1:] {
		s0 = s[0]
		ndx := s0 % 64
		c := node.Children[ndx]
		if c == nil {
			c = &Node{}
			node.Children[ndx] = c
		}
		node = c
	}
	node.Values[s[0]-0x80] = 1
}

func (t *Trie) ProcessValues(node *Node) {
	for _, c := range &node.Children {
		if c != nil {
			t.ProcessValues(c)
		}
	}

	if IsEmpty(node.Values) {
		return
	}

	var pvalue uint64
	for i, v := range node.Values {
		if v != 0 {
			pvalue |= uint64(1) << i
		}
	}
	t.Values = append(t.Values, pvalue)
	node.ValIndex = len(t.Values)
}

func (t *Trie) Calc(node *Node) int {
	count := 0
	for _, c := range &node.Children {
		if c != nil {
			count += t.Calc(c) + 1
		}
	}
	return count
}

func IsEmpty(list [64]byte) bool {
	for _, v := range list {
		if v != 0 {
			return false
		}
	}
	return true
}

const (
	s2 = 0x02 // size 2
	s3 = 0x03 // size 3
	s4 = 0x04 // size 4
)

var utfSizes = [64]byte{
	00, 00, s2, s2, s2, s2, s2, s2, s2, s2, s2, s2, s2, s2, s2, s2, // 0xC0-0xCF
	s2, s2, s2, s2, s2, s2, s2, s2, s2, s2, s2, s2, s2, s2, s2, s2, // 0xD0-0xDF
	s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, // 0xE0-0xEF
	s4, s4, s4, s4, s4, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, // 0xF0-0xFF
}

func (t *Trie) ProcessIndex(node *Node) int {
	var index [64]byte

	for i, c := range &node.Children {
		if c != nil {
			index[i] = byte(t.ProcessIndex(c))
		}
	}

	if node.ValIndex != 0 {
		return node.ValIndex - 1
	}

	t.Index = append(t.Index, index)
	return len(t.Index) - 3
}

func (t *Trie) Process() {
	// first block in values and index are always zero, to simplify lookup code
	t.Values = append(t.Values, 0)
	t.ProcessValues(t.Root)
	t.Index = append(t.Index, [64]byte{}, [64]byte{}, [64]byte{}, utfSizes)

	for i, c := range &t.Root.Children {
		if c != nil {
			t.Index[3][i] |= (byte(t.ProcessIndex(c)) << 3)
		}
	}
}

func (t *Trie) PrintIndex() {
	var slen int
	fmt.Printf("var indexTable = [%d]byte{\n", len(t.Index)*64)
	for i := range t.Index {
		for j := 0; j < 64; j++ {
			if t.Index[i][j] != 0 {
				s := fmt.Sprintf("0x%03x:0x%02x, ", i*64+j, t.Index[i][j])
				fmt.Print(s)
				slen += len(s)
				if slen > 70 {
					fmt.Println("")
					slen = 0
				}
			}
		}
	}
	fmt.Printf("\n}\n")
}

func (t *Trie) PrintValues() {
	var slen int
	fmt.Printf("var valuesTable = [%d]uint64{\n", len(t.Values))
	for _, v := range t.Values {
		s := fmt.Sprintf("0x%016x,", v)
		fmt.Print(s)
		slen += len(s)
		if slen > 75 {
			fmt.Println("")
			slen = 0
		}
	}
	fmt.Printf("}\n")
}
