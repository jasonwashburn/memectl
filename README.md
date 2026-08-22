# memectl

<p align="center">
  <img src="docs/images/drake-memectl-create.jpg" alt="Drake prefers memectl create over manually clicking through a meme generator" width="320">
</p>

`memectl` is a command-line tool for generating memes through Imgflip.
It is being built as an experiment with [OpenSpec](https://github.com/Fission-AI/OpenSpec)
and spec-driven development.

The first available command lists the public templates that can be used to generate
a meme. It requires no configuration or Imgflip credentials.

## Prerequisites

Install [mise](https://mise.jdx.dev/), then install the pinned toolchain:

```sh
mise install
```

OpenSpec requires Node.js 20.19.0 or newer; the pinned Node.js release satisfies
that requirement.

## Usage

List available Imgflip meme templates:

```sh
go run . get templates
```

## Development

```sh
hk install --mise
mise run fmt
mise run lint
mise run vet
mise run test
mise run coverage
mise run build
mise run release-snapshot
```

`mise run release-snapshot` creates the supported release archives and checksum
manifest in `dist/` without publishing anything to GitHub.

Run all hk checks manually:

```sh
hk check --all
```

Run the CLI from source:

```sh
go run . --help
go run . --version
```
