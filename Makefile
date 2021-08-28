
.PHONY: example
example:
	@mkdir -p output
	@cd example;go build -o ../output/example

.PHONY: clean
clean:
	@rm -rf output