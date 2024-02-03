# emoji remover
[![Go](https://github.com/levmv/emoji/actions/workflows/go.yml/badge.svg)](https://github.com/levmv/emoji/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/levmv/emoji)](https://goreportcard.com/report/github.com/levmv/emoji)

Simple library to remove all emoji from string. Pretty fast and with 
very moderate memory consumption (total size of used static tables ~1kb)

Important to know:
- Input must be valid UTF8 (assumption is that emoji removal is most likely not the first thing done on string, so it's must be already validated and so we can cut corners and gain some speed)
- `©️`, `®️` and `™️` is not recognised as emoji and therefore are not removed (but `0xfe0f` rune removed if present)