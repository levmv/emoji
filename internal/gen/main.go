package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/levmv/emoji/internal/data"
)

const emojiTestURL = "https://unicode.org/Public/emoji/15.0/emoji-test.txt"

const dataTpl = `
package data
// Size: {{ .Size }}
var AllEmojiRunes = []rune { 
{{ .List }} 
}`

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

func iterateFirstCol(in string, f func(string)) {
	lines := strings.Split(in, "\n")

	for _, line := range lines {
		if line == "" || line[0] == '#' {
			continue
		}
		parts := strings.Split(line, ";")
		f(strings.Trim(parts[0], " "))
	}
}

func formatList(list []rune) string {
	var (
		r      strings.Builder
		strlen int
	)
	for _, emoji := range list {
		if strlen > 70 {
			r.WriteString("\n")
			strlen = 0
		}
		str := fmt.Sprintf("0x%x,", emoji)
		r.WriteString(str)
		strlen += len(str)
	}
	return r.String()
}

func gendata() {
	in, err := fetchURL(emojiTestURL)
	if err != nil {
		panic(err)
	}

	emmap := make(map[rune]bool, 1400)

	iterateFirstCol(string(in), func(text string) {
		cps := strings.Split(text, " ")

		//		var emoji []rune
		for _, s := range cps {
			i, err := strconv.ParseInt(s, 16, 32)
			if err != nil {
				panic(err)
			}
			//			emoji = append(emoji, rune(i))

			if i > 0xae && i != 0x2122 {
				emmap[rune(i)] = true
			}
		}
	})
	var (
		size int
		list []rune
	)
	for e := range emmap {
		size += len(string(e))
		list = append(list, e)
	}
	slices.Sort(list)

	t := template.Must(template.New("tpl").Parse(dataTpl))

	f, err := os.Create("../data/data.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_ = t.Execute(f, struct {
		Size int
		List string
	}{size, formatList(list)})
}

func gentrie() {
	trie := Trie{Root: &Node{}}
	for _, r := range data.AllEmojiRunes {
		trie.Put(r, 1)
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
