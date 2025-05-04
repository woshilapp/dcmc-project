package main

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
	// proto "github.com/woshilapp/dcmc-project/protocol"
)

// Replace returns a copy of the string s with the first n non-overlapping instances of old replaced by new,
// excluding replacements within specified intervals. If old is empty, it matches at the beginning of the string
// and after each UTF-8 sequence. If n < 0, all valid instances are replaced.
func Replace(s, old, new string, n int, exclude [][]int) string {
	if old == new || n == 0 {
		return s
	}

	lenOld := len(old)
	positions := findAllIndex(s, old)
	filtered := filterPositions(positions, lenOld, exclude)
	m := len(filtered)
	if m == 0 {
		return s
	}

	if n < 0 || m < n {
		n = m
	}

	var b strings.Builder
	b.Grow(len(s) + n*(len(new)-lenOld))
	prev := 0
	for i := 0; i < n; i++ {
		pos := filtered[i]
		b.WriteString(s[prev:pos])
		b.WriteString(new)
		prev = pos + lenOld
	}
	b.WriteString(s[prev:])
	return b.String()
}

// ReplaceAll returns a copy of s with all non-overlapping instances of old replaced by new,
// excluding replacements within specified intervals.
func ReplaceAll(s, old, new string, exclude [][]int) string {
	return Replace(s, old, new, -1, exclude)
}

// findAllIndex finds all start positions of old in s.
func findAllIndex(s, old string) []int {
	var positions []int
	lenOld := len(old)
	start := 0
	if lenOld == 0 {
		positions = append(positions, 0)
		for start < len(s) {
			_, wid := utf8.DecodeRuneInString(s[start:])
			start += wid
			positions = append(positions, start)
		}
	} else {
		for {
			idx := strings.Index(s[start:], old)
			if idx == -1 {
				break
			}
			pos := start + idx
			positions = append(positions, pos)
			start = pos + lenOld
		}
	}
	return positions
}

// filterPositions removes positions where replacement overlaps with excluded intervals.
func filterPositions(positions []int, lenOld int, exclude [][]int) []int {
	var filtered []int
	for _, pos := range positions {
		excluded := false
		for _, ex := range exclude {
			exStart, exEnd := ex[0], ex[1]
			if pos < exEnd && exStart < pos+lenOld {
				excluded = true
				break
			}
		}
		if !excluded {
			filtered = append(filtered, pos)
		}
	}
	return filtered
}

func main2() {
	// proto.Encode(1, 2, 3, "asd")
	// proto.Run()

	s := `asBBda"BB"syasddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd`
	start := time.Now()
	result := ReplaceAll(s, "d", "G", [][]int{{22, 33}})
	ends := time.Since(start)
	fmt.Println(result, int64(ends/time.Millisecond))
	// 输出: asCCda"BB"sy
}
