install-completions:
	@echo Run: complete -C \"go run ./cmd/argv-complete\" example

# The below will fail if "example" already contains a .argv section!
# watch out, objcopy does some string transformations! It converts tabs to periods -- what else?
build:
	rm -f example
	go build -o example ./cmd/example
	go run ./cmd/argv-embed example cmd/example/arguments.yaml

read:
	objdump -sj .argv example | tail -n +5 | xxd -r

test:
	COMP_LINE="example --" COMP_POINT=10 go run ./cmd/argv-complete ./example
