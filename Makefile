lint:
	golangci-lint run .

test:
	go test -v -race -timeout 60s -count 10 .