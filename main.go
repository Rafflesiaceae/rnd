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

	fmt.Fprintln(f, "usage: <|yaml|0..9|9|-|> [count]")
	os.Exit(retCode)
}
func usageAndDie(retCode int, a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	usage(retCode)
}

func cli() {
	var err error

	args := os.Args[1:]

	for _, arg := range args {
		switch arg {
		case "-h", "--help":
			usage(0)
		}
	}

	if len(args) < 1 && len(args) > 2 {
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

	looksLikeRange := false
	{
		if inputsLineCount == 1 {
			dotCount := 0
			for _, v := range inputs {
				switch v {
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				case '.':
					dotCount++
				default:
					looksLikeRange = false
					break
				}
			}
			if dotCount != 2 {
				looksLikeRange = false
			} else {
				looksLikeRange = true
			}
		}
	}

	parseRange := func() error {
		var err error
		splitResult := strings.Split(inputs, "..")
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
	if looksLikeRange {
		err = parseRange()
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
		if choices == nil {
			result := start + (rand.Intn(end - start + 1))
			fmt.Println(result)
		} else {
			result := rand.Intn(len(choices))
			fmt.Println(choices[result])
		}
	}
}
