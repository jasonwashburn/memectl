# Contributing to memectl

`memectl` is a Go project managed with [mise](https://mise.jdx.dev/). It was
built primarily to evaluate OpenSpec and spec-driven development through real
feature work. Most development follows that workflow; small tooling and
non-spec-affecting updates may be made outside it.

## Development Setup

Clone the repository and install [mise](https://mise.jdx.dev/). mise installs
the pinned Go, Node.js, and project tools:

```sh
mise install
```

### Pre-commit Hooks

mise installs the repository hooks automatically after `mise install`. Run the
following command to install or refresh them manually:

```sh
hk install --mise
```

Run all installed hooks manually with:

```sh
hk check --all
```

## Run From Source

Run the CLI directly from a checkout:

```sh
go run . --help
go run . get templates
go run . create meme writing-memes --template 181913649 --text "Writing memes manually" --text "Using memectl"
```

`create meme` requires `IMGFLIP_USERNAME` and `IMGFLIP_PASSWORD` in the
environment. Template listing does not require credentials.

## Quality Checks

Run the full validation suite before submitting a change:

```sh
mise run check
```

For focused checks, use the available mise tasks:

```sh
mise run fmt
mise run lint
mise run vet
mise run test
mise run coverage
mise run build
```

## Release Snapshot

Build the release archives and checksum manifest locally without publishing:

```sh
mise run release-snapshot
```

The snapshot writes its output to `dist/`.
