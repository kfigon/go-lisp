.PHONY:
test:
	go test -timeout=3s ./...

.PHONY:
test-v:
	go test -v -timeout=3s ./...
.PHONY:
clean:
	go clean -testcache
