package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Threshold struct {
	yellow int
	green int
}

type Params struct {
	label string
	threshold Threshold
	color string
	value string
}

func main() {
	source := flag.String("filename", "output.out", "File containing the tests output")
	label := flag.String("text", "Coverage", "Text on the left side of the badge")
	yellowThreshold := flag.Int("yellow", 30, "At what percentage does the badge becomes yellow instead of red")
	greenThreshold := flag.Int("green", 70, "At what percentage does the badge becomes green instead of yellow")
	color := flag.String("color", "", "Color of the badge - green/yellow/red")
	target := flag.String("target", "README.md", "Target file")
	value := flag.String("value", "", "Text on the right side of the badge")

	flag.Parse()

	params := &Params{
		*label,
		Threshold{*yellowThreshold, *greenThreshold},
		*color,
		*value,
	}

	err := generateBadge(*source, *target, params)

	if err != nil { log.Fatal(err) }
}

func generateBadge(source string, target string, params *Params) error {
	var coverage string
	var err error

	if params.value != "" {
		coverage = params.value
	} else {
		coverage, err = retrieveTotalCoverage(source)
	}

	if err != nil { return err }

	badgeColor := setColor(coverage, params.threshold.yellow, params.threshold.green, params.color)
	err = updateReadme(target, coverage, params.label, badgeColor)

	if err != nil { return err }

	fmt.Println("\033[0;36mGoBadge: Coverage badge updated to " + coverage + " in " + target + "\033[0m")

	return nil
}

func setColor(coverage string, yellowThreshold int, greenThreshold int, color string) string {
	coverageNumber, _ := strconv.ParseFloat(strings.Replace(coverage, "%", "", 1), 4)
	if color != "" { return color }
	if coverageNumber >= float64(greenThreshold) { return "brightgreen" }
	if coverageNumber >= float64(yellowThreshold) { return "yellow" }
	return "red"
}

func retrieveTotalCoverage(filename string) (string, error) {
	// Read coverage file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("\033[1;31mGoBadge: Error while opening the coverage file\033[0m")
		return "", err
	}
	defer file.Close()


	// split content by words and grab the last one (total percentage)
	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("\033[1;31mGoBadge: Error while reading the coverage file\033[0m")
		return "", err
	}
	words := strings.Fields(string(b))
	last := words[len(words)-1]

	return last, nil
}

func updateReadme(target string, coverage string, label string, color string) error {
	found := false
	encodedLabel := url.QueryEscape(label)
	encodedCoverage := url.QueryEscape(coverage)

	input, err := ioutil.ReadFile(target)
	if err != nil {
		fmt.Println("\033[1;31mGoBadge: Error while reading the target file\033[0m")
		return err
	}

	lines := strings.Split(string(input), "\n")
	newLine := "![" + label + "](https://img.shields.io/badge/" + encodedLabel + "-" + encodedCoverage + "-" + color + ")"

	for i, line := range lines {
		if strings.Contains(line, "![" + label + "](https://img.shields.io/badge/" + encodedLabel) {
			found = true
			lines[i] = newLine
		}
	}

	// If badge not found, insert the badge on line 2 (right after the title)
	if found == false {
		lines = append(lines, "")
		copy(lines[2:], lines[1:])
		lines[1] = newLine
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(target, []byte(output), 0644)
	if err != nil {
		fmt.Println("\033[1;31mGoBadge: Error while updating the target file\033[0m")
		return err
	}

	return nil
}
