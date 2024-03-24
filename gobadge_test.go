package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createFile(path string, content string) {
	var _, err = os.Stat(path)

	// If file exists, delete it
	if err == nil {
		deleteFile(path)
	}

	file, _ := os.Create(path)
	file.WriteString(content)
	file.Sync()

	defer file.Close()
}

func deleteFile(path string) {
	os.Remove(path)
}

func validateFileContent(t *testing.T, path string, content string) {
	input, _ := ioutil.ReadFile(path)
	assert.Equal(t, content, string(input))
}

func TestGenerateBadge(t *testing.T) {
	var err error

	// Create a badge and use value from the output file as coverage
	createFile("test.md", "## Header\nDescription...\n")
	createFile("test.txt", "github.com/AlexBeauchemin/go-coverage-badge/gobadge.go:53:	updateReadme\n80.0%total:                                                          (statements)            33.3%")
	err = generateBadge("test.txt", "test.md", &Params{"Coverage", Threshold{50, 70}, "", "", ""})
	assert.Equal(t, nil, err)
	validateFileContent(t, "test.md", "## Header\n![Coverage](https://img.shields.io/badge/Coverage-33.3%25-red)\nDescription...\n")

	// Create a badge and use the provided value as coverage
	createFile("test.md", "## Header\nDescription...\n")
	err = generateBadge("unknown.out", "test.md", &Params{"Coverage", Threshold{50, 70}, "", "55%", ""})
	assert.Equal(t, nil, err)
	validateFileContent(t, "test.md", "## Header\n![Coverage](https://img.shields.io/badge/Coverage-55%25-yellow)\nDescription...\n")

	// Update the badge if it already exists
	err = generateBadge("unknown.out", "test.md", &Params{"Coverage", Threshold{50, 70}, "green", "56%", ""})
	assert.Equal(t, nil, err)
	validateFileContent(t, "test.md", "## Header\n![Coverage](https://img.shields.io/badge/Coverage-56%25-green)\nDescription...\n")

	// error if invalid source file
	err = generateBadge("unknown.out", "test.md", &Params{"Coverage", Threshold{50, 70}, "", "", ""})
	assert.Error(t, err)

	// error if invalid target file
	err = generateBadge("coverage.out", "unknown.md", &Params{"Coverage", Threshold{50, 70}, "", "", ""})
	assert.Error(t, err)
}

func TestSetColor(t *testing.T) {
	color := setColor("10%", 30, 70, "")
	assert.Equal(t, "red", color)
	color = setColor("30%", 30, 70, "")
	assert.Equal(t, "yellow", color)
	color = setColor("35.51%", 30, 70, "")
	assert.Equal(t, "yellow", color)
	color = setColor("100%", 30, 70, "")
	assert.Equal(t, "brightgreen", color)
	color = setColor("100%", 30, 70, "blue")
	assert.Equal(t, "blue", color)
}

func TestRetrieveCoverage(t *testing.T) {
	createFile("test.txt", "github.com/AlexBeauchemin/go-coverage-badge/gobadge.go:53:	updateReadme\n80.0%total:                                                          (statements)            33.3%")
	coverage, err := retrieveTotalCoverage("test.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, "33.3%", coverage)
	deleteFile("test.txt")

	coverage, err = retrieveTotalCoverage("test2.txt")
	assert.Equal(t, "open test2.txt: no such file or directory", err.Error())
	assert.Equal(t, "", coverage)
}

func TestUpdateReadme(t *testing.T) {
	var err error
	target := "test.md"

	createFile(target, "## A Title\nA description\n")
	err = updateReadme(target, "50.01%", "Coverage", "green", "")
	validateFileContent(t, target, "## A Title\n![Coverage](https://img.shields.io/badge/Coverage-50.01%25-green)\nA description\n")
	assert.Equal(t, nil, err)
	deleteFile(target)

	createFile(target, "## A Title\nA description\n")
	err = updateReadme(target, "50.01%", "Coverage", "green", "https://www.cnn.com")
	validateFileContent(t, target, "## A Title\n[![Coverage](https://img.shields.io/badge/Coverage-50.01%25-green)](https://www.cnn.com)\nA description\n")
	assert.Equal(t, nil, err)
	deleteFile(target)

	// Test updating badge with link to badge without link
	createFile(target, "## A Title\nA description\nContent [![Coverage](https://img.shields.io/badge/Coverage-50.01%25-green)](https://www.cnn.com)")
	err = updateReadme(target, "40.01%", "Coverage", "green", "")
	validateFileContent(t, target, "## A Title\nA description\nContent ![Coverage](https://img.shields.io/badge/Coverage-40.01%25-green)")
	assert.Equal(t, nil, err)
	deleteFile(target)

	// Test updating badge without link to badge with link
	createFile(target, "## A Title\nA description\nContent ![Coverage](https://img.shields.io/badge/Coverage-50.01%25-green)")
	err = updateReadme(target, "50.01%", "Coverage", "green", "https://www.cnn.com")
	validateFileContent(t, target, "## A Title\nA description\nContent [![Coverage](https://img.shields.io/badge/Coverage-50.01%25-green)](https://www.cnn.com)")
	assert.Equal(t, nil, err)
	deleteFile(target)

	err = updateReadme("unknown.md", "50.01%", "Coverage", "green", "")
	assert.Equal(t, "open unknown.md: no such file or directory", err.Error())
}
