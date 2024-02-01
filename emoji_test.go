package emoji

import (
	"slices"
	"testing"

	"github.com/levmv/emoji/internal/data"
)

var findCases = []struct {
	in     string
	offset int
	size   int
}{
	{"", 0, 0}, {"a", 0, 0}, {"foo", 0, 0}, {"—", 0, 0},
	{"🐨", 0, 4},
	{"🐨🐨", 0, 8},
	{"🐨🐨🐨", 0, 12},
	{"foo⏰", 3, 3}, {"foo⏰ ⏰", 3, 3},
	{"foo⏰⏰", 3, 6},
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

var utfSizeCases = []struct {
	in     string
	expect int
}{
	{"a", 1},
	{"Г", 2}, {"П", 2},
	{"—", 3},
	{"🐨", 4},
}

func TestUTFSize(t *testing.T) {
	for _, c := range utfSizeCases {
		sz, _ := lookup([]byte(c.in))
		if sz != c.expect {
			t.Errorf("for '%v' expected size %v, got %v", c.in, c.expect, sz)
		}
	}
}

func TestAllRunes(t *testing.T) {
	rmap := make(map[rune]bool, len(data.AllEmojiRunes))
	for _, c := range data.AllEmojiRunes {
		rmap[c] = true

		if _, ok := lookup([]byte(string(c))); !ok {
			t.Errorf("not matched %U (%v)", c, string(c))
		}
	}
	for r := '\u1000'; r < rune(0xea07f); r++ {
		if _, ok := rmap[r]; ok {
			continue
		}

		if _, ok := lookup([]byte(string(r))); ok {
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

		if slices.Compare(e, s) != 0 {
			t.Errorf("not matched `%s`(%v) and `%s`(%v)", e, e, s, s)
		}
	}
}
