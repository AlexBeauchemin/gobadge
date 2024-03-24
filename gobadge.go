package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Threshold struct {
	yellow int
	green  int
}

type Params struct {
	label     string
	threshold Threshold
	color     string
	value     string
	link      string
}

func main() {
	source := flag.String("filename", "output.out", "File containing the tests output")
	label := flag.String("text", "Coverage", "Text on the left side of the badge")
	yellowThreshold := flag.Int("yellow", 30, "At what percentage does the badge becomes yellow instead of red")
	greenThreshold := flag.Int("green", 70, "At what percentage does the badge becomes green instead of yellow")
	color := flag.String("color", "", "Color of the badge - green/yellow/red")
	target := flag.String("target", "README.md", "Target file")
	value := flag.String("value", "", "Text on the right side of the badge")
	link := flag.String("link", "", "Link the badge goes to")

	flag.Parse()

	params := &Params{
		*label,
		Threshold{*yellowThreshold, *greenThreshold},
		*color,
		*value,
		*link,
	}

	err := generateBadge(*source, *target, params)

	if err != nil {
		log.Fatal(err)
	}
}

func generateBadge(source string, target string, params *Params) error {
	var coverage string
	var err error

	if params.value != "" {
		coverage = params.value
	} else {
		coverage, err = retrieveTotalCoverage(source)
	}

	if err != nil {
		return err
	}

	badgeColor := setColor(coverage, params.threshold.yellow, params.threshold.green, params.color)
	err = updateReadme(target, coverage, params.label, badgeColor, params.link)

	if err != nil {
		return err
	}

	fmt.Println("\033[0;36mGoBadge: Coverage badge updated to " + coverage + " in " + target + "\033[0m")

	return nil
}

func setColor(coverage string, yellowThreshold int, greenThreshold int, color string) string {
	coverageNumber, _ := strconv.ParseFloat(strings.Replace(coverage, "%", "", 1), 64)
	if color != "" {
		return color
	}
	if coverageNumber >= float64(greenThreshold) {
		return "brightgreen"
	}
	if coverageNumber >= float64(yellowThreshold) {
		return "yellow"
	}
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
	b, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("\033[1;31mGoBadge: Error while reading the coverage file\033[0m")
		return "", err
	}
	words := strings.Fields(string(b))
	last := words[len(words)-1]

	return last, nil
}

func updateReadme(target string, coverage string, label string, color string, link string) error {
	encodedLabel := url.QueryEscape(label)
	encodedCoverage := url.QueryEscape(coverage)

	input, err := os.ReadFile(target)
	if err != nil {
		fmt.Println("\033[1;31mGoBadge: Error while reading the target file\033[0m")
		return err
	}

	// Possible regex exprs with and without link
	// Playground: https://goplay.tools/snippet/GWvkx43QndT
	badgeRegexes := []*regexp.Regexp{
		regexp.MustCompile(`!\[(\w+)\]\(https:\/\/img\.shields\.io\/badge\/(\w+)-([\d\.%]+)-(\w+)\)`), 
		regexp.MustCompile(`\[!\[(\w+)\]\(https:\/\/img\.shields\.io\/badge\/(\w+)-([\d\.%]+)-(\w+)\)\]\((.*)\)`),
	}
	// badgeRegex := regexp.MustCompile(`!\[(\w+)\]\(https:\/\/img\.shields\.io\/badge\/(\w+)-([\d\.%]+)-(\w+)\)`)
	// if link != "" {
	// 	badgeRegex = regexp.MustCompile(`\[!\[(\w+)\]\(https:\/\/img\.shields\.io\/badge\/(\w+)-([\d\.%]+)-(\w+)\)\]\((.*)\)`)
	// }

	newBadge := "![" + label + "](https://img.shields.io/badge/" + encodedLabel + "-" + encodedCoverage + "-" + color + ")"
	if link != "" {
		newBadge = "[![" + label + "](https://img.shields.io/badge/" + encodedLabel + "-" + encodedCoverage + "-" + color + ")](" + link + ")"
	}

	// Check if badge is already in README. If matches, replace it with new badge
	var output string
	var found = false
	for _, re := range badgeRegexes {
		if re.MatchString(string(input)) {
			output = re.ReplaceAllString(string(input), newBadge)
			found = true
			goto outside
		}
	}

outside:
	// If no matches found for regex exprs, it means there is no badge in README
	// If badge not found, insert the badge on line 2 (right after the title)
	if !found {
		lines := strings.Split(string(input), "\n")
		lines = append(lines, "")
		copy(lines[2:], lines[1:])
		lines[1] = newBadge
		output = strings.Join(lines, "\n")
	}

	err = os.WriteFile(target, []byte(output), 0644)
	if err != nil {
		fmt.Println("\033[1;31mGoBadge: Error while updating the target file\033[0m")
		return err
	}

	return nil
}
