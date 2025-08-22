# argv

Universal shell completions.

This an extremely early prototype.

## Install

```shell
go install go.wit.com/apps/argv/cmd/argv-complete
go install go.wit.com/apps/argv/cmd/argv-embed
```

## Describe your arguments in YAML

```yaml
name: example
arguments:

- placeholder: PATH
  positional: true
  cardinality: one

- long: --verbose
  placeholder: VERBOSE
  cardinality: zero

- long: --ids
  placeholder: IDS
  cardinality: multiple
```

## Embed the arguments

```shell
$ argv-embed myprogram myarguments.yaml
```

This command validates myarguments.yaml, then embeds it in a text segment within myprogram, which must be an ELF binary. The ELF binary is modified in-place.

## Tell your shell to use arg-complete

```shell
complete -C argv-complete myprogram
```

Currently this only works in bash. More coming soon.

## Get completions

Type `myprogram --` then press TAB, or `./myprogram --` then press TAB and you will see completions. Currently it only completes option names (more coming soon).

## Language support

You can use argv to describe command line options, and get completions, for programs written in any language. Currently, we only support ELF binaries, but we will soon support binaries on other platforms, as well as scripts that are not represented as binaries. Although argv is written in Go, it can provide completions for programs written in any language.


