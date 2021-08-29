
.PHONY: example
example:
	@mkdir -p output
	@cd example;go build -gcflags='-l -N' -o ../output/example

.PHONY: clean
clean:
	@rm -rf output \
		coverage

.PHONY: test
test:
	@mkdir -p coverage
	@go test -race -cover ./...