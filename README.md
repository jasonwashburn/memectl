# memectl

<p align="center">
  <img src="docs/images/drake-memectl-create.jpg" alt="Drake prefers memectl create over manually clicking through a meme generator" width="320">
</p>

`memectl` is a command-line tool for generating memes through Imgflip.
It is being built as an experiment with [OpenSpec](https://github.com/Fission-AI/OpenSpec)
and spec-driven development.

Use memectl to list public templates or create a captioned meme from a selected
template. Listing requires no configuration or Imgflip credentials; creation
uses your Imgflip account credentials from the environment.

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

Create a captioned meme from a template. Set your Imgflip account credentials
in the environment first; they are required only for creation, not template
listing:

```sh
export IMGFLIP_USERNAME="your-imgflip-username"
export IMGFLIP_PASSWORD="your-imgflip-password"
go run . create meme 181913649 --text "First caption" --text "Second caption"
```

Use `--text` once for each template text box, in the order Imgflip should
apply them. On success, memectl prints the direct hosted image URL and its
Imgflip page URL. Hosted image URLs are publicly accessible to anyone who
knows the URL.

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
