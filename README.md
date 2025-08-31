# emoji remover
[![Go](https://github.com/levmv/emoji/actions/workflows/go.yml/badge.svg)](https://github.com/levmv/emoji/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/levmv/emoji)](https://goreportcard.com/report/github.com/levmv/emoji)

Simple library to remove all emoji from string.

### Basic Usage

```go
import "github.com/levmv/emoji"

// Remove emojis from string
text := "I love coding! 🚀💻"
result := emoji.Remove([]byte(text))
// result: "I love coding! "
```


## Performance

- **Time Complexity**: O(n) 
- **Memory Usage**: ~1KB static lookup tables  
- **Memory Efficiency**: When no emojis are found, returns the original slice. 
- **UTF-8 Optimized**: Direct byte parsing, no rune conversion



## Important Considerations
-  The input must be valid UTF-8. This is by design - emoji removal is typically not the first text processing step, so the string should already be validated.

- `©️`, `®️` and `™️` are **not** considered emojis and will be preserved(but `0xfe0f` rune removed if present).