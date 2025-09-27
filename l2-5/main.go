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

	// реализация флага -i
	if i {
		inString = strings.ToLower(inString)
		grepString = strings.ToLower(grepString)
	}

	linesLowerCase := strings.Split(inString, "\n")

	var outputString string
	var isEscape bool
	for i := range linesLowerCase {
		if F {
			isEscape = strings.Contains(linesLowerCase[i], grepString)
		} else {
			isEscape, _ = regexp.MatchString(grepString, linesLowerCase[i])
		}

		// реализация флага -c
		if v {
			if !isEscape {
				if n {
					outputString += AddLineNumber(i) + highlight(lines[i], grepString, F) + "\n"
					count++
					continue
				}
				outputString += highlight(lines[i], grepString, F) + "\n"
				count++
				continue
			}
			continue
		}
		if isEscape {
			if c {
				count++
				continue
			}

			if A > 0 {
				outputString += FlagA(A, i, grepString, lines, F)
				continue
			}
			if B > 0 {
				outputString += FlagB(B, i, grepString, lines, F)
				continue
			}
			if C > 0 {
				outputString += FlagC(C, i, grepString, lines, F)
				continue
			}

			if n {
				outputString += AddLineNumber(i) + highlight(lines[i], grepString, F) + "\n"
			} else {
				outputString += highlight(lines[i], grepString, F) + "\n"
			}
		}
	}

	// реализация флага -c
	if c {
		fmt.Println(count)
		os.Exit(0)
	}

	fmt.Println(strings.TrimSpace(outputString))
}

func highlight(line, grepString string, fixed bool) string {
	var re *regexp.Regexp
	if fixed {
		// ищем буквально, без regexp-метасимволов
		re = regexp.MustCompile("(?i)" + regexp.QuoteMeta(grepString))
	} else {
		// как есть (регулярка)
		re = regexp.MustCompile("(?i)" + grepString)
	}
	return re.ReplaceAllString(line, "\033[31m$0\033[0m")
}

func FlagA(aVal int, lineNumber int, grepString string, lines []string, fixed bool) string {
	var output string
	if n {
		output = AddLineNumber(lineNumber) + highlight(lines[lineNumber], grepString, fixed) + "\n"
	} else {
		output = highlight(lines[lineNumber], grepString, fixed) + "\n"
	}

	for i := 1; i <= aVal; i++ {
		if lineNumber+i >= len(lines) {
			output += "\n"
			continue
		}

		if n {
			output += AddLineNumber(lineNumber+i) + highlight(lines[lineNumber+i], grepString, fixed) + "\n"
			continue
		}
		output += highlight(lines[lineNumber+i], grepString, fixed) + "\n"
	}
	return output
}

func FlagB(bVal int, lineNumber int, grepString string, lines []string, fixed bool) string {
	output := ""
	for i := bVal; i >= 0; i-- {
		if lineNumber-i < 0 {
			output += "\n"
			continue
		}

		if n {
			output += AddLineNumber(lineNumber-i) + highlight(lines[lineNumber-i], grepString, fixed) + "\n"
			continue
		}
		output += highlight(lines[lineNumber-i], grepString, fixed) + "\n"
	}
	return output
}

func FlagC(cVal int, lineNumber int, grepString string, lines []string, fixed bool) string {
	// before line and line
	output := ""
	for i := cVal; i >= 0; i-- {
		if lineNumber-i < 0 {
			output += "\n"
			continue
		}

		if n {
			output += AddLineNumber(lineNumber-i) + highlight(lines[lineNumber-i], grepString, fixed) + "\n"
			continue
		}
		output += highlight(lines[lineNumber-i], grepString, fixed) + "\n"
		continue
	}

	// after line
	for i := 1; i <= cVal; i++ {
		if lineNumber+i >= len(lines) {
			output += "\n"
			continue
		}

		if n {
			output += AddLineNumber(lineNumber+i) + highlight(lines[lineNumber+i], grepString, fixed) + "\n"
			continue
		}
		output += highlight(lines[lineNumber+i], grepString, fixed) + "\n"
	}
	return output
}

func AddLineNumber(lineNumber int) string {
	return fmt.Sprintf("%d:", lineNumber+1)
}
