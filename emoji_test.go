package emoji

import (
	"bytes"
	"testing"

	"github.com/levmv/emoji/internal/data"
)

type findCase struct {
	in     string
	offset int
	size   int
}

var findCases = []findCase{
	{"", 0, 0}, {"a", 0, 0}, {"foo", 0, 0}, {"—", 0, 0},
	{"©\ufe0f", 2, 3}, {"®\ufe0f", 2, 3}, {"™\ufe0f", 3, 3},
	{"©", 0, 0}, {"®", 0, 0}, {"™", 0, 0},
	{"🐨", 0, 4},
	{"🐨🐨", 0, 8},
	{"🐨🐨🐨", 0, 12},
	{"7️⃣7️⃣", 0, 7}, // can't detect sequence with ascii first rune
	{"7️⃣🐨", 0, 11},
	{"র‍্য", 0, 0}, // must not match \u200d if no emoji around
}

var removeTests = []string{
	"🐨", "",
	"—", "—",
	"foo📨bar", "foobar",
	"0️⃣", "", "1⃣", "",
	".📝.❇️.#️⃣.✳.", ".....",
	"ыℕ⊆₀næδέοგაიაγιγντ🛋️ρομ™€kΩ⍎⍕'", "ыℕ⊆₀næδέοგაიაγιγντρομ™€kΩ⍎⍕'",
	"zw🫱🏼‍🫲🏾🫱🏼‍🫲🏾j", "zwj",
	"🧔🏻‍♀.👱🏻‍♀️.👩🏼‍❤️‍💋‍👨🏿🇧🇴", "..",
}

func init() {

	for _, e := range data.AllEmojies {
		if e == "" {
			continue
		}
		findCases = append(findCases,
			findCase{e, 0, len(e)},
			findCase{"foo" + e, 3, len(e)},
		)
	}
}

func TestMissmatches(t *testing.T) {
	rmap := data.AllRunesMap()
	for r := '\u0000'; r < rune(0x10ffff); r++ {
		if _, ok := rmap[r]; ok {
			continue
		}

		if _, sz := find([]byte(string(r))); sz != 0 {
			t.Errorf(" matched %v (%U)", string(r), r)
		}
	}
}

func TestFind(t *testing.T) {
	for _, c := range findCases {
		o, s := find([]byte(c.in))
		if o != c.offset || s != c.size {
			t.Errorf("for '%s' got: %v %v, expect: %v %v", c.in, o, s, c.offset, c.size)
		}
	}
}

func TestRemove(t *testing.T) {
	for i := 0; i < len(removeTests); i += 2 {
		g := []byte(removeTests[i])
		e := []byte(removeTests[i+1])

		s := Remove(g)

		if !bytes.Equal(e, s) {
			t.Errorf("not matched `%s`(%v) and `%s`(%v)", e, e, s, s)
		}
	}
}
