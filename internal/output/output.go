package output

import (
	"fmt"
	"strings"
)

func PrintAligned(lines []string, longestKey int) {
	for _, line := range lines {
		parts := strings.SplitN(line, " - ", 2)
		if len(parts) == 2 {
			fmt.Printf("%-*s %s\n", longestKey, parts[0], parts[1])
		} else {
			fmt.Println(line)
		}
	}
}

func FindLongest(items []string) int {
	longest := 0
	for _, item := range items {
		parts := strings.SplitN(item, " - ", 2)
		if len(parts) > 0 && len(parts[0]) > longest {
			longest = len(parts[0])
		}
	}
	return longest
}
