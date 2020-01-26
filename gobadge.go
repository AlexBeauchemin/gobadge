package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	filename := flag.String("filename", "output.out", "File containing the tests output")
	label := flag.String("text", "Coverage", "Text on the left side of the badge")
	yellowThreshold := flag.Int("yellow", 30, "At what percentage does the badge becomes yellow instead of red")
	greenThreshold := flag.Int("green", 70, "At what percentage does the badge becomes green instead of yellow")
	color := flag.String("color", "", "Color of the badge - green/yellow/red")
	target := flag.String("target", "README.md", "Target file")

	flag.Parse()

	coverage := retrieveTotalCoverage(*filename)
	badgeColor := setColor(coverage, *yellowThreshold, *greenThreshold, *color)
	updateReadme(*target, coverage, *label, badgeColor)
}

func setColor(coverage string, yellowThreshold int, greenThreshold int, color string) string {
	coverageNumber, _ := strconv.ParseFloat(strings.Replace(coverage, "%", "", 1), 4)
	if color != "" { return color }
	if coverageNumber >= float64(greenThreshold) { return "brightgreen" }
	if coverageNumber >= float64(yellowThreshold) { return "yellow" }
	return "red"
}

func retrieveTotalCoverage(filename string) string {
	// Read coverage file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()


	// split content by words and grab the last one (total percentage)
	b, err := ioutil.ReadAll(file)
	words := strings.Fields(string(b))
	last := words[len(words)-1]

	return last
}

func updateReadme(target string, coverage string, label string, color string) {
	found := false
	encodedLabel := url.QueryEscape(label)
	encodedCoverage := url.QueryEscape(coverage)

	input, err := ioutil.ReadFile(target)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
}
