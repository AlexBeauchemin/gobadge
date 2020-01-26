package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	coverage := retrieveTotalCoverage("test.txt")
	assert.Equal(t, "33.3%", coverage)
}

func TestUpdateReadme(t *testing.T) {
	target := "test.md"
	input, _ := ioutil.ReadFile(target)
	bk := string(input)

	updateReadme(target, "50.01%", "Coverage", "green")
	input, _ = ioutil.ReadFile(target)
	assert.Equal(t, "## A Title\n![Coverage](https://img.shields.io/badge/Coverage-50.01%25-green)\nA description\n", string(input))

	_ = ioutil.WriteFile("test.md", []byte(bk), 0644)
}
