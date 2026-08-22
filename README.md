# memectl

`memectl` will be a command-line tool for generating memes through Imgflip.
It is being built as an experiment with [OpenSpec](https://github.com/Fission-AI/OpenSpec)
and spec-driven development.

The project currently contains only its bootstrap. It has no Imgflip integration,
configuration, credential handling, or feature commands yet.

## Prerequisites

Install [mise](https://mise.jdx.dev/), then install the pinned toolchain:

```sh
mise install
```

OpenSpec requires Node.js 20.19.0 or newer; the pinned Node.js release satisfies
that requirement.

## Development

```sh
mise run fmt
mise run vet
mise run test
mise run build
```

Run the CLI from source:

```sh
go run . --help
go run . --version
```
