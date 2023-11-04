lint:
	golangci-lint run --disable-all -E gocritic

test:
	go test -v -race -timeout 60s -count 10 .