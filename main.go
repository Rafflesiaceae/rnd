package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// still in disbelieve that we are actually doing this...

var (
	maxPositiveIntValue = int((^uint(0)) >> 1)

	start   int = 0
	end     int = 0
	choices []string

	count int = 1

	optJoin = false
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func die(retCode int, a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(retCode)
}
func dief(retCode int, format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(retCode)
}
func usage(retCode int) {
	f := os.Stderr
	if retCode == 0 {
		f = os.Stdout
	}

	fmt.Fprintln(f, "usage: [-j|--join] <|yaml|0..9|0-9|9|-|> [count]")
	os.Exit(retCode)
}
func usageAndDie(retCode int, a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	usage(retCode)
}

func cli() {
	var err error

	rawArgs := os.Args[1:]

	args := make([]string, 0)
	for _, arg := range rawArgs {
		switch arg {
		case "-h", "--help":
			usage(0)
		case "-j", "--join": // don't print newlines
			optJoin = true
		default:
			args = append(args, arg)
		}
	}

	if len(args) < 1 || len(args) > 2 {
		usageAndDie(1, "pass me 1-2 args")
	}

	if len(args) == 2 {
		count, err = strconv.Atoi(args[1])
		if err != nil {
			usageAndDie(1, fmt.Sprintf("failed to parse [count] as int → %s", args[1]))
		}
	}

	var inputs string

	{ // assign inputs
		firstArg := args[0]
		if firstArg == "-" { // read from stdin
			inputsRaw, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				usageAndDie(1, err)
			}
			inputs = string(inputsRaw)
		} else { // input is firstArg
			inputs = firstArg
		}
	}

	var inputsLineCount = 0
	{ // count lines
		scanner := bufio.NewScanner(strings.NewReader(inputs))
		for scanner.Scan() {
			inputsLineCount++
		}
		if err := scanner.Err(); err != nil {
			inputsLineCount = -1
		}
	}

	if inputsLineCount < 0 { // smells like binary
		usageAndDie(1, "we don't handle binary inputs for now")
	}

	var looksLikeRange string
	{
		if inputsLineCount == 1 {
			dotCount := 0
			dashCount := 0
			endashCount := 0
			for _, v := range inputs {
				switch v {
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				case '.':
					dotCount++
				case '-':
					dashCount++
				case '–':
					endashCount++
				default:
					looksLikeRange = ""
					break
				}
			}
			if dotCount == 2 && dashCount == 0 && endashCount == 0 {
				looksLikeRange = ".."
			} else if dotCount == 0 && dashCount == 1 && endashCount == 0 {
				looksLikeRange = "-"
			} else if dotCount == 0 && dashCount == 0 && endashCount == 1 {
				looksLikeRange = "–"
			} else {
				looksLikeRange = ""
			}
		}
	}

	parseRange := func(rangeSeperator string) error {
		var err error
		splitResult := strings.Split(inputs, rangeSeperator)
		if splitResult[0] > splitResult[1] {
			start, err = strconv.Atoi(splitResult[1])
			if err != nil {
				return err
			}
			end, err = strconv.Atoi(splitResult[0])
			if err != nil {
				return err
			}
		} else {
			start, err = strconv.Atoi(splitResult[0])
			if err != nil {
				return err
			}
			end, err = strconv.Atoi(splitResult[1])
			if err != nil {
				return err
			}
		}
		return nil
	}

	// attempt to parse a..b range
	if looksLikeRange != "" {
		err = parseRange(looksLikeRange)
		if err == nil {
			return
		}
	}

	// attempt to parse yaml, these will be choices
	// var parsedYaml interface{}
	err = yaml.Unmarshal([]byte(inputs), &choices)
	if err == nil {
		return
	}

	// attempt to parse as single int ← end
	end, err = strconv.Atoi(inputs)
	if err == nil {
		return
	}

	usageAndDie(1, "unable to read inputs")
}

func main() {
	cli()

	if start > end {
		dief(1, "start (%d) must not be greater than end (%d)", start, end)
	}

	for i := 0; i < count; i++ {
		var result string
		if choices == nil {
			result = strconv.Itoa(start + (rand.Intn(end - start + 1)))
		} else {
			result = choices[rand.Intn(len(choices))]
		}

		if optJoin {
			fmt.Print(result)
		} else {
			fmt.Println(result)
		}
	}
}
