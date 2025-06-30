install-completions:
	@echo Run: complete -C \"go run ./cmd/argv-complete\" example

# The below will fail if "example" already contains a .argv section!
build:
	rm -f example
	go build -o example ./cmd/example
	go run ./cmd/argv-prepare cmd/example/arguments.json /tmp/args.json
	objcopy \
		--add-section .argv=/tmp/args.json \
		--set-section-flags .argv=noload,readonly \
		example example

read:
	objdump -sj .argv example | tail -n +5 | xxd -r
