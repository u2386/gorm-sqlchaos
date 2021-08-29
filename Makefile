
.PHONY: example
example:
	@mkdir -p output
	@cd example/simple;go build -gcflags='-l -N' -o ../../output/example

.PHONY: clean
clean:
	@rm -rf output

.PHONY: test
test:
	@go test -gcflags='all=-l -N' -v -race -cover ./...