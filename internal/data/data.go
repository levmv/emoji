// package data provides emoji lists for trie gen and for tests
// other than that it's not used in the library
package data

func AllRunesMap() map[rune]struct{} {
	emmap := make(map[rune]struct{}, 1400)

	for _, emoji := range AllEmojies {
		for _, r := range emoji {
			if r > 0xae && r != 0x2122 {
				emmap[r] = struct{}{}
			}
		}
	}
	return emmap
}
