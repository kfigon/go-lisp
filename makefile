.PHONY:
test:
	go test -v -timeout=3s ./...

.PHONY:
clean:
	go clean -testcache
