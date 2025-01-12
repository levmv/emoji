package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/levmv/emoji/internal/data"
)

const emojiTestURL = "https://unicode.org/Public/emoji/16.0/emoji-test.txt"

func fetchURL(url string) (body []byte, err error) {
	r, err := http.Get(url)
	if err != nil {
		return body, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return body, fmt.Errorf("status: %s", r.Status)
	}

	return io.ReadAll(r.Body)
}

func skipEmoji(r int64) bool {
	return r == 0x00a9 || r == 0x00ae || r == 0x2122
}

func gendata() {
	in, err := fetchURL(emojiTestURL)
	if err != nil {
		panic(err)
	}
	emojies := []string{}

	lines := strings.Split(string(in), "\n")
	for _, line := range lines {
		if line == "" || line[0] == '#' {
			continue
		}
		parts := strings.Split(line, ";")
		cps := strings.Split(strings.Trim(parts[0], " "), " ")

		var emoji []rune
		for _, s := range cps {
			i, err := strconv.ParseInt(s, 16, 32)
			if err != nil {
				panic(err)
			}
			if skipEmoji(i) {
				break
			}
			emoji = append(emoji, rune(i))
		}
		if len(emoji) > 0 {
			emojies = append(emojies, fmt.Sprintf("\"%s\"", string(emoji)))
		}

	}

	f, err := os.Create("../data/list.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("package data\n\nvar AllEmojies = []string{\n")

	for i, emoji := range emojies {
		f.Write([]byte(emoji))
		if (i+1)%30 == 0 {
			f.WriteString(",\n")
		} else {
			f.WriteString(",")
		}
	}
	f.WriteString("}")
}

func gentrie() {
	trie := Trie{Root: &Node{}}

	rmap := data.AllRunesMap()
	var list []rune

	for r := range rmap {
		list = append(list, r)
	}
	slices.Sort(list)

	for _, r := range list {
		trie.Put(r)
	}
	trie.Process()
	trie.PrintIndex()
	trie.PrintValues()
}

func help() {
	fmt.Println(`available commands:
		gen - download unicode test file and generates emoji list
		trie - generates trie data tables		
		`)
	os.Exit(1)
}

func main() {
	cmd := os.Args[1:]
	if len(cmd) == 0 {
		help()
	}

	switch cmd[0] {
	case "gen":
		gendata()
	case "trie":
		gentrie()
	default:
		help()
	}
}
