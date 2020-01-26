go vet ./... &&
go test ./... -covermode=count -coverprofile=coverage.out fmt &&
go tool cover -func=coverage.out -o=coverage.out &&
go build -o gobadge gobadge.go &&
./gobadge -filename coverage.out -yellow=50 -green=60
