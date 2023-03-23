# GOBadge
![Coverage](https://img.shields.io/badge/Coverage-76.5%25-brightgreen)

#### 👆 Easily create and insert coverage badge (or any other badge) in your readme with go

Using [https://shields.io](https://shields.io), the executable will generate a coverage badge based on the coverage output. The badge will be inserted on the second line of the readme if it doesn't exist, and will be updated if already present.

## How to use
Install the executable
```
go get github.com/AlexBeauchemin/gobadge
go install github.com/AlexBeauchemin/gobadge
```
Make sure you generate a coverage file with your total coverage, something like
```go
go test ./... -covermode=count -coverprofile=coverage.out fmt
go tool cover -func=coverage.out -o=coverage.out
```

Your output should looks like this, with the total as the last line
```
...
github.com/AlexBeauchemin/go-coverage-badge/gobadge.go:36:	retrieveTotalCoverage	87.5%
github.com/AlexBeauchemin/go-coverage-badge/gobadge.go:53:	updateReadme		80.0%
total:								(statements)		67.4%
```
Then run the executable
```
gobadge -filename=coverage.out
```

## Flags
|Flag                          |Default   |Description|
|------------------------------|----------|-----------|
|filename                      |output.out|File to scan for the coverage total|
|label                         |Coverage  |Left-side content of the badge|
|value                         |          |Right-side content of the badge|
|yellow                        |30        |At what percentage the badge will become yellow instead of red|
|green                         |70        |At what percentage the badge becomes green instead of yellow|
|color                         |          |Force a color for the badge|
|target                        |README.md |Where to insert the badge|
|link                          |          |Optional URL when you click the badge|

## Examples
```
gobadge -filename=coverage.out
gobadge -label="Go Coverage" -value=55.6% -color=blue -target=OTHER_README.md
gobadge -yellow=60 -green=80
gobadge -color=ff69b4
gobadge -link=https://github.com/project/repo/actions/workflows/test.yml
```

## TODO

- Add a silent mode with no output to console/stdout
- Allow to specify a line number in the target file for badge creation
- Allow to specify a template to look for in the target file for badge creation
