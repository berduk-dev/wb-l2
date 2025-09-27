package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	afterN       int
	beforeN      int
	contextN     int
	onlyCount    bool
	shouldIgnore bool
	shouldInvert bool
	shouldFixed  bool
	shouldNumber bool
)

var count int

func init() {
	flag.IntVar(&afterN, "A", 0, "печатать +N строк после совпадения")
	flag.IntVar(&beforeN, "B", 0, "печатать +N строк до совпадения")
	flag.IntVar(&contextN, "C", 0, "печатать ±N строк вокруг совпадения")
	flag.BoolVar(&onlyCount, "c", false, "количество строк")
	flag.BoolVar(&shouldIgnore, "i", false, "игнорировать регистр")
	flag.BoolVar(&shouldInvert, "v", false, "вместо совпадения, исключать")
	flag.BoolVar(&shouldFixed, "F", false, "точное совпадение со строкой, не паттерн")
	flag.BoolVar(&shouldNumber, "n", false, "напечатать номер строки")

	flag.Parse()
}

func main() {
	grepString := flag.Arg(0)
	filePath := flag.Arg(1)

	if afterN < 0 {
		fmt.Println("-afterN can't be < 0")
		os.Exit(0)
	}

	if beforeN < 0 {
		fmt.Println("-beforeN can't be < 0")
		os.Exit(0)
	}

	if contextN < 0 {
		fmt.Println("-contextN can't be < 0")
		os.Exit(0)
	}

	inBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("error reading file:", err)
		os.Exit(1)
	}

	inString := string(inBytes)

	lines := strings.Split(inString, "\n")

	// реализация флага -shouldIgnore
	if shouldIgnore {
		inString = strings.ToLower(inString)
		grepString = strings.ToLower(grepString)
	}

	linesLowerCase := strings.Split(inString, "\n")

	var outputString string
	var isEscape bool
	for i := range linesLowerCase {
		if shouldFixed {
			isEscape = strings.Contains(linesLowerCase[i], grepString)
		} else {
			isEscape, _ = regexp.MatchString(grepString, linesLowerCase[i])
		}

		// реализация флага -onlyCount
		if shouldInvert {
			if isEscape {
				continue
			}

			if shouldNumber {
				outputString += AddLineNumber(i) + highlight(lines[i], grepString, shouldFixed) + "\n"
				count++
				continue
			}
			outputString += highlight(lines[i], grepString, shouldFixed) + "\n"
			count++
			continue
		}
		if isEscape {
			if onlyCount {
				count++
				continue
			}

			if afterN > 0 {
				outputString += FlagA(afterN, i, grepString, lines, shouldFixed)
				continue
			}
			if beforeN > 0 {
				outputString += FlagB(beforeN, i, grepString, lines, shouldFixed)
				continue
			}
			if contextN > 0 {
				outputString += FlagC(contextN, i, grepString, lines, shouldFixed)
				continue
			}

			if shouldNumber {
				outputString += AddLineNumber(i) + highlight(lines[i], grepString, shouldFixed) + "\n"
			} else {
				outputString += highlight(lines[i], grepString, shouldFixed) + "\n"
			}
		}
	}

	// реализация флага -onlyCount
	if onlyCount {
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
		// как есть
		re = regexp.MustCompile("(?i)" + grepString)
	}
	return re.ReplaceAllString(line, "\033[31m$0\033[0m")
}

func FlagA(aVal int, lineNumber int, grepString string, lines []string, fixed bool) string {
	var output string
	if shouldNumber {
		output = AddLineNumber(lineNumber) + highlight(lines[lineNumber], grepString, fixed) + "\n"
	} else {
		output = highlight(lines[lineNumber], grepString, fixed) + "\n"
	}

	for i := 1; i <= aVal; i++ {
		if lineNumber+i >= len(lines) {
			output += "\n"
			continue
		}

		if shouldNumber {
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

		if shouldNumber {
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

		if shouldNumber {
			output += AddLineNumber(lineNumber-i) + highlight(lines[lineNumber-i], grepString, fixed) + "\n"
			continue
		}
		output += highlight(lines[lineNumber-i], grepString, fixed) + "\n"
		continue
	}

	// afterN line
	for i := 1; i <= cVal; i++ {
		if lineNumber+i >= len(lines) {
			output += "\n"
			continue
		}

		if shouldNumber {
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
