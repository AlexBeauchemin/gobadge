go test ./... -coverprofile=coverage.out &&
go tool cover -html=coverage.out &&
go build -o gobadge gobadge.go
