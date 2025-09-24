package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	k int
	n bool
	r bool
	u bool
)

func init() {
	flag.IntVar(&k, "k", -1, "указание колонки для сортировки")
	flag.BoolVar(&n, "n", false, "сортировать по числовому значению")
	flag.BoolVar(&r, "r", false, "сортировать в обратном порядке")
	flag.BoolVar(&u, "u", false, "не выводить повторяющиеся строки")

	flag.Parse()
}

func main() {
	fmt.Println(k, n, r, u)
	filePath := flag.Arg(0)

	inBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("error reading file:", err)
		os.Exit(1)
	}

	lines := strings.Split(string(inBytes), "\n")

	sort.Slice(lines, func(i int, j int) bool {
		wordI := strings.Split(lines[i], " ")
		wordJ := strings.Split(lines[j], " ")

		// Реализация флага -k
		if k >= 0 {
			if r {
				return wordI[k] > wordJ[k]
			}

			return wordI[k] < wordJ[k]
		}

		// Реализация флага -n
		if n {
			num := 0
			for i, el := range wordI {
				_, err := strconv.Atoi(el)
				if err == nil {
					num = i
					break
				}
			}
			numI, _ := strconv.Atoi(wordI[num])
			numJ, _ := strconv.Atoi(wordJ[num])

			if r {
				return numI > numJ
			}

			return numI < numJ
		}

		if r {
			return lines[i] > lines[j]
		}

		return lines[i] < lines[j]
	})

	// Проверка на уникальные строки и вывод
	var unicLines []string
	if u {
		for i := range lines {
			if i > 0 {
				if lines[i] == lines[i-1] {
					continue
				}
			}
			unicLines = append(unicLines, lines[i])
		}
		outBytes := []byte(strings.Join(unicLines, "\n"))
		err = os.WriteFile("sorted.txt", outBytes, 0644)
	} else { // Если флаг -u не передан
		outBytes := []byte(strings.Join(lines, "\n"))
		err = os.WriteFile("sorted.txt", outBytes, 0644)
	}

	if err != nil {
		fmt.Println("error writing file:", err)
		os.Exit(1)
	}
}
