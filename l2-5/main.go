package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	A int
	B int
	C int
	c bool
	i bool
	v bool
	F bool
	n bool
)

var count int

func init() {
	flag.IntVar(&A, "A", 0, "")
	flag.IntVar(&B, "B", 0, "")
	flag.IntVar(&C, "C", 0, "")
	flag.BoolVar(&c, "c", false, "")
	flag.BoolVar(&i, "i", false, "")
	flag.BoolVar(&v, "v", false, "")
	flag.BoolVar(&F, "F", false, "")
	flag.BoolVar(&n, "n", false, "")

	flag.Parse()
}

func main() {
	fmt.Println(A, B, C, c, i, v, F, n)

	grepString := flag.Arg(0)
	filePath := flag.Arg(1)

	if A < 0 {
		fmt.Println("-A can't be < 0")
		os.Exit(0)
	}

	if B < 0 {
		fmt.Println("-B can't be < 0")
		os.Exit(0)
	}

	if C < 0 {
		fmt.Println("-C can't be < 0")
		os.Exit(0)
	}

	inBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("error reading file:", err)
		os.Exit(1)
	}

	inString := string(inBytes)

	lines := strings.Split(inString, "\n")

	// Реализация флага -i
	if i {
		inString = strings.ToLower(inString)
		grepString = strings.ToLower(grepString)
	}

	linesLowerCase := strings.Split(inString, "\n")

	var outputString string
	for i := range linesLowerCase {
		var isEscape bool
		if F {
			isEscape = strings.Contains(linesLowerCase[i], grepString)
		} else {
			isEscape, _ = regexp.MatchString(grepString, linesLowerCase[i])
		}

		if v {
			if !isEscape {
				if n {
					outputString += Flagn(i) + lines[i] + "\n"
					continue
				}
				outputString += lines[i] + "\n"
				count++
				continue
			}
			continue
		}
		if isEscape {
			count++

			if A > 0 {
				outputString += FlagA(A, i, lines)
				continue
			}
			if B > 0 {
				outputString += FlagB(B, i, lines)
				continue
			}
			if C > 0 {
				outputString += FlagC(C, i, lines)
				continue
			}
			if n {
				outputString += Flagn(i) + lines[i] + "\n"
				continue
			}
			outputString += lines[i] + "\n"
		}
	}

	// FlagA Реализация флага -c
	if c {
		fmt.Println(count)
		os.Exit(0)
	}

	fmt.Println(strings.TrimSpace(outputString))
}

func FlagA(aVal int, lineNumber int, lines []string) string {
	var output string
	if n {
		output = Flagn(lineNumber) + lines[lineNumber] + "\n"
	} else {
		output = lines[lineNumber] + "\n"
	}

	for i := 1; i <= aVal; i++ {
		if lineNumber+i >= len(lines) {
			output += "\n"
			continue
		}

		if n {
			output += Flagn(lineNumber+i) + lines[lineNumber+i] + "\n"
			continue
		}
		output += lines[lineNumber+i] + "\n"
	}
	return output
}

func FlagB(bVal int, lineNumber int, lines []string) string {
	output := ""
	for i := bVal; i >= 0; i-- {
		if lineNumber-i < 0 {
			output += "\n"
			continue
		}

		if n {
			output += Flagn(lineNumber-i) + lines[lineNumber-i] + "\n"
			continue
		}
		output += lines[lineNumber-i] + "\n"
	}
	return output
}

func FlagC(cVal int, lineNumber int, lines []string) string {
	output := ""
	for i := cVal; i >= 0; i-- {
		if lineNumber-i < 0 {
			output += "\n"
			continue
		}

		if n {
			output += Flagn(lineNumber-i) + lines[lineNumber-i] + "\n"
			continue
		}
		output += lines[lineNumber-i] + "\n"
		continue
	}

	for i := 1; i <= cVal; i++ {
		if lineNumber+i >= len(lines) {
			output += "\n"
			continue
		}

		if n {
			output += Flagn(lineNumber+i) + lines[lineNumber+i] + "\n"
			continue
		}
		output += lines[lineNumber+i] + "\n"
	}
	return output
}

func Flagn(lineNumber int) string {
	return fmt.Sprintf("%d:", lineNumber)
}
